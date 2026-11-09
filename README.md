# MemeKita — Aplikasi Meme ala 9GAG (Android + Go)

Monorepo untuk komunitas meme dengan feed vertikal, vote, komentar bersarang, unggah media, section, profil, moderasi, plus fitur unik: Meme Duel, generator meme, kanal daerah, tantangan harian, sawer koin, reaksi Indonesia, balas bergambar, share WA.

## Struktur

```
proto/          # kontrak gRPC (Protobuf, buf)
backend/        # Go 1.22 + gRPC + GORM (postgres) + River + AWS S3 (R2)
android/        # Kotlin native, ViewBinding, MVVM, Hilt, Paging 3, Coil, ExoPlayer, gRPC
admin/          # Vite + React + TypeScript + Astryx
docs/           # arsitektur, asumsi, ADR
```

## Prasyarat

- Go 1.22+, Node 20, JDK 17, Android SDK 34
- PostgreSQL 15, MinIO (dev) / Cloudflare R2 (prod)
- Buf CLI (`buf lint`, `buf breaking`), golang-migrate, golangci-lint

## Menjalankan Backend Lokal

```bash
docker compose up -d postgres minio   # 256m + 128m (hemat RAM 2GB)
cd backend
go mod download
# migrasi
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/memekita_dev?sslmode=disable" up
go run ./cmd/seed   # 10 user, 50 post orisinal, duels, coin packages, tantangan, templates
go run ./cmd/server # gRPC :50051, HTTP :4000 (health + /api/dashboard)
# atau via Docker
docker compose up backend
```

Kredensial R2 dibaca dari `/home/demo/Downloads/cloudflare` (jangan commit). FCM/Google Auth dari `/home/demo/Downloads/takmungkin` (service account `ceoshapp-e7ada`). Di prod, set via env `R2_*` / `DATABASE_URL` (lihat `internal/config/config.go`). GORM hanya untuk query — migrasi via `golang-migrate` (DILARANG AutoMigrate di prod).

## Membangun Android

```bash
cd android
./gradlew assembleDebug -Pabi=armeabi-v7a   # dev: satu ABI, no R8, incremental
# APK: android/app/build/outputs/apk/debug/app-armeabi-v7a-debug.apk
# Rilis AAB: ./gradlew bundleRelease (hanya di CI, bukan mesin 2GB)
```

`gradle.properties` sudah dioptimasi untuk 2GB: `Xmx1024m`, `workers.max=1`, `parallel=false`, `caching=true`.

Login Google: `google-services.json` dari `/home/demo/Downloads/takmungkin/firebase_fcm/` sudah mendukung package `com.meme_kita.mobile` (project `ceoshapp-e7ada`). FCM memakai service account yang sama.

Font Manrope sudah dibundel di `android/app/src/main/res/font/` (dari `/home/demo/Downloads/imo_fonts`).

## Menjalankan Panel Admin

```bash
cd admin
npm install
npm run dev    # http://localhost:5173  (NODE_OPTIONS=--max-old-space-size=512)
npm run build  # preview produksi
```

Admin memakai komponen **Astryx** (https://astryx.atmeta.com/components): App Shell + Side Nav, Table, Dialog, Field, Button, Badge, Tabs, Pagination, Card. Warna MemeKita `#6C4CF1` sebagai aksen.

Login admin: `admin@memekita.id / admin12345` (seed). Semua tindakan tercatat di `audit_logs`.

## Kontrak API

```bash
buf lint proto
buf breaking proto --against '.git#branch=main,subdir=proto'
# generate Go + Kotlin
buf generate proto --template proto/buf.gen.yaml
# Go stubs: backend/proto/gen, Kotlin: android/app/src/main/java/com/meme_kita/mobile/proto
```

Stub Go dan Kotlin dihasilkan dari `proto/` — jangan edit manual.

## Tes

```bash
# backend per package, -race (hemat RAM, jangan semua sekaligus)
go test ./internal/auth -run TestHash -count=1
go test ./internal/post -run TestHotScore -count=1 -race
go test ./internal/media -run TestValidate -count=1
go test ./... -p 1 -race   # full di CI

# Android per modul
./gradlew :app:testDebugUnitTest
./gradlew :app:connectedDebugAndroidTest   # Espresso, butuh HP fisik armeabi-v7a
```

CI (`.github/workflows/ci.yml`) menjalankan: `buf lint/breaking`, `gofmt`, `go vet`, `golangci-lint` (errcheck, gosec, govet, staticcheck), `go test -race`, `gradlew lint + test`, `vite build`, cek stub (`grep TODO/notImplemented`).

## Fitur Unik

- **Meme Duel**: duel 24 jam, satu vote per user, bracket mingguan, lencana Juara Duel
- **Generator**: pilih template → teks atas/bawah → render di perangkat → unggah
- **Kanal Daerah**: section Sunda/Jawa/Betawi (is_regional), moderator per kanal
- **Tantangan Harian**: River job tengah malam, ganti tema otomatis
- **Sawer Koin**: ledger `user_coins` + `coin_transactions`, top-up Play Billing, leaderboard mingguan
- **Reaksi Indonesia**: enum `ngakak, relate, anjay, sedih, marah`, persentase top, tidak pengaruhi hot_score
- **Balas Bergambar**: comment `media_url` via alur unggah sama
- **Share WA**: `Intent.ACTION_SEND` + tautan `https://memekita.id/post/{id}`

## Clean Architecture

- **Android**: `ui → ViewModel → UseCase → Repository → DataSource (gRPC, Room)`; Hilt DI; `Result` sealed
- **Backend Go**: `handler gRPC (tipis) → usecase → repository (interface GORM) → models (GORM tag)`; `main.go` hanya wiring, tidak ada logika; River worker tipis
- **Admin**: Feature-Sliced, Astryx, service terpusat

Skor Hot: `hot_score = (upvotes - downvotes) / (umur_jam + 2)^1.5`

Keamanan: Argon2, TLS, validasi server, rate limit interceptor gRPC, batas ukuran, verifikasi magic bytes, presigned 900s, EncryptedSharedPreferences.

Lihat `docs/ARCHITECTURE.md`, `docs/ASSUMPTIONS.md`, `docs/adr/`.
