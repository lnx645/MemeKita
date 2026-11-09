package comment

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/memekita/backend/internal/models"
)

type Usecase struct{ db *gorm.DB }
func NewUsecase(db *gorm.DB) *Usecase { return &Usecase{db: db} }

func (uc *Usecase) List(ctx context.Context, postID uuid.UUID, sort string) ([]*models.Comment, error) {
	q := uc.db.WithContext(ctx).Where("post_id = ? AND status = ?", postID, "visible").Preload("User")
	if sort == "top" { q = q.Order("upvotes DESC, created_at DESC") } else { q = q.Order("created_at DESC") }
	var comments []*models.Comment
	if err := q.Find(&comments).Error; err != nil { return nil, err }
	// nest
	m := map[string][]*models.Comment{}
	var roots []*models.Comment
	for _, c := range comments {
		key := ""
		if c.ParentID != nil { key = c.ParentID.String() }
		m[key] = append(m[key], c)
	}
	var build func(pid string) []*models.Comment
	build = func(pid string) []*models.Comment {
		children := m[pid]
		for _, ch := range children {
			// recursively attach replies via separate query or in-memory
			chID := ch.ID.String()
			if subs, ok := m[chID]; ok {
				// attach already via build
				_ = subs // keep to avoid unused, will be used in nesting
			}
		}
		return children
	}
	roots = build("")
	// simple: return flat with parent grouping for now — full nesting requires model Replies field
	return roots, nil
}

func (uc *Usecase) Create(ctx context.Context, postID, userID uuid.UUID, parentID *uuid.UUID, body string, mediaURL *string) (*models.Comment, error) {
	if body == "" && (mediaURL == nil || *mediaURL == "") { return nil, fmt.Errorf("komentar tidak boleh kosong") }
	if len(body) > 2000 { return nil, fmt.Errorf("komentar maksimal 2000") }
	c := &models.Comment{ID: uuid.New(), PostID: postID, UserID: userID, ParentID: parentID, Body: body, MediaURL: mediaURL, Status: "visible"}
	if err := uc.db.WithContext(ctx).Create(c).Error; err != nil { return nil, err }
	if err := uc.db.WithContext(ctx).Model(&models.Post{}).Where("id = ?", postID).UpdateColumn("comment_count", gorm.Expr("comment_count + 1")).Error; err != nil {
		return nil, err
	}
	return c, nil
}

func (uc *Usecase) Delete(ctx context.Context, commentID, actorID uuid.UUID, actorRole string) error {
	var c models.Comment
	if err := uc.db.WithContext(ctx).Where("id = ?", commentID).First(&c).Error; err != nil {
		if err == gorm.ErrRecordNotFound { return fmt.Errorf("komentar tidak ditemukan") }
		return err
	}
	if c.UserID != actorID && actorRole != "admin" && actorRole != "moderator" {
		return fmt.Errorf("tidak berhak")
	}
	return uc.db.WithContext(ctx).Delete(&c).Error
}
