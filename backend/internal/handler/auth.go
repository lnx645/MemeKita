package handler

import (
	"context"
	"github.com/google/uuid"
	"github.com/memekita/backend/internal/auth"
	gen "github.com/memekita/backend/proto/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	uc *auth.Usecase
}

func NewAuthHandler(uc *auth.Usecase) *AuthHandler { return &AuthHandler{uc: uc} }

func (h *AuthHandler) Register(ctx context.Context, req *gen.RegisterRequest) (*gen.AuthResponse, error) {
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "username, email, password wajib")
	}
	if len(req.Username) < 3 {
		return nil, status.Error(codes.InvalidArgument, "username minimal 3")
	}
	user, access, refresh, err := h.uc.Register(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &gen.AuthResponse{AccessToken: access, RefreshToken: refresh, ExpiresIn: 3600, User: &gen.UserInfo{Id: user.ID.String(), Username: user.Username, Email: user.Email, Role: user.Role}}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *gen.LoginRequest) (*gen.AuthResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email dan password wajib")
	}
	user, access, refresh, err := h.uc.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "email atau password salah")
	}
	return &gen.AuthResponse{AccessToken: access, RefreshToken: refresh, ExpiresIn: 3600, User: &gen.UserInfo{Id: user.ID.String(), Username: user.Username, Email: user.Email, Role: user.Role}}, nil
}

func (h *AuthHandler) LoginWithGoogle(ctx context.Context, req *gen.LoginWithGoogleRequest) (*gen.AuthResponse, error) {
	if req.IdToken == "" {
		return nil, status.Error(codes.InvalidArgument, "id_token wajib")
	}
	// In real, verify via Google certs; here we parse unverified for dev
	// For now, reject
	return nil, status.Error(codes.Unimplemented, "Google Auth belum dikonfigurasi")
}

func (h *AuthHandler) RefreshToken(ctx context.Context, req *gen.RefreshTokenRequest) (*gen.AuthResponse, error) {
	if req.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh_token wajib")
	}
	claims, err := h.uc.VerifyAccessToken(req.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "token tidak valid")
	}
	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		return nil, status.Error(codes.Unauthenticated, "token tidak valid")
	}
	uid, err := uuid.Parse(sub)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "token tidak valid")
	}
	user, err := h.uc.FindByID(ctx, uid)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "token tidak valid")
	}
	if user.Status == "blocked" {
		return nil, status.Error(codes.PermissionDenied, "akun diblokir")
	}
	return nil, status.Error(codes.Unimplemented, "refresh belum diimplementasi")
}
