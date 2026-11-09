package tip

import (
	"context"
	"gorm.io/gorm"
	"github.com/google/uuid"
	"github.com/memekita/backend/internal/models"
)

type Usecase struct{ db *gorm.DB }
func NewUsecase(db *gorm.DB) *Usecase { return &Usecase{db: db} }

func (uc *Usecase) List(ctx context.Context, userID uuid.UUID, typ string, limit int) ([]models.Tip, error) {
	q := uc.db.WithContext(ctx).Order("created_at DESC").Limit(limit)
	if typ == "sent" { q = q.Where("sender_id = ?", userID) } else if typ == "received" { q = q.Where("receiver_id = ?", userID) } else { q = q.Where("sender_id = ? OR receiver_id = ?", userID, userID) }
	var out []models.Tip
	if err := q.Find(&out).Error; err != nil { return nil, err }
	return out, nil
}
