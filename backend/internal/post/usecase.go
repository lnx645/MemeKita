package post

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/memekita/backend/internal/models"
	"gorm.io/gorm"
)

func HotScore(up, down int, ageHours float64) float64 {
	score := float64(up - down)
	denom := math.Pow(ageHours+2, 1.5)
	if denom == 0 {
		return score
	}
	return score / denom
}

type Usecase struct{ db *gorm.DB }

func NewUsecase(db *gorm.DB) *Usecase { return &Usecase{db: db} }

func (uc *Usecase) Create(ctx context.Context, userID, sectionID uuid.UUID, title, mediaType, mediaKey string, width, height int, isNsfw bool, tags []string, challengeID *uuid.UUID) (*models.Post, error) {
	if title == "" {
		return nil, fmt.Errorf("judul wajib diisi")
	}
	if len(title) > 200 {
		return nil, fmt.Errorf("judul maksimal 200 karakter")
	}
	p := &models.Post{
		ID: uuid.New(), UserID: userID, SectionID: sectionID, Title: title, MediaType: mediaType, MediaKey: &mediaKey, Width: &width, Height: &height, Status: "processing", IsNsfw: isNsfw, ChallengeID: challengeID,
	}
	if err := uc.db.WithContext(ctx).Create(p).Error; err != nil {
		return nil, err
	}
	for _, t := range tags {
		if err := uc.db.WithContext(ctx).Exec("INSERT INTO post_tags (post_id, tag) VALUES (?,?) ON CONFLICT DO NOTHING", p.ID, t).Error; err != nil {
			return nil, err
		}
	}
	if challengeID != nil {
		if err := uc.db.WithContext(ctx).Exec("INSERT INTO challenge_entries (challenge_id, post_id) VALUES (?,?) ON CONFLICT DO NOTHING", *challengeID, p.ID).Error; err != nil {
			return nil, err
		}
	}
	return p, nil
}

func (uc *Usecase) Get(ctx context.Context, id uuid.UUID) (*models.Post, error) {
	var p models.Post
	err := uc.db.WithContext(ctx).Preload("User").Preload("Section").Where("id = ?", id).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (uc *Usecase) ListFeed(ctx context.Context, tab, sectionSlug, cursor string, limit int, includeNsfw bool) ([]models.Post, string, bool, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 50 {
		limit = 50
	}
	q := uc.db.WithContext(ctx).Where("status = ?", "published")
	if !includeNsfw {
		q = q.Where("is_nsfw = ?", false)
	}
	if sectionSlug != "" {
		var sec models.Section
		if err := uc.db.WithContext(ctx).Where("slug = ?", sectionSlug).First(&sec).Error; err == nil {
			q = q.Where("section_id = ?", sec.ID)
		}
	}
	if cursor != "" {
		if t, err := time.Parse(time.RFC3339Nano, cursor); err == nil {
			q = q.Where("created_at < ?", t.UTC())
		}
	}
	switch tab {
	case "hot":
		q = q.Order("hot_score DESC, created_at DESC")
	case "fresh":
		q = q.Order("created_at DESC")
	case "trending":
		q = q.Order("hot_score DESC")
	default:
		q = q.Order("hot_score DESC")
	}
	var posts []models.Post
	err := q.Limit(limit + 1).Preload("User").Preload("Section").Find(&posts).Error
	if err != nil {
		return nil, "", false, err
	}
	hasMore := len(posts) > limit
	if hasMore {
		posts = posts[:limit]
	}
	next := ""
	if hasMore && len(posts) > 0 {
		next = posts[len(posts)-1].CreatedAt.UTC().Format(time.RFC3339Nano)
	}
	return posts, next, hasMore, nil
}
