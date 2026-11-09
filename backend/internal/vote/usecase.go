package vote

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/memekita/backend/internal/models"
	"github.com/memekita/backend/internal/post"
)

type Usecase struct{ db *gorm.DB }
func NewUsecase(db *gorm.DB) *Usecase { return &Usecase{db: db} }

func (uc *Usecase) VotePost(ctx context.Context, userID, postID uuid.UUID, value int) (int, int, int, error) {
	if value != 1 && value != -1 && value != 0 {
		return 0,0,0, fmt.Errorf("value harus -1,0,1")
	}
	var up, down, myVote int
	err := uc.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var p models.Post
		if err := tx.WithContext(ctx).Where("id = ?", postID).First(&p).Error; err != nil {
			if err == gorm.ErrRecordNotFound { return fmt.Errorf("post tidak ditemukan") }
			return err
		}
		var existing models.Vote
		err := tx.WithContext(ctx).Where("user_id = ? AND post_id = ?", userID, postID).First(&existing).Error
		exists := err == nil
		old := 0
		if exists { old = existing.Value }
		up = p.Upvotes; down = p.Downvotes

		switch {
		case !exists && value == 0:
			myVote = 0
		case !exists && value != 0:
			if err := tx.WithContext(ctx).Create(&models.Vote{UserID: userID, PostID: postID, Value: value}).Error; err != nil { return err }
			up += boolToInt(value==1); down += boolToInt(value==-1); myVote = value
		case exists && value == 0:
			if err := tx.WithContext(ctx).Where("user_id = ? AND post_id = ?", userID, postID).Delete(&models.Vote{}).Error; err != nil { return err }
			up -= boolToInt(old==1); down -= boolToInt(old==-1); myVote = 0
		case exists && old == value:
			myVote = old
		case exists && old != value:
			if err := tx.WithContext(ctx).Model(&models.Vote{}).Where("user_id = ? AND post_id = ?", userID, postID).Update("value", value).Error; err != nil { return err }
			up += -boolToInt(old==1) + boolToInt(value==1)
			down += -boolToInt(old==-1) + boolToInt(value==-1)
			myVote = value
		}
		if up < 0 { up = 0 }
		if down < 0 { down = 0 }
		ageHours := time.Since(p.CreatedAt.UTC()).Hours()
		hot := post.HotScore(up, down, ageHours)
		if err := tx.WithContext(ctx).Model(&models.Post{}).Where("id = ?", postID).Updates(map[string]interface{}{"upvotes": up, "downvotes": down, "hot_score": hot}).Error; err != nil { return err }
		return nil
	})
	return up, down, myVote, err
}
func boolToInt(b bool) int { if b { return 1 }; return 0 }
