package config

import (
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL  string
	GRPCPort     string
	HTTPPort     string
	JWTSecret    string
	R2Endpoint   string
	R2Bucket     string
	R2AccessKey  string
	R2SecretKey  string
	CDNURL       string
	FCMProjectID string
	FCMCredPath  string
	AccessTTL    int
	RefreshTTL   int
	MaxImageSize int64
	MaxGifSize   int64
	MaxVideoSize int64
}

func Load() Config {
	return Config{
		DatabaseURL:  env("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/memekita_dev?sslmode=disable"),
		GRPCPort:     env("GRPC_PORT", "50051"),
		HTTPPort:     env("HTTP_PORT", "4000"),
		JWTSecret:    loadJWTSecret(),
		R2Endpoint:   env("R2_ENDPOINT", "http://localhost:9000"),
		R2Bucket:     env("R2_BUCKET", "memekita-media"),
		R2AccessKey:  env("R2_ACCESS_KEY", "minioadmin"),
		R2SecretKey:  env("R2_SECRET_KEY", "minioadmin123"),
		CDNURL:       env("CDN_URL", "http://localhost:9000/memekita-media"),
		FCMProjectID: env("FCM_PROJECT_ID", "ceoshapp-e7ada"),
		FCMCredPath:  env("FCM_CREDENTIALS_PATH", "/home/demo/Downloads/takmungkin/firebase_fcm"),
		AccessTTL:    envInt("JWT_ACCESS_TTL", 3600),
		RefreshTTL:   envInt("JWT_REFRESH_TTL", 86400*7),
		MaxImageSize: 10 * 1024 * 1024,
		MaxGifSize:   15 * 1024 * 1024,
		MaxVideoSize: 50 * 1024 * 1024,
	}
}

func loadJWTSecret() string {
	if s := os.Getenv("JWT_SECRET"); s != "" {
		return s
	}
	if b, err := os.ReadFile("/home/demo/Downloads/takmungkin/firebase_fcm/jwt_secret"); err == nil {
		if s := string(b); s != "" {
			return s
		}
	}
	return "dev_secret_key_change_in_prod_32chars!!"
}

func env(k, fallback string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fallback
}

func envInt(k string, fallback int) int {
	if v := os.Getenv(k); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}
