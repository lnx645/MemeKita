package handler

import (
	"context"
	"github.com/google/uuid"
	"github.com/memekita/backend/internal/duel"
	gen "github.com/memekita/backend/proto/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DuelHandler struct{ uc *duel.Usecase }

func NewDuelHandler(uc *duel.Usecase) *DuelHandler { return &DuelHandler{uc: uc} }

func (h *DuelHandler) ListDuels(ctx context.Context, req *gen.ListDuelsRequest) (*gen.ListDuelsResponse, error) {
	return &gen.ListDuelsResponse{}, nil
}
func (h *DuelHandler) VoteDuel(ctx context.Context, req *gen.VoteDuelRequest) (*gen.VoteDuelResponse, error) {
	claims, ok := ctx.Value("claims").(map[string]interface{})
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "perlu login")
	}
	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		return nil, status.Error(codes.Unauthenticated, "perlu login")
	}
	uid, errParse := uuid.Parse(sub)
	if errParse != nil {
		return nil, status.Error(codes.Unauthenticated, "token tidak valid")
	}
	duelID, err := uuid.Parse(req.DuelId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "duel_id tidak valid")
	}
	postID, err := uuid.Parse(req.PostId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "post_id tidak valid")
	}
	if err := h.uc.Vote(ctx, uid, duelID, postID); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &gen.VoteDuelResponse{Success: true}, nil
}
