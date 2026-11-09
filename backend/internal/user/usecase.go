package user

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/memekita/backend/internal/models"
)

type Usecase struct{ db *gorm.DB }
func NewUsecase(db *gorm.DB) *Usecase { return &Usecase{db: db} }

func (uc *Usecase) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var u models.User
	if err := uc.db.WithContext(ctx).Where("id = ?", id).First(&u).Error; err != nil { return nil, err }
	return &u, nil
}
func (uc *Usecase) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var u models.User
	if err := uc.db.WithContext(ctx).Where("username = ?", username).First(&u).Error; err != nil { return nil, err }
	return &u, nil
}
func (uc *Usecase) Follow(ctx context.Context, followerID, targetID uuid.UUID) error {
	return uc.db.WithContext(ctx).Exec("INSERT INTO follows (follower_id, target_type, target_id) VALUES (?,'user',?) ON CONFLICT DO NOTHING", followerID, targetID).Error
}
func (uc *Usecase) Unfollow(ctx context.Context, followerID, targetID uuid.UUID) error {
	return uc.db.WithContext(ctx).Exec("DELETE FROM follows WHERE follower_id=? AND target_type='user' AND target_id=?", followerID, targetID).Error
}
