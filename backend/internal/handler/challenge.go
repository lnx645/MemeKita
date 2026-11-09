package handler

import (
	"context"
	"github.com/memekita/backend/internal/challenge"
	gen "github.com/memekita/backend/proto/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ChallengeHandler struct{ uc *challenge.Usecase }

func NewChallengeHandler(uc *challenge.Usecase) *ChallengeHandler { return &ChallengeHandler{uc: uc} }

func (h *ChallengeHandler) ListChallenges(ctx context.Context, req *gen.ListChallengesRequest) (*gen.ListChallengesResponse, error) {
	c, err := h.uc.Current(ctx)
	if err != nil {
		return &gen.ListChallengesResponse{}, nil
	}
	return &gen.ListChallengesResponse{Challenges: []*gen.Challenge{{Id: c.ID.String(), Title: c.Title, Theme: c.Theme, Status: c.Status}}}, nil
}
func (h *ChallengeHandler) GetChallenge(ctx context.Context, req *gen.GetChallengeRequest) (*gen.GetChallengeResponse, error) {
	if req.ChallengeId == "" {
		return nil, status.Error(codes.InvalidArgument, "challenge_id wajib")
	}
	return &gen.GetChallengeResponse{}, nil
}
