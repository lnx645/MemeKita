package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/memekita/backend/internal/config"
	"github.com/memekita/backend/internal/db"
	"github.com/memekita/backend/internal/handler"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	gormDB, err := db.NewGORM(ctx, cfg)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	// wiring — only constructors, no business logic in main

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(handler.AuthInterceptor(cfg), handler.RateLimitInterceptor()),
	)
	// Register gRPC services here (generated from proto/)
	// Example: pb.RegisterAuthServiceServer(s, &AuthHandler{uc: auth.NewUsecase(repo, gormDB, cfg)})
	log.Printf("gRPC listening on :%s", cfg.GRPCPort)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("grpc serve: %v", err)
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{"status":"ok"}`)); err != nil {
			log.Printf("write error: %v", err)
		}
	})
	mux.HandleFunc("/api/dashboard", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{"users":124,"posts":87,"reports":5}`)); err != nil {
			log.Printf("write error: %v", err)
		}
	})
	if gormDB != nil {
		log.Printf("db ready: postgres connected")
	}
	log.Printf("HTTP listening on :%s", cfg.HTTPPort)
	if err := http.ListenAndServe(":"+cfg.HTTPPort, mux); err != nil {
		log.Fatalf("http: %v", err)
	}
}
