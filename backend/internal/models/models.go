package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GORM models — persistence only, no business logic.
// Entity domain ada di internal/domain, repository memetakan entity ↔ model.

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Username     string    `gorm:"uniqueIndex;not null"`
	Email        string    `gorm:"uniqueIndex;not null"`
	PasswordHash *string
	GoogleID     *string `gorm:"uniqueIndex"`
	AvatarURL    *string
	Bio          *string
	Website      *string
	Role         string `gorm:"default:user"`
	Status       string `gorm:"default:active"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

type Section struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Slug        string    `gorm:"uniqueIndex;not null"`
	Name        string    `gorm:"not null"`
	Description *string
	IsNsfw      bool `gorm:"default:false"`
	IsRegional  bool `gorm:"default:false"`
	Icon        *string
	PostCount   int       `gorm:"default:0"`
	ModeratorID *uuid.UUID `gorm:"type:uuid"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Challenge struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Title     string
	Theme     string
	StartsAt  time.Time
	EndsAt    time.Time
	Status    string `gorm:"default:open"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Post struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID       uuid.UUID  `gorm:"type:uuid;not null;index"`
	SectionID    uuid.UUID  `gorm:"type:uuid;not null;index"`
	ChallengeID  *uuid.UUID `gorm:"type:uuid;index"`
	Title        string     `gorm:"not null"`
	MediaType    string     `gorm:"not null"`
	MediaURL     *string
	ThumbURL     *string
	MediaKey     *string
	Width        *int
	Height       *int
	Status       string `gorm:"default:processing;index"`
	IsNsfw       bool   `gorm:"default:false"`
	Upvotes      int    `gorm:"default:0"`
	Downvotes    int    `gorm:"default:0"`
	CommentCount int    `gorm:"default:0"`
	HotScore     float64 `gorm:"default:0;index"`
	CreatedAt    time.Time `gorm:"autoCreateTime;index"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`

	User    User    `gorm:"foreignKey:UserID"`
	Section Section `gorm:"foreignKey:SectionID"`
}

type PostTag struct {
	PostID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Tag    string    `gorm:"primaryKey"`
}

type Vote struct {
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	PostID    uuid.UUID `gorm:"type:uuid;primaryKey;index"`
	Value     int
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type Comment struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	PostID    uuid.UUID  `gorm:"type:uuid;not null;index"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index"`
	ParentID  *uuid.UUID `gorm:"type:uuid;index"`
	Body      string
	MediaURL  *string
	MediaKey  *string
	Upvotes   int    `gorm:"default:0"`
	Downvotes int    `gorm:"default:0"`
	Status    string `gorm:"default:visible"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	User      User `gorm:"foreignKey:UserID"`
}

type CommentVote struct {
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	CommentID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Value     int
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type Follow struct {
	FollowerID uuid.UUID `gorm:"type:uuid;primaryKey"`
	TargetType string    `gorm:"primaryKey"`
	TargetID   uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}

type Notification struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID  `gorm:"type:uuid;index;not null"`
	Type      string
	RefID     *uuid.UUID `gorm:"type:uuid"`
	ActorID   *uuid.UUID `gorm:"type:uuid"`
	Message   *string
	ReadAt    *time.Time
	CreatedAt time.Time `gorm:"autoCreateTime;index"`
}

type Report struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ReporterID  uuid.UUID  `gorm:"type:uuid;not null"`
	TargetType  string
	TargetID    uuid.UUID `gorm:"type:uuid;not null"`
	Reason      string
	Description *string
	Status      string `gorm:"default:pending;index"`
	HandledBy   *uuid.UUID `gorm:"type:uuid"`
	HandledAt   *time.Time
	Note        *string
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

type ModAction struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ModeratorID uuid.UUID `gorm:"type:uuid"`
	TargetType  string
	TargetID    uuid.UUID `gorm:"type:uuid"`
	Action      string
	Note        *string
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

type MemeTemplate struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name       string
	ImageURL   string
	Width      *int
	Height     *int
	TextBoxes  string `gorm:"type:jsonb;default:'{}'"`
	IsOfficial bool   `gorm:"default:false"`
	Status     string `gorm:"default:pending"`
	CreatedBy  *uuid.UUID `gorm:"type:uuid"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

type Duel struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	PostAID      uuid.UUID  `gorm:"type:uuid;not null;column:post_a_id"`
	PostBID      uuid.UUID  `gorm:"type:uuid;not null;column:post_b_id"`
	StartsAt     time.Time
	EndsAt       time.Time
	Status       string     `gorm:"default:open;index"`
	WinnerPostID *uuid.UUID `gorm:"type:uuid"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type DuelVote struct {
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	DuelID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	PostID    uuid.UUID `gorm:"type:uuid"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type ChallengeEntry struct {
	ChallengeID uuid.UUID `gorm:"type:uuid;primaryKey"`
	PostID      uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

type Tip struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	SenderID   uuid.UUID `gorm:"type:uuid;not null"`
	ReceiverID uuid.UUID `gorm:"type:uuid;not null"`
	PostID     *uuid.UUID `gorm:"type:uuid"`
	AmountCoins int
	Status     string `gorm:"default:completed"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}

type UserCoin struct {
	UserID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	BalanceCoins int       `gorm:"default:0"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

type CoinPackage struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Coins     int
	PriceIdr  int `gorm:"column:price_idr"`
	PlaySku   string `gorm:"uniqueIndex"`
	IsActive  bool `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CoinTransaction struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID       uuid.UUID  `gorm:"type:uuid;index;not null"`
	Type         string
	AmountCoins  int
	BalanceAfter int
	RefID        *uuid.UUID `gorm:"type:uuid"`
	CreatedAt    time.Time `gorm:"autoCreateTime;index"`
}

type PostReaction struct {
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	PostID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	Reaction  string
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type Bookmark struct {
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	PostID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type AuditLog struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ActorID    *uuid.UUID `gorm:"type:uuid"`
	Action     string
	TargetType *string
	TargetID   *uuid.UUID `gorm:"type:uuid"`
	Metadata   string `gorm:"type:jsonb;default:'{}'"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}

func (User) TableName() string { return "users" }
func (Section) TableName() string { return "sections" }
func (Post) TableName() string { return "posts" }
func (PostTag) TableName() string { return "post_tags" }
func (Vote) TableName() string { return "votes" }
func (Comment) TableName() string { return "comments" }
func (CommentVote) TableName() string { return "comment_votes" }
func (Follow) TableName() string { return "follows" }
func (Notification) TableName() string { return "notifications" }
func (Report) TableName() string { return "reports" }
func (ModAction) TableName() string { return "mod_actions" }
func (MemeTemplate) TableName() string { return "meme_templates" }
func (Duel) TableName() string { return "duels" }
func (DuelVote) TableName() string { return "duel_votes" }
func (ChallengeEntry) TableName() string { return "challenge_entries" }
func (Tip) TableName() string { return "tips" }
func (UserCoin) TableName() string { return "user_coins" }
func (CoinPackage) TableName() string { return "coin_packages" }
func (CoinTransaction) TableName() string { return "coin_transactions" }
func (PostReaction) TableName() string { return "post_reactions" }
func (Bookmark) TableName() string { return "bookmarks" }
func (AuditLog) TableName() string { return "audit_logs" }

// BeforeCreate hooks to set UUID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil { u.ID = uuid.New() }
	u.CreatedAt = time.Now().UTC()
	u.UpdatedAt = time.Now().UTC()
	return nil
}
func (s *Section) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil { s.ID = uuid.New() }
	return nil
}
