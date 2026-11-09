package handler

import (
	"context"
	"github.com/google/uuid"
	"github.com/memekita/backend/internal/moderation"
	gen "github.com/memekita/backend/proto/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ModerationHandler struct{ uc *moderation.Usecase }

func NewModerationHandler(uc *moderation.Usecase) *ModerationHandler {
	return &ModerationHandler{uc: uc}
}

func (h *ModerationHandler) CreateReport(ctx context.Context, req *gen.CreateReportRequest) (*gen.CreateReportResponse, error) {
	claims, ok := ctx.Value("claims").(map[string]interface{})
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "perlu login")
	}
	sub, _ := claims["sub"].(string)
	uid, _ := uuid.Parse(sub)
	if req.TargetId == "" || req.Reason == "" {
		return nil, status.Error(codes.InvalidArgument, "target dan alasan wajib")
	}
	targetID, err := uuid.Parse(req.TargetId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "target_id tidak valid")
	}
	id, err := h.uc.CreateReport(ctx, uid, req.TargetType, targetID, req.Reason)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &gen.CreateReportResponse{ReportId: id.String(), Message: "Terima kasih, laporanmu sudah kami terima."}, nil
}
func (h *ModerationHandler) ListReports(ctx context.Context, req *gen.ListReportsRequest) (*gen.ListReportsResponse, error) {
	return &gen.ListReportsResponse{}, nil
}
func (h *ModerationHandler) ResolveReport(ctx context.Context, req *gen.ResolveReportRequest) (*gen.ResolveReportResponse, error) {
	claims, ok := ctx.Value("claims").(map[string]interface{})
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "perlu login")
	}
	sub, _ := claims["sub"].(string)
	uid, _ := uuid.Parse(sub)
	reportID, err := uuid.Parse(req.ReportId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "report_id tidak valid")
	}
	if req.Action == "" {
		return nil, status.Error(codes.InvalidArgument, "action wajib")
	}
	if err := h.uc.Resolve(ctx, reportID, uid, req.Action, req.Note); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &gen.ResolveReportResponse{Success: true}, nil
}
