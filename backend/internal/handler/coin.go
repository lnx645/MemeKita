package handler

import (
	"context"
	"github.com/google/uuid"
	"github.com/memekita/backend/internal/coin"
	gen "github.com/memekita/backend/proto/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CoinHandler struct{ uc *coin.Usecase }

func NewCoinHandler(uc *coin.Usecase) *CoinHandler { return &CoinHandler{uc: uc} }

func (h *CoinHandler) ListPackages(ctx context.Context, req *gen.ListPackagesRequest) (*gen.ListPackagesResponse, error) {
	pkgs, err := h.uc.ListPackages(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var out []*gen.CoinPackage
	for _, p := range pkgs {
		out = append(out, &gen.CoinPackage{Id: p.ID.String(), Coins: int64(p.Coins), PriceIdr: int64(p.PriceIdr), PlaySku: p.PlaySku, IsActive: p.IsActive})
	}
	return &gen.ListPackagesResponse{Packages: out}, nil
}

func (h *CoinHandler) ConfirmTopUp(ctx context.Context, req *gen.ConfirmTopUpRequest) (*gen.ConfirmTopUpResponse, error) {
	claims, ok := ctx.Value("claims").(map[string]interface{})
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "perlu login")
	}
	sub, _ := claims["sub"].(string)
	uid, _ := uuid.Parse(sub)
	pkgID, err := uuid.Parse(req.PackageId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "package_id tidak valid")
	}
	if req.PurchaseToken == "" {
		return nil, status.Error(codes.InvalidArgument, "purchase_token wajib")
	}
	newBal, txID, err := h.uc.Topup(ctx, uid, pkgID, req.PurchaseToken)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &gen.ConfirmTopUpResponse{Success: true, NewBalance: int64(newBal), TransactionId: txID}, nil
}

func (h *CoinHandler) GetBalance(ctx context.Context, req *gen.GetBalanceRequest) (*gen.GetBalanceResponse, error) {
	claims, ok := ctx.Value("claims").(map[string]interface{})
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "perlu login")
	}
	sub, _ := claims["sub"].(string)
	uid, _ := uuid.Parse(sub)
	bal := h.uc.GetBalance(ctx, uid)
	return &gen.GetBalanceResponse{Balance: int64(bal)}, nil
}
