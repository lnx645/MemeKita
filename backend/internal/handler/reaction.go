package handler

import (
	"context"
	"github.com/google/uuid"
	"github.com/memekita/backend/internal/reaction"
	gen "github.com/memekita/backend/proto/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ReactionHandler struct{ uc *reaction.Usecase }

func NewReactionHandler(uc *reaction.Usecase) *ReactionHandler { return &ReactionHandler{uc: uc} }

func (h *ReactionHandler) ReactPost(ctx context.Context, req *gen.ReactPostRequest) (*gen.ReactPostResponse, error) {
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
	reactions := map[int32]string{1: "ngakak", 2: "relate", 3: "anjay", 4: "sedih", 5: "marah"}
	r, ok := reactions[req.Reaction]
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "reaksi tidak valid")
	}
	if err := h.uc.React(ctx, uid, postID, r); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &gen.ReactPostResponse{Success: true, MyReaction: req.Reaction}, nil
}

func (h *ReactionHandler) GetReactions(ctx context.Context, req *gen.GetReactionsRequest) (*gen.GetReactionsResponse, error) {
	postID, err := uuid.Parse(req.PostId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "post_id tidak valid")
	}
	counts, top, percent := h.uc.Counts(ctx, postID)
	return &gen.GetReactionsResponse{Counts: &gen.ReactionCounts{Ngakak: int64(counts["ngakak"]), Relate: int64(counts["relate"]), Anjay: int64(counts["anjay"]), Sedih: int64(counts["sedih"]), Marah: int64(counts["marah"]), Total: int64(counts["ngakak"] + counts["relate"] + counts["anjay"] + counts["sedih"] + counts["marah"])}, TopReaction: top, TopPercent: percent}, nil
}
