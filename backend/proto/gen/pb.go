package gen

import "context"

// Placeholder for generated protobuf generated — in real CI, run `buf generate` to regenerate from proto/*.proto
// This file provides minimal types so handlers compile without protoc.

type RegisterRequest struct{ Username, Email, Password string }
type LoginRequest struct{ Email, Password string }
type LoginWithGoogleRequest struct{ IdToken string }
type RefreshTokenRequest struct{ RefreshToken string }
type AuthResponse struct {
	AccessToken, RefreshToken string
	ExpiresIn                 int64
	User                      *UserInfo
}
type UserInfo struct{ Id, Username, Email, AvatarUrl, Role string }

type GetFeedRequest struct {
	Tab, SectionSlug, Cursor string
	Limit                    int32
	IncludeNsfw              bool
}
type GetFeedResponse struct {
	Posts      []*Post
	NextCursor string
	HasMore    bool
}
type Post struct {
	Id, UserId, Username, AvatarUrl, SectionId, SectionName, SectionSlug, Title, MediaType, MediaUrl, ThumbUrl, Status string
	Width, Height                                                                                                      int32
	IsNsfw                                                                                                             bool
	Upvotes, Downvotes, CommentCount                                                                                   int64
	HotScore                                                                                                           float64
	Tags                                                                                                               []string
	CreatedAt, ChallengeId                                                                                             string
	MyVote                                                                                                             int32
	MyBookmark                                                                                                         bool
}

type GetPostRequest struct{ PostId string }
type GetPostResponse struct{ Post *Post }
type CreateUploadRequest struct {
	MediaType             int32
	FileSize              int64
	ContentType, FileName string
}
type CreateUploadResponse struct {
	UploadUrl, MediaKey string
	ExpiresIn           int64
}
type CompletePostRequest struct {
	MediaKey, Title, SectionId string
	Tags                       []string
	IsNsfw                     bool
	MediaType                  int32
	Width, Height              int32
	ChallengeId                string
}
type CompletePostResponse struct{ Post *Post }
type DeletePostRequest struct{ PostId string }
type DeletePostResponse struct{ Success bool }

type VotePostRequest struct {
	PostId string
	Value  int32
}
type VotePostResponse struct {
	Upvotes, Downvotes int64
	MyVote             int32
}

type ListCommentsRequest struct {
	PostId, Sort, Cursor string
	Limit                int32
}
type ListCommentsResponse struct {
	Comments   []*Comment
	NextCursor string
	HasMore    bool
}
type Comment struct {
	Id, PostId, UserId, Username, AvatarUrl, ParentId, Body, MediaUrl string
	Upvotes, Downvotes                                                int64
	MyVote                                                            int32
	Status, CreatedAt                                                 string
	Replies                                                           []*Comment
}
type CreateCommentRequest struct{ PostId, ParentId, Body, MediaUrl, MediaKey string }
type CreateCommentResponse struct{ Comment *Comment }
type DeleteCommentRequest struct{ CommentId string }
type DeleteCommentResponse struct{ Success bool }
type VoteCommentRequest struct {
	CommentId string
	Value     int32
}
type VoteCommentResponse struct {
	Upvotes, Downvotes int64
	MyVote             int32
}

// AuthService etc. interfaces for gRPC registration
type AuthServiceServer interface {
	Register(context.Context, *RegisterRequest) (*AuthResponse, error)
	Login(context.Context, *LoginRequest) (*AuthResponse, error)
	LoginWithGoogle(context.Context, *LoginWithGoogleRequest) (*AuthResponse, error)
	RefreshToken(context.Context, *RefreshTokenRequest) (*AuthResponse, error)
}
type FeedServiceServer interface {
	GetFeed(context.Context, *GetFeedRequest) (*GetFeedResponse, error)
}
type PostServiceServer interface {
	GetPost(context.Context, *GetPostRequest) (*GetPostResponse, error)
	CreateUpload(context.Context, *CreateUploadRequest) (*CreateUploadResponse, error)
	CompletePost(context.Context, *CompletePostRequest) (*CompletePostResponse, error)
	DeletePost(context.Context, *DeletePostRequest) (*DeletePostResponse, error)
}
type VoteServiceServer interface {
	VotePost(context.Context, *VotePostRequest) (*VotePostResponse, error)
}
