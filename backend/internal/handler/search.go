package handler

import (
	"context"
	"github.com/memekita/backend/internal/search"
	gen "github.com/memekita/backend/proto/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SearchHandler struct{ uc *search.Usecase }

func NewSearchHandler(uc *search.Usecase) *SearchHandler { return &SearchHandler{uc: uc} }

func (h *SearchHandler) Search(ctx context.Context, req *gen.SearchRequest) (*gen.SearchResponse, error) {
	if req.Query == "" {
		return nil, status.Error(codes.InvalidArgument, "query wajib")
	}
	res, err := h.uc.Search(ctx, req.Query, req.Type, int(req.Limit))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if _, ok := res["posts"]; ok {
		// posts available in result map
	}
	return &gen.SearchResponse{}, nil
}
