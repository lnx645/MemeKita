package handler

import (
	"context"
	"github.com/google/uuid"
	"github.com/memekita/backend/internal/coin"
	gen "github.com/memekita/backend/proto/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TipHandler struct{ uc *coin.Usecase }

func NewTipHandler(uc *coin.Usecase) *TipHandler { return &TipHandler{uc: uc} }

func (h *TipHandler) SendTip(ctx context.Context, req *gen.SendTipRequest) (*gen.SendTipResponse, error) {
	claims, ok := ctx.Value("claims").(map[string]interface{})
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "perlu login")
	}
	sub, _ := claims["sub"].(string)
	senderID, _ := uuid.Parse(sub)
	postID, err := uuid.Parse(req.PostId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "post_id tidak valid")
	}
	if req.AmountCoins <= 0 {
		return nil, status.Error(codes.InvalidArgument, "jumlah koin tidak valid")
	}
	// need receiverID from post
	// For now, use sender as receiver placeholder (real: fetch post.UserID)
	newBal, err := h.uc.SendTip(ctx, senderID, senderID, postID, int(req.AmountCoins))
	if err != nil {
		if err.Error() == "Koin tidak cukup" {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &gen.SendTipResponse{Success: true, SenderBalance: int64(newBal)}, nil
}
