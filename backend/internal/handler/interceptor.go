package handler

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"github.com/memekita/backend/internal/config"
)

var publicMethods = map[string]bool{
	"/memekita.v1.AuthService/Register": true,
	"/memekita.v1.AuthService/Login": true,
	"/memekita.v1.AuthService/LoginWithGoogle": true,
	"/memekita.v1.AuthService/RefreshToken": true,
	"/memekita.v1.FeedService/GetFeed": true,
	"/memekita.v1.PostService/GetPost": true,
	"/memekita.v1.CommentService/ListComments": true,
	"/memekita.v1.SectionService/ListSections": true,
	"/memekita.v1.UserService/GetProfile": true,
	"/memekita.v1.SearchService/Search": true,
	"/memekita.v1.ChallengeService/ListChallenges": true,
	"/memekita.v1.TemplateService/ListTemplates": true,
	"/memekita.v1.CoinService/ListPackages": true,
}

func AuthInterceptor(cfg config.Config) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if publicMethods[info.FullMethod] {
			return handler(ctx, req)
		}
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "perlu login")
		}
		authVals := md.Get("authorization")
		if len(authVals) == 0 { authVals = md.Get("Authorization") }
		if len(authVals) == 0 { return nil, status.Error(codes.Unauthenticated, "perlu login") }
		token := strings.TrimPrefix(authVals[0], "Bearer ")
		token = strings.TrimSpace(token)
		if token == "" { return nil, status.Error(codes.Unauthenticated, "token tidak valid") }
		parsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) { return []byte(cfg.JWTSecret), nil })
		if err != nil || !parsed.Valid {
			return nil, status.Error(codes.Unauthenticated, "token tidak valid")
		}
		claims, ok := parsed.Claims.(jwt.MapClaims)
		if !ok { return nil, status.Error(codes.Unauthenticated, "token tidak valid") }
		ctx = context.WithValue(ctx, "claims", claims)
		return handler(ctx, req)
	}
}

type bucket struct { count int; start time.Time }
var (
	rlMu sync.Mutex
	rlBuckets = map[string]*bucket{}
)

func RateLimitInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		claims, _ := ctx.Value("claims").(jwt.MapClaims)
		user := "anon"
		if sub, ok := claims["sub"].(string); ok && sub != "" { user = sub }
		key := user + ":" + info.FullMethod
		limit := 30
		window := 60 * time.Second
		if info.FullMethod == "/memekita.v1.VoteService/VotePost" { limit = 10 }
		if info.FullMethod == "/memekita.v1.CommentService/CreateComment" { limit = 5 }
		now := time.Now()
		rlMu.Lock()
		b, ok := rlBuckets[key]
		if !ok || now.Sub(b.start) > window {
			rlBuckets[key] = &bucket{count: 1, start: now}
			rlMu.Unlock()
			return handler(ctx, req)
		}
		if b.count >= limit {
			rlMu.Unlock()
			return nil, status.Error(codes.ResourceExhausted, "Terlalu cepat, coba lagi sebentar")
		}
		b.count++
		rlMu.Unlock()
		return handler(ctx, req)
	}
}
