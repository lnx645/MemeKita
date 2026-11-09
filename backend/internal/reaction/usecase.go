package reaction

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var valid = map[string]bool{"ngakak": true, "relate": true, "anjay": true, "sedih": true, "marah": true}

type Usecase struct{ db *gorm.DB }

func NewUsecase(db *gorm.DB) *Usecase { return &Usecase{db: db} }

func (uc *Usecase) React(ctx context.Context, userID, postID uuid.UUID, reaction string) error {
	if !valid[reaction] {
		return fmt.Errorf("reaksi tidak valid")
	}
	return uc.db.WithContext(ctx).Exec("INSERT INTO post_reactions (user_id, post_id, reaction) VALUES (?,?,?) ON CONFLICT (user_id, post_id) DO UPDATE SET reaction = ?", userID, postID, reaction, reaction).Error
}
func (uc *Usecase) Counts(ctx context.Context, postID uuid.UUID) (map[string]int, string, float64) {
	var rows []struct {
		Reaction string
		Count    int
	}
	uc.db.WithContext(ctx).Raw("SELECT reaction, count(*) as count FROM post_reactions WHERE post_id = ? GROUP BY reaction", postID).Scan(&rows)
	m := map[string]int{"ngakak": 0, "relate": 0, "anjay": 0, "sedih": 0, "marah": 0}
	top := ""
	max := 0
	total := 0
	for _, r := range rows {
		m[r.Reaction] = r.Count
		total += r.Count
		if r.Count > max {
			max = r.Count
			top = r.Reaction
		}
	}
	percent := 0.0
	if total > 0 {
		percent = float64(max) / float64(total) * 100
	}
	return m, top, percent
}
