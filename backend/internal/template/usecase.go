package template

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/memekita/backend/internal/models"
)

type Usecase struct{ db *gorm.DB }
func NewUsecase(db *gorm.DB) *Usecase { return &Usecase{db: db} }

func (uc *Usecase) List(ctx context.Context, query string, officialOnly bool, limit int) ([]models.MemeTemplate, error) {
	q := uc.db.WithContext(ctx).Where("status = ? OR is_official = ?", "approved", true)
	if query != "" { q = q.Where("name ILIKE ?", "%"+query+"%") }
	if officialOnly { q = q.Where("is_official = ?", true) }
	var out []models.MemeTemplate
	if err := q.Order("created_at DESC").Limit(limit).Find(&out).Error; err != nil { return nil, err }
	return out, nil
}
func (uc *Usecase) Submit(ctx context.Context, userID uuid.UUID, name, imageURL string, w, h int) (uuid.UUID, error) {
	id := uuid.New()
	if err := uc.db.WithContext(ctx).Create(&models.MemeTemplate{ID: id, Name: name, ImageURL: imageURL, Width: &w, Height: &h, Status: "pending", CreatedBy: &userID}).Error; err != nil { return uuid.Nil, err }
	return id, nil
}
