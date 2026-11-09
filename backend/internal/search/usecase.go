package search

import (
	"context"
	"gorm.io/gorm"
	"github.com/memekita/backend/internal/models"
)

type Usecase struct{ db *gorm.DB }
func NewUsecase(db *gorm.DB) *Usecase { return &Usecase{db: db} }

func (uc *Usecase) Search(ctx context.Context, query, typ string, limit int) (map[string]interface{}, error) {
	out := map[string]interface{}{"posts":[]interface{}{},"users":[]interface{}{},"tags":[]interface{}{}}
	if query == "" { return out, nil }
	if typ == "all" || typ == "post" {
		var posts []models.Post
		if err := uc.db.WithContext(ctx).Where("status = ? AND title ILIKE ?", "published", "%"+query+"%").Order("hot_score DESC").Limit(limit).Preload("User").Find(&posts).Error; err != nil {
			return nil, err
		}
		out["posts"]=posts
	}
	if typ == "all" || typ == "user" {
		var users []models.User
		if err := uc.db.WithContext(ctx).Where("username ILIKE ?", "%"+query+"%").Limit(10).Find(&users).Error; err != nil {
			return nil, err
		}
		out["users"]=users
	}
	return out, nil
}
