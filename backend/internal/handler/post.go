package handler

import (
	"context"
	"github.com/google/uuid"
	"github.com/memekita/backend/internal/config"
	"github.com/memekita/backend/internal/media"
	"github.com/memekita/backend/internal/post"
	gen "github.com/memekita/backend/proto/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PostHandler struct {
	uc  *post.Usecase
	cfg config.Config
}

func NewPostHandler(uc *post.Usecase, cfg config.Config) *PostHandler {
	return &PostHandler{uc: uc, cfg: cfg}
}

func (h *PostHandler) GetPost(ctx context.Context, req *gen.GetPostRequest) (*gen.GetPostResponse, error) {
	if req.PostId == "" {
		return nil, status.Error(codes.InvalidArgument, "post_id wajib")
	}
	id, err := uuid.Parse(req.PostId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "post_id tidak valid")
	}
	p, err := h.uc.Get(ctx, id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "post tidak ditemukan")
	}
	mediaURL := ""
	if p.MediaURL != nil {
		mediaURL = *p.MediaURL
	}
	return &gen.GetPostResponse{Post: &gen.Post{Id: p.ID.String(), Title: p.Title, MediaUrl: mediaURL, Status: p.Status}}, nil
}

func (h *PostHandler) CreateUpload(ctx context.Context, req *gen.CreateUploadRequest) (*gen.CreateUploadResponse, error) {
	claims, ok := ctx.Value("claims").(map[string]interface{})
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "perlu login")
	}
	if len(claims) == 0 {
		return nil, status.Error(codes.Unauthenticated, "perlu login")
	}
	mediaType := map[int32]string{1: "image", 2: "gif", 3: "video"}[req.MediaType]
	if mediaType == "" {
		mediaType = "image"
	}
	if err := media.Validate(mediaType, req.FileSize, req.ContentType); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	// Get userID from claims
	userIDStr, _ := claims["sub"].(string)
	uid, _ := uuid.Parse(userIDStr)
	key := media.GenerateKey(uid, req.FileName)
	url := media.PresignedURL(h.cfg.R2Endpoint, h.cfg.R2Bucket, key)
	return &gen.CreateUploadResponse{UploadUrl: url, MediaKey: key, ExpiresIn: 900}, nil
}

func (h *PostHandler) CompletePost(ctx context.Context, req *gen.CompletePostRequest) (*gen.CompletePostResponse, error) {
	claims, ok := ctx.Value("claims").(map[string]interface{})
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "perlu login")
	}
	userIDStr, _ := claims["sub"].(string)
	uid, _ := uuid.Parse(userIDStr)
	if req.Title == "" {
		return nil, status.Error(codes.InvalidArgument, "judul wajib")
	}
	if req.SectionId == "" {
		return nil, status.Error(codes.InvalidArgument, "section wajib")
	}
	sectionID, err := uuid.Parse(req.SectionId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "section tidak valid")
	}
	mediaType := map[int32]string{1: "image", 2: "gif", 3: "video"}[req.MediaType]
	if mediaType == "" {
		mediaType = "image"
	}
	var challengeID *uuid.UUID
	if req.ChallengeId != "" {
		cid, err := uuid.Parse(req.ChallengeId)
		if err == nil {
			challengeID = &cid
		}
	}
	p, err := h.uc.Create(ctx, uid, sectionID, req.Title, mediaType, req.MediaKey, int(req.Width), int(req.Height), req.IsNsfw, req.Tags, challengeID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &gen.CompletePostResponse{Post: &gen.Post{Id: p.ID.String(), Title: p.Title}}, nil
}

func (h *PostHandler) DeletePost(ctx context.Context, req *gen.DeletePostRequest) (*gen.DeletePostResponse, error) {
	if req.PostId == "" {
		return nil, status.Error(codes.InvalidArgument, "post_id wajib")
	}
	return &gen.DeletePostResponse{Success: true}, nil
}
