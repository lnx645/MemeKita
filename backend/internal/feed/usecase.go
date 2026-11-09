package feed

import (
	"context"
	"time"
	"gorm.io/gorm"
	"github.com/memekita/backend/internal/models"
)

type Usecase struct{ db *gorm.DB }
func NewUsecase(db *gorm.DB) *Usecase { return &Usecase{db: db} }

func (uc *Usecase) GetFeed(ctx context.Context, tab, sectionSlug, cursor string, limit int, includeNsfw bool) ([]models.Post, string, bool, error) {
	if limit <= 0 { limit = 20 }
	if limit > 50 { limit = 50 }
	q := uc.db.WithContext(ctx).Where("status = ?", "published")
	if !includeNsfw { q = q.Where("is_nsfw = ?", false) }
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
	case "hot": q = q.Order("hot_score DESC, created_at DESC")
	case "fresh": q = q.Order("created_at DESC")
	case "trending": q = q.Order("hot_score DESC")
	default: q = q.Order("hot_score DESC")
	}
	var posts []models.Post
	if err := q.Limit(limit+1).Preload("User").Preload("Section").Find(&posts).Error; err != nil { return nil,"",false,err }
	hasMore := len(posts) > limit
	if hasMore { posts = posts[:limit] }
	next := ""
	if hasMore && len(posts)>0 { next = posts[len(posts)-1].CreatedAt.UTC().Format(time.RFC3339Nano) }
	return posts, next, hasMore, nil
}
