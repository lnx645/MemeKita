package duel

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/memekita/backend/internal/models"
	"gorm.io/gorm"
	"time"
)

type Usecase struct{ db *gorm.DB }

func NewUsecase(db *gorm.DB) *Usecase { return &Usecase{db: db} }

func (uc *Usecase) Create(ctx context.Context, a, b uuid.UUID) (uuid.UUID, error) {
	id := uuid.New()
	now := time.Now().UTC()
	ends := now.Add(24 * time.Hour)
	if err := uc.db.WithContext(ctx).Create(&models.Duel{ID: id, PostAID: a, PostBID: b, StartsAt: now, EndsAt: ends, Status: "open"}).Error; err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
func (uc *Usecase) Vote(ctx context.Context, userID, duelID, postID uuid.UUID) error {
	var d models.Duel
	if err := uc.db.WithContext(ctx).Where("id = ?", duelID).First(&d).Error; err != nil {
		return err
	}
	if d.Status != "open" {
		return fmt.Errorf("duel sudah berakhir")
	}
	if postID != d.PostAID && postID != d.PostBID {
		return fmt.Errorf("post bukan bagian duel")
	}
	return uc.db.WithContext(ctx).Exec("INSERT INTO duel_votes (user_id, duel_id, post_id) VALUES (?,?,?) ON CONFLICT DO NOTHING", userID, duelID, postID).Error
}
func (uc *Usecase) CloseExpired(ctx context.Context) error {
	var duels []models.Duel
	if err := uc.db.WithContext(ctx).Where("status = ? AND ends_at < ?", "open", time.Now().UTC()).Find(&duels).Error; err != nil {
		return err
	}
	for _, d := range duels {
		var aCount, bCount int64
		uc.db.WithContext(ctx).Model(&models.DuelVote{}).Where("duel_id = ? AND post_id = ?", d.ID, d.PostAID).Count(&aCount)
		uc.db.WithContext(ctx).Model(&models.DuelVote{}).Where("duel_id = ? AND post_id = ?", d.ID, d.PostBID).Count(&bCount)
		var winner *uuid.UUID
		if aCount > bCount {
			winner = &d.PostAID
		} else if bCount > aCount {
			winner = &d.PostBID
		}
		if err := uc.db.WithContext(ctx).Model(&models.Duel{}).Where("id = ?", d.ID).Updates(map[string]interface{}{"status": "closed", "winner_post_id": winner}).Error; err != nil {
			return err
		}
	}
	return nil
}
