# Arsitektur MemeKita

## Ringkasan

Monorepo `proto/`, `backend/`, `android/`, `admin/`, `docs/`. Transport gRPC/HTTP2 + TLS, JWT via metadata + interceptor, media via presigned R2, worker River di PostgreSQL, feed hot cache in-memory TTL, vote transaksi GORM.

```
Android (Kotlin) ──gRPC/TLS──> Go (gRPC) ──GORM──> PostgreSQL
      │                              │  River worker
      │                              ├─ thumbnail & transcode
      │                              ├─ hitung hot_score & trending
      │                              ├─ cron tantangan harian
      │                              └─ kirim notifikasi (FCM)
      └─ unggah langsung ──URL bertanda tangan──> Cloudflare R2
Panel Admin (Vite+React) ──> antrean laporan ──> PostgreSQL
CDN ──> media
```

## Keputusan Desain (lihat `docs/adr/`)

- Kontrak sekali di `.proto`, stub Go/Kotlin, `buf breaking`
- Upload presigned, bukan stream gRPC
- River di Postgres, tanpa Redis
- Feed cursor pagination, Hot cache in-memory TTL, Fresh/Trending dari DB
- Vote + hot_score satu transaksi GORM

## Lapisan

### Android (MVVM + Clean, modular)

`app → feature:feed, feature:duel ... → core`
- Presentation: Fragment + ViewModel (StateFlow, UDF)
- Domain: UseCase murni, interface Repository
- Data: Repository impl, gRPC DataSource, Room/cache
- Hilt DI

### Backend Go (Clean + GORM + Ports & Adapters)

```
handler gRPC (tipis) → usecase → repository (interface) → models (GORM)
```

- Usecase adalah port (logika bisnis), adapter (GORM repo, gRPC handler, River, FCM, R2) di batas luar
- Model GORM `internal/models` terpisah dari entity `internal/domain` (tanpa tag gorm di entity)
- `main.go` hanya wiring (inject db, cfg, usecase, handler)
- DILARANG AutoMigrate di production; migrasi hanya via `golang-migrate`
- Semua query `WithContext(ctx)`, `Updates` bukan `Save`, cek `ErrRecordNotFound`, transaksi untuk vote+counter

### Admin (Feature-Sliced)

`app/ → features/ → entities/ → shared/`

## Skema (21 tabel, `migrations/000001_create_all.up.sql`)

`users, sections, challenges, posts, post_tags, votes, comments, comment_votes, follows, notifications, reports, mod_actions, meme_templates, duels, duel_votes, challenge_entries, tips, user_coins, coin_packages, coin_transactions, post_reactions` + `bookmarks, audit_logs, app_settings`

Indeks: `posts(status, inserted_at)`, `posts(section_id, hot_score)`, `votes(post_id)`.

## Feed

- Fresh: `published` order `created_at DESC`
- Hot: `hot_score DESC` (`(up-down)/(age+2)^1.5`)
- Trending: jendela pendek
- NSFW: `is_nsfw` + setting user

## Sawer Koin

`user_coins`, `coin_transactions(type: topup, tip_sent, tip_received)`, `tips`. Top-up verifikasi Play token → `balance_after`. Kirim atomik, saldo tidak minus, ledger selalu konsisten.

## Moderasi

Laporan `pending` → `mod_actions` + `audit_logs` → post/comment `hidden/removed`, blokir user, notifikasi penulis.

## Keamanan

Argon2, TLS, validasi handler, rate limit interceptor (in-memory bucket), presigned 900s, Keystore.

## CI

`buf lint/breaking`, `gofmt -l`, `go vet`, `golangci-lint` (errcheck, gosec, govet, staticcheck, gocritic), `go test -race -p 1`, `gradlew`, `vite build`, cek stub.

## Mesin Dev 2GB

- Satu pekerjaan berat sekaligus
- Gradle `Xmx1024m`, `workers.max=1`, satu ABI
- Docker `postgres 256m, minio 128m, backend 256m`, River 1-2 worker
- Go test `-p 1 -race` per package
- Node `max-old-space-size=512`
