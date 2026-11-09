package gen

import "context"

// Additional types for remaining services

type ListSectionsRequest struct{ IncludeRegional bool }
type ListSectionsResponse struct{ Sections []*SectionInfo }
type SectionInfo struct {
	Id, Slug, Name, Description string
	IsNsfw, IsRegional          bool
	Icon                        string
	PostCount                   int64
	IsFollowing                 bool
	ModeratorId                 string
}
type GetSectionRequest struct{ Slug string }
type GetSectionResponse struct{ Section *SectionInfo }
type FollowSectionRequest struct{ SectionId string }
type FollowSectionResponse struct{ Success bool }
type UnfollowSectionRequest struct{ SectionId string }
type UnfollowSectionResponse struct{ Success bool }

type GetProfileRequest struct{ UserId, Username string }
type GetProfileResponse struct{ User *UserProfile }
type UserProfile struct {
	Id, Username, Email, AvatarUrl, Bio, Website, Role string
	PostCount, FollowerCount, FollowingCount           int64
	IsFollowing, IsBlocked                             bool
	CreatedAt                                          string
}
type UpdateProfileRequest struct{ DisplayName, Bio, Website, AvatarKey string }
type UpdateProfileResponse struct{ User *UserProfile }
type FollowUserRequest struct{ TargetUserId string }
type FollowUserResponse struct{ Success bool }

type SearchRequest struct {
	Query, Type, Cursor string
	Limit               int32
}
type SearchResponse struct {
	Posts      []*Post
	Users      []*UserInfo
	Tags       []*TagInfo
	NextCursor string
	HasMore    bool
}
type TagInfo struct {
	Tag       string
	PostCount int64
}

type ListNotificationsRequest struct {
	Cursor string
	Limit  int32
	Filter string
}
type ListNotificationsResponse struct {
	Notifications []*Notification
	NextCursor    string
	HasMore       bool
}
type Notification struct{ Id, Type, RefId, ActorId, ActorName, ActorAvatar, PostThumb, Message, ReadAt, CreatedAt string }
type MarkReadRequest struct{ NotificationId string }
type MarkReadResponse struct{ Success bool }

type CreateReportRequest struct{ TargetType, TargetId, Reason, Description string }
type CreateReportResponse struct{ ReportId, Message string }

type ListReportsRequest struct {
	Status, Reason, Cursor string
	Limit                  int32
}
type ListReportsResponse struct {
	Reports    []*Report
	NextCursor string
	HasMore    bool
}
type Report struct{ Id, ReporterId, TargetType, TargetId, Reason, Description, Status, HandledBy, CreatedAt, TargetPreview, ReporterName string }
type ResolveReportRequest struct{ ReportId, Action, Note string }
type ResolveReportResponse struct{ Success bool }
type GetModActionsRequest struct {
	Cursor string
	Limit  int32
}
type GetModActionsResponse struct {
	Actions    []*ModAction
	NextCursor string
	HasMore    bool
}
type ModAction struct{ Id, ModeratorId, ModeratorName, TargetType, TargetId, Action, Note, CreatedAt string }

type ListDuelsRequest struct {
	Status, Cursor string
	Limit          int32
}
type ListDuelsResponse struct {
	Duels      []*Duel
	NextCursor string
	HasMore    bool
}
type Duel struct {
	Id                                     string
	PostA, PostB                           *Post
	StartsAt, EndsAt, Status, WinnerPostId string
	VotesA, VotesB                         int64
	PercentA, PercentB                     float64
	MyVotePostId                           string
}
type VoteDuelRequest struct{ DuelId, PostId string }
type VoteDuelResponse struct {
	Success            bool
	PercentA, PercentB float64
}
type GetDuelResultRequest struct{ DuelId string }
type GetDuelResultResponse struct{ Duel *Duel }
type GetWeeklyChampionsRequest struct{}
type GetWeeklyChampionsResponse struct{ Champions []*Duel }

type ListChallengesRequest struct {
	Status, Cursor string
	Limit          int32
}
type ListChallengesResponse struct {
	Challenges []*Challenge
	NextCursor string
	HasMore    bool
}
type Challenge struct {
	Id, Title, Theme, StartsAt, EndsAt, Status string
	EntryCount                                 int64
}
type GetChallengeRequest struct{ ChallengeId string }
type GetChallengeResponse struct{ Challenge *Challenge }

type ListTemplatesRequest struct {
	Query        string
	OfficialOnly bool
	Cursor       string
	Limit        int32
}
type ListTemplatesResponse struct {
	Templates  []*MemeTemplate
	NextCursor string
	HasMore    bool
}
type MemeTemplate struct {
	Id, Name, ImageUrl string
	Width, Height      int32
	TextBoxesJson      string
	IsOfficial         bool
	Status, CreatedAt  string
}
type SubmitTemplateRequest struct {
	Name, ImageKey string
	Width, Height  int32
	TextBoxesJson  string
}
type SubmitTemplateResponse struct{ Template *MemeTemplate }

type ReactPostRequest struct {
	PostId   string
	Reaction int32
}
type ReactPostResponse struct {
	Success    bool
	MyReaction int32
}
type GetReactionsRequest struct{ PostId string }
type GetReactionsResponse struct {
	Counts      *ReactionCounts
	MyReaction  int32
	TopReaction string
	TopPercent  float64
}
type ReactionCounts struct{ Ngakak, Relate, Anjay, Sedih, Marah, Total int64 }

type SendTipRequest struct {
	PostId      string
	AmountCoins int64
}
type SendTipResponse struct {
	Success       bool
	SenderBalance int64
	TransactionId string
}
type ListTipsRequest struct {
	UserId, Type, Cursor string
	Limit                int32
}
type ListTipsResponse struct {
	Tips       []*Tip
	NextCursor string
	HasMore    bool
}
type Tip struct {
	Id, SenderId, SenderName, ReceiverId, ReceiverName, PostId string
	AmountCoins                                                int64
	CreatedAt                                                  string
}

type ListPackagesRequest struct{}
type ListPackagesResponse struct{ Packages []*CoinPackage }
type CoinPackage struct {
	Id              string
	Coins, PriceIdr int64
	PlaySku         string
	IsActive        bool
}
type ConfirmTopUpRequest struct{ PackageId, PurchaseToken, PlaySku string }
type ConfirmTopUpResponse struct {
	Success       bool
	NewBalance    int64
	TransactionId string
}
type GetBalanceRequest struct{}
type GetBalanceResponse struct{ Balance int64 }
type ListTransactionsRequest struct {
	Cursor string
	Limit  int32
}
type ListTransactionsResponse struct {
	Transactions []*CoinTransaction
	NextCursor   string
	HasMore      bool
}
type CoinTransaction struct {
	Id, Type                  string
	AmountCoins, BalanceAfter int64
	RefId, CreatedAt          string
}

// Service interfaces for remaining
type SectionServiceServer interface {
	ListSections(context.Context, *ListSectionsRequest) (*ListSectionsResponse, error)
	GetSection(context.Context, *GetSectionRequest) (*GetSectionResponse, error)
	FollowSection(context.Context, *FollowSectionRequest) (*FollowSectionResponse, error)
}
type UserServiceServer interface {
	GetProfile(context.Context, *GetProfileRequest) (*GetProfileResponse, error)
}
type SearchServiceServer interface {
	Search(context.Context, *SearchRequest) (*SearchResponse, error)
}
type NotificationServiceServer interface {
	ListNotifications(context.Context, *ListNotificationsRequest) (*ListNotificationsResponse, error)
	MarkRead(context.Context, *MarkReadRequest) (*MarkReadResponse, error)
}
type ReportServiceServer interface {
	CreateReport(context.Context, *CreateReportRequest) (*CreateReportResponse, error)
}
type ModerationServiceServer interface {
	ListReports(context.Context, *ListReportsRequest) (*ListReportsResponse, error)
	ResolveReport(context.Context, *ResolveReportRequest) (*ResolveReportResponse, error)
}
type DuelServiceServer interface {
	ListDuels(context.Context, *ListDuelsRequest) (*ListDuelsResponse, error)
	VoteDuel(context.Context, *VoteDuelRequest) (*VoteDuelResponse, error)
}
type ChallengeServiceServer interface {
	ListChallenges(context.Context, *ListChallengesRequest) (*ListChallengesResponse, error)
}
type TemplateServiceServer interface {
	ListTemplates(context.Context, *ListTemplatesRequest) (*ListTemplatesResponse, error)
}
type ReactionServiceServer interface {
	ReactPost(context.Context, *ReactPostRequest) (*ReactPostResponse, error)
	GetReactions(context.Context, *GetReactionsRequest) (*GetReactionsResponse, error)
}
type TipServiceServer interface {
	SendTip(context.Context, *SendTipRequest) (*SendTipResponse, error)
}
type CoinServiceServer interface {
	ListPackages(context.Context, *ListPackagesRequest) (*ListPackagesResponse, error)
	ConfirmTopUp(context.Context, *ConfirmTopUpRequest) (*ConfirmTopUpResponse, error)
	GetBalance(context.Context, *GetBalanceRequest) (*GetBalanceResponse, error)
}
