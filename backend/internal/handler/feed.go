package handler

import (
	"context"
	"github.com/memekita/backend/internal/feed"
	gen "github.com/memekita/backend/proto/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FeedHandler struct{ uc *feed.Usecase }

func NewFeedHandler(uc *feed.Usecase) *FeedHandler { return &FeedHandler{uc: uc} }

func (h *FeedHandler) GetFeed(ctx context.Context, req *gen.GetFeedRequest) (*gen.GetFeedResponse, error) {
	if req.Tab != "hot" && req.Tab != "fresh" && req.Tab != "trending" && req.Tab != "" {
		return nil, status.Error(codes.InvalidArgument, "tab harus hot, fresh, atau trending")
	}
	tab := req.Tab
	if tab == "" {
		tab = "hot"
	}
	posts, next, hasMore, err := h.uc.GetFeed(ctx, tab, req.SectionSlug, req.Cursor, int(req.Limit), req.IncludeNsfw)
	if err != nil {
		return nil, status.Error(codes.Internal, "gagal memuat feed")
	}
	var out []*gen.Post
	for _, p := range posts {
		mediaURL := ""
		if p.MediaURL != nil {
			mediaURL = *p.MediaURL
		}
		thumbURL := ""
		if p.ThumbURL != nil {
			thumbURL = *p.ThumbURL
		}
		out = append(out, &gen.Post{Id: p.ID.String(), UserId: p.UserID.String(), Title: p.Title, MediaType: p.MediaType, MediaUrl: mediaURL, ThumbUrl: thumbURL, Status: p.Status, IsNsfw: p.IsNsfw, Upvotes: int64(p.Upvotes), Downvotes: int64(p.Downvotes), HotScore: p.HotScore})
	}
	return &gen.GetFeedResponse{Posts: out, NextCursor: next, HasMore: hasMore}, nil
}
