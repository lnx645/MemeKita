package challenge

import (
	"context"
	"github.com/google/uuid"
	"github.com/memekita/backend/internal/models"
	"gorm.io/gorm"
	"time"
)

type Usecase struct{ db *gorm.DB }

func NewUsecase(db *gorm.DB) *Usecase { return &Usecase{db: db} }

func (uc *Usecase) Current(ctx context.Context) (*models.Challenge, error) {
	var c models.Challenge
	if err := uc.db.WithContext(ctx).Where("status = ?", "open").Order("starts_at DESC").First(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}
func (uc *Usecase) CreateDaily(ctx context.Context, theme string) error {
	now := time.Now().UTC()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	tomorrow := todayStart.Add(24 * time.Hour)
	if err := uc.db.WithContext(ctx).Model(&models.Challenge{}).Where("status = ? AND ends_at < ?", "open", now).Update("status", "closed").Error; err != nil {
		return err
	}
	return uc.db.WithContext(ctx).Create(&models.Challenge{ID: uuid.New(), Title: theme, Theme: theme, StartsAt: todayStart, EndsAt: tomorrow, Status: "open"}).Error
}
