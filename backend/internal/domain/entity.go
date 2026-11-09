package domain

import "time"

// Domain entities — pure data, no gorm tags, no DB logic.

type User struct {
	ID        string
	Username  string
	Email     string
	AvatarURL string
	Role      string
	Status    string
	CreatedAt time.Time
}

type Post struct {
	ID           string
	UserID       string
	Username     string
	AvatarURL    string
	SectionID    string
	SectionName  string
	SectionSlug  string
	Title        string
	MediaType    string
	MediaURL     string
	ThumbURL     string
	Status       string
	IsNsfw       bool
	Upvotes      int
	Downvotes    int
	CommentCount int
	HotScore     float64
	Tags         []string
	CreatedAt    time.Time
	MyVote       int
}

type Comment struct {
	ID        string
	PostID    string
	UserID    string
	Username  string
	AvatarURL string
	ParentID  *string
	Body      string
	MediaURL  *string
	Upvotes   int
	Downvotes int
	Status    string
	CreatedAt time.Time
	Replies   []*Comment
}

type ReactionCounts struct {
	Ngakak int
	Relate int
	Anjay  int
	Sedih  int
	Marah  int
	Total  int
}
