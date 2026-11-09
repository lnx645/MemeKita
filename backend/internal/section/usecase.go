package section

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/memekita/backend/internal/models"
)

type Usecase struct{ db *gorm.DB }
func NewUsecase(db *gorm.DB) *Usecase { return &Usecase{db: db} }

func (uc *Usecase) List(ctx context.Context, includeRegional bool) ([]models.Section, error) {
	var secs []models.Section
	q := uc.db.WithContext(ctx).Order("name ASC")
	if !includeRegional { q = q.Where("is_regional = ?", false) }
	if err := q.Find(&secs).Error; err != nil { return nil, err }
	return secs, nil
}
func (uc *Usecase) GetBySlug(ctx context.Context, slug string) (*models.Section, error) {
	var s models.Section
	if err := uc.db.WithContext(ctx).Where("slug = ?", slug).First(&s).Error; err != nil { return nil, err }
	return &s, nil
}
func (uc *Usecase) Follow(ctx context.Context, userID, sectionID uuid.UUID) error {
	return uc.db.WithContext(ctx).Exec("INSERT INTO follows (follower_id, target_type, target_id) VALUES (?,'section',?) ON CONFLICT DO NOTHING", userID, sectionID).Error
}
func (uc *Usecase) Unfollow(ctx context.Context, userID, sectionID uuid.UUID) error {
	return uc.db.WithContext(ctx).Exec("DELETE FROM follows WHERE follower_id=? AND target_type='section' AND target_id=?", userID, sectionID).Error
}
