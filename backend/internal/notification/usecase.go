package notification

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/memekita/backend/internal/models"
)

type Usecase struct{ db *gorm.DB }
func NewUsecase(db *gorm.DB) *Usecase { return &Usecase{db: db} }

func (uc *Usecase) Notify(ctx context.Context, userID uuid.UUID, typ string, refID *uuid.UUID, actorID *uuid.UUID, message string) error {
	return uc.db.WithContext(ctx).Create(&models.Notification{ID: uuid.New(), UserID: userID, Type: typ, RefID: refID, ActorID: actorID, Message: &message}).Error
}
func (uc *Usecase) List(ctx context.Context, userID uuid.UUID, filter string, limit int) ([]models.Notification, error) {
	q := uc.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Limit(limit)
	if filter == "unread" { q = q.Where("read_at IS NULL") }
	var out []models.Notification
	if err := q.Find(&out).Error; err != nil { return nil, err }
	return out, nil
}
func (uc *Usecase) MarkRead(ctx context.Context, userID, notifID uuid.UUID) error {
	return uc.db.WithContext(ctx).Model(&models.Notification{}).Where("id = ? AND user_id = ?", notifID, userID).Update("read_at", gorm.Expr("NOW()")).Error
}
