package moderation

import (
	"context"
	"fmt"
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/memekita/backend/internal/models"
)

type Usecase struct{ db *gorm.DB }
func NewUsecase(db *gorm.DB) *Usecase { return &Usecase{db: db} }

func (uc *Usecase) CreateReport(ctx context.Context, reporterID uuid.UUID, targetType string, targetID uuid.UUID, reason string) (uuid.UUID, error) {
	if reason == "" { return uuid.Nil, fmt.Errorf("alasan wajib") }
	id := uuid.New()
	if err := uc.db.WithContext(ctx).Create(&models.Report{ID: id, ReporterID: reporterID, TargetType: targetType, TargetID: targetID, Reason: reason, Status: "pending"}).Error; err != nil { return uuid.Nil, err }
	return id, nil
}
func (uc *Usecase) Resolve(ctx context.Context, reportID, moderatorID uuid.UUID, action, note string) error {
	return uc.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var r models.Report
		if err := tx.WithContext(ctx).Where("id = ?", reportID).First(&r).Error; err != nil { return err }
		if r.Status != "pending" { return fmt.Errorf("sudah ditangani") }
		now := time.Now().UTC()
		if err := tx.WithContext(ctx).Model(&models.Report{}).Where("id = ?", reportID).Updates(map[string]interface{}{"status":"handled","handled_by":moderatorID,"handled_at":now,"note":note}).Error; err != nil { return err }
		if err := tx.WithContext(ctx).Create(&models.ModAction{ID: uuid.New(), ModeratorID: moderatorID, TargetType: r.TargetType, TargetID: r.TargetID, Action: action, Note: &note}).Error; err != nil { return err }
		meta := fmt.Sprintf(`{"report_id":"%s"}`, reportID.String())
		if err := tx.WithContext(ctx).Create(&models.AuditLog{ID: uuid.New(), ActorID: &moderatorID, Action: "moderation:"+action, TargetType: &r.TargetType, TargetID: &r.TargetID, Metadata: meta}).Error; err != nil {
			return err
		}
		switch action {
		case "hide":
			if r.TargetType=="post" { if err := tx.WithContext(ctx).Model(&models.Post{}).Where("id = ?", r.TargetID).Update("status","hidden").Error; err != nil { return err } }
			if r.TargetType=="comment" { if err := tx.WithContext(ctx).Model(&models.Comment{}).Where("id = ?", r.TargetID).Update("status","hidden").Error; err != nil { return err } }
		case "remove":
			if r.TargetType=="post" { if err := tx.WithContext(ctx).Model(&models.Post{}).Where("id = ?", r.TargetID).Update("status","removed").Error; err != nil { return err } }
			if r.TargetType=="comment" { if err := tx.WithContext(ctx).Model(&models.Comment{}).Where("id = ?", r.TargetID).Update("status","removed").Error; err != nil { return err } }
		case "remove_and_block":
			if r.TargetType=="post" {
				if err := tx.WithContext(ctx).Model(&models.Post{}).Where("id = ?", r.TargetID).Update("status","removed").Error; err != nil { return err }
				var p models.Post
				if err := tx.WithContext(ctx).Where("id = ?", r.TargetID).First(&p).Error; err == nil {
					if err := tx.WithContext(ctx).Model(&models.User{}).Where("id = ?", p.UserID).Update("status","blocked").Error; err != nil { return err }
				}
			}
		}
		var author uuid.UUID
		if r.TargetType=="post" { var p models.Post; if err := tx.WithContext(ctx).Where("id = ?", r.TargetID).First(&p).Error; err == nil { author = p.UserID } }
		if r.TargetType=="comment" { var c models.Comment; if err := tx.WithContext(ctx).Where("id = ?", r.TargetID).First(&c).Error; err == nil { author = c.UserID } }
		if author != uuid.Nil {
			if err := tx.WithContext(ctx).Create(&models.Notification{ID: uuid.New(), UserID: author, Type: "moderation", RefID: &r.TargetID, Message: stringPtr("Kontenmu ditindak moderasi")}).Error; err != nil { return err }
		}
		return nil
	})
}
func stringPtr(s string) *string { return &s }
