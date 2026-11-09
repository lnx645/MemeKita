package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"

	"github.com/memekita/backend/internal/config"
	"github.com/memekita/backend/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, user *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	FindByGoogleID(ctx context.Context, googleID string) (*models.User, error)
}

type GormRepo struct{ db *gorm.DB }

func NewGormRepo(db *gorm.DB) *GormRepo { return &GormRepo{db: db} }

func (r *GormRepo) Create(ctx context.Context, u *models.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *GormRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *GormRepo) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var u models.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *GormRepo) FindByGoogleID(ctx context.Context, gid string) (*models.User, error) {
	var u models.User
	err := r.db.WithContext(ctx).Where("google_id = ?", gid).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

type Usecase struct {
	repo Repository
	db   *gorm.DB
	cfg  config.Config
}

func NewUsecase(repo Repository, db *gorm.DB, cfg config.Config) *Usecase {
	return &Usecase{repo: repo, db: db, cfg: cfg}
}

func (uc *Usecase) Register(ctx context.Context, username, email, password string) (*models.User, string, string, error) {
	if len(username) < 3 || len(username) > 30 {
		return nil, "", "", errors.New("username harus 3-30 karakter")
	}
	if len(password) < 8 {
		return nil, "", "", errors.New("password minimal 8 karakter")
	}
	if _, err := uc.repo.FindByEmail(ctx, email); err == nil {
		return nil, "", "", errors.New("email sudah terdaftar")
	}
	hash := hashPassword(password)
	u := &models.User{
		ID:           uuid.New(),
		Username:     username,
		Email:        email,
		PasswordHash: &hash,
		Role:         "user",
		Status:       "active",
	}
	if err := uc.repo.Create(ctx, u); err != nil {
		return nil, "", "", fmt.Errorf("register: %w", err)
	}
	if err := uc.db.WithContext(ctx).Exec("INSERT INTO user_coins (user_id, balance_coins) VALUES (?,0) ON CONFLICT DO NOTHING", u.ID).Error; err != nil {
		return nil, "", "", err
	}
	access, refresh, err := uc.generateTokens(u)
	if err != nil {
		return nil, "", "", err
	}
	return u, access, refresh, nil
}

func (uc *Usecase) Login(ctx context.Context, email, password string) (*models.User, string, string, error) {
	u, err := uc.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", "", errors.New("email atau password salah")
		}
		return nil, "", "", err
	}
	if u.PasswordHash == nil || !verifyPassword(*u.PasswordHash, password) {
		return nil, "", "", errors.New("email atau password salah")
	}
	if u.Status == "blocked" {
		return nil, "", "", errors.New("akun diblokir")
	}
	access, refresh, err := uc.generateTokens(u)
	return u, access, refresh, err
}

func (uc *Usecase) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *Usecase) VerifyAccessToken(token string) (jwt.MapClaims, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) { return []byte(uc.cfg.JWTSecret), nil })
	if err != nil || !t.Valid {
		return nil, errors.New("token tidak valid")
	}
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("token tidak valid")
	}
	return claims, nil
}

func (uc *Usecase) generateTokens(u *models.User) (string, string, error) {
	now := time.Now().UTC()
	accessClaims := jwt.MapClaims{
		"sub": u.ID.String(), "username": u.Username, "role": u.Role,
		"exp": now.Add(time.Duration(uc.cfg.AccessTTL) * time.Second).Unix(),
		"iat": now.Unix(),
	}
	refreshClaims := jwt.MapClaims{
		"sub": u.ID.String(), "type": "refresh",
		"exp": now.Add(time.Duration(uc.cfg.RefreshTTL) * time.Second).Unix(),
		"iat": now.Unix(),
	}
	access, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(uc.cfg.JWTSecret))
	if err != nil {
		return "", "", err
	}
	refresh, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(uc.cfg.JWTSecret))
	return access, refresh, err
}

func (uc *Usecase) FindOrCreateGoogleUser(ctx context.Context, email, googleID, name, avatar string) (*models.User, error) {
	if u, err := uc.repo.FindByGoogleID(ctx, googleID); err == nil {
		return u, nil
	}
	if u, err := uc.repo.FindByEmail(ctx, email); err == nil {
		if err := uc.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", u.ID).Updates(map[string]interface{}{"google_id": googleID}).Error; err != nil {
			return nil, err
		}
		return u, nil
	}
	username := generateUsername(name, email)
	hash := hashPassword(randomPassword())
	u := &models.User{
		ID:           uuid.New(),
		Username:     username,
		Email:        email,
		PasswordHash: &hash,
		GoogleID:     &googleID,
		AvatarURL:    &avatar,
		Role:         "user",
		Status:       "active",
	}
	if err := uc.repo.Create(ctx, u); err != nil {
		return nil, err
	}
	if err := uc.db.WithContext(ctx).Exec("INSERT INTO user_coins (user_id, balance_coins) VALUES (?,0) ON CONFLICT DO NOTHING", u.ID).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func generateUsername(name, email string) string {
	base := name
	if base == "" {
		base = email
		if at := indexOf(email, "@"); at > 0 {
			base = email[:at]
		}
	}
	sanitized := ""
	for _, r := range base {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			sanitized += string(r)
		} else {
			sanitized += "_"
		}
	}
	if len(sanitized) > 20 {
		sanitized = sanitized[:20]
	}
	if sanitized == "" {
		sanitized = "user"
	}
	b := make([]byte, 3)
	if _, err := rand.Read(b); err != nil {
		return sanitized + "_0000"
	}
	suffix := base64.RawURLEncoding.EncodeToString(b)[:4]
	return sanitized + "_" + suffix
}

func indexOf(s, sep string) int {
	for i := 0; i < len(s); i++ {
		if string(s[i]) == sep {
			return i
		}
	}
	return -1
}

func randomPassword() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "fallback-password-123"
	}
	return base64.StdEncoding.EncodeToString(b)
}

func hashPassword(pw string) string {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		// fallback salt
		for i := range salt {
			salt[i] = byte(i)
		}
	}
	hash := argon2.IDKey([]byte(pw), salt, 1, 64*1024, 1, 32)
	return base64.StdEncoding.EncodeToString(append(salt, hash...))
}

func verifyPassword(hash, pw string) bool {
	b, err := base64.StdEncoding.DecodeString(hash)
	if err != nil || len(b) < 16+32 {
		return false
	}
	salt := b[:16]
	expected := b[16:]
	got := argon2.IDKey([]byte(pw), salt, 1, 64*1024, 1, 32)
	if len(got) != len(expected) {
		return false
	}
	for i := range got {
		if got[i] != expected[i] {
			return false
		}
	}
	return true
}
