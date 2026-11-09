package handler

import (
	"context"
	"github.com/google/uuid"
	"github.com/memekita/backend/internal/vote"
	gen "github.com/memekita/backend/proto/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type VoteHandler struct{ uc *vote.Usecase }

func NewVoteHandler(uc *vote.Usecase) *VoteHandler { return &VoteHandler{uc: uc} }

func (h *VoteHandler) VotePost(ctx context.Context, req *gen.VotePostRequest) (*gen.VotePostResponse, error) {
	claims, ok := ctx.Value("claims").(map[string]interface{})
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "perlu login")
	}
	sub, _ := claims["sub"].(string)
	uid, _ := uuid.Parse(sub)
	postID, err := uuid.Parse(req.PostId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "post_id tidak valid")
	}
	if req.Value != 1 && req.Value != -1 && req.Value != 0 {
		return nil, status.Error(codes.InvalidArgument, "value harus 1, -1, atau 0")
	}
	up, down, my, err := h.uc.VotePost(ctx, uid, postID, int(req.Value))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &gen.VotePostResponse{Upvotes: int64(up), Downvotes: int64(down), MyVote: int32(my)}, nil
}
