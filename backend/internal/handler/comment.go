package handler

import (
	"context"
	"github.com/google/uuid"
	"github.com/memekita/backend/internal/comment"
	gen "github.com/memekita/backend/proto/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CommentHandler struct{ uc *comment.Usecase }

func NewCommentHandler(uc *comment.Usecase) *CommentHandler { return &CommentHandler{uc: uc} }

func (h *CommentHandler) ListComments(ctx context.Context, req *gen.ListCommentsRequest) (*gen.ListCommentsResponse, error) {
	if req.PostId == "" {
		return nil, status.Error(codes.InvalidArgument, "post_id wajib")
	}
	pid, err := uuid.Parse(req.PostId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "post_id tidak valid")
	}
	comments, err := h.uc.List(ctx, pid, req.Sort)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var out []*gen.Comment
	for _, c := range comments {
		out = append(out, &gen.Comment{Id: c.ID.String(), PostId: c.PostID.String(), UserId: c.UserID.String(), Body: c.Body, Status: c.Status})
	}
	return &gen.ListCommentsResponse{Comments: out}, nil
}

func (h *CommentHandler) CreateComment(ctx context.Context, req *gen.CreateCommentRequest) (*gen.CreateCommentResponse, error) {
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
	var parentID *uuid.UUID
	if req.ParentId != "" {
		pid, err := uuid.Parse(req.ParentId)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "parent_id tidak valid")
		}
		parentID = &pid
	}
	if req.Body == "" && req.MediaUrl == "" {
		return nil, status.Error(codes.InvalidArgument, "komentar tidak boleh kosong")
	}
	c, err := h.uc.Create(ctx, postID, uid, parentID, req.Body, &req.MediaUrl)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &gen.CreateCommentResponse{Comment: &gen.Comment{Id: c.ID.String(), Body: c.Body}}, nil
}
