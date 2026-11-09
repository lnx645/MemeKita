package handler

import (
	"context"
	"github.com/google/uuid"
	"github.com/memekita/backend/internal/section"
	gen "github.com/memekita/backend/proto/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SectionHandler struct{ uc *section.Usecase }

func NewSectionHandler(uc *section.Usecase) *SectionHandler { return &SectionHandler{uc: uc} }

func (h *SectionHandler) ListSections(ctx context.Context, req *gen.ListSectionsRequest) (*gen.ListSectionsResponse, error) {
	secs, err := h.uc.List(ctx, req.IncludeRegional)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var out []*gen.SectionInfo
	for _, s := range secs {
		desc := ""
		if s.Description != nil {
			desc = *s.Description
		}
		icon := ""
		if s.Icon != nil {
			icon = *s.Icon
		}
		out = append(out, &gen.SectionInfo{Id: s.ID.String(), Slug: s.Slug, Name: s.Name, Description: desc, IsNsfw: s.IsNsfw, IsRegional: s.IsRegional, Icon: icon})
	}
	return &gen.ListSectionsResponse{Sections: out}, nil
}

func (h *SectionHandler) GetSection(ctx context.Context, req *gen.GetSectionRequest) (*gen.GetSectionResponse, error) {
	if req.Slug == "" {
		return nil, status.Error(codes.InvalidArgument, "slug wajib")
	}
	s, err := h.uc.GetBySlug(ctx, req.Slug)
	if err != nil {
		return nil, status.Error(codes.NotFound, "section tidak ditemukan")
	}
	desc := ""
	if s.Description != nil {
		desc = *s.Description
	}
	return &gen.GetSectionResponse{Section: &gen.SectionInfo{Id: s.ID.String(), Slug: s.Slug, Name: s.Name, Description: desc, IsRegional: s.IsRegional}}, nil
}

func (h *SectionHandler) FollowSection(ctx context.Context, req *gen.FollowSectionRequest) (*gen.FollowSectionResponse, error) {
	claims, ok := ctx.Value("claims").(map[string]interface{})
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "perlu login")
	}
	sub, _ := claims["sub"].(string)
	uid, _ := uuid.Parse(sub)
	sid, err := uuid.Parse(req.SectionId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "section_id tidak valid")
	}
	if err := h.uc.Follow(ctx, uid, sid); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &gen.FollowSectionResponse{Success: true}, nil
}
