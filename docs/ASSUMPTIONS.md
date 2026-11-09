# Asumsi & Keputusan Terbuka

1. **Go + GORM** — Spec Go(1) mewajibkan GORM (bukan pgx/sqlc). Model GORM di `internal/models` terpisah dari entity `internal/domain`. Migrasi hanya via `golang-migrate` (`migrations/*.sql`), DILARANG `AutoMigrate` di prod. Semua query `WithContext(ctx)`, `Updates` untuk partial, cek `ErrRecordNotFound`, transaksi untuk vote+counter.

2. **Google Auth / FCM** — Folder `/home/demo/Downloads/takmungkin` (project `ceoshapp-e7ada`, package `com.meme_kita.mobile` di `google-services(3).json`). Fallback ke env `JWT_SECRET`, `DATABASE_URL`, `R2_*`. Tidak commit.

3. **Cloudflare R2** — Folder `/home/demo/Downloads/cloudflare` tidak ada di mesin dev → pakai MinIO lokal (`docker-compose.yml`); prod via env `R2_*`. Presigned via AWS SDK v2.

4. **Buf** — Kontrak `.proto` + `buf.yaml`/`buf.gen.yaml` disediakan; CI `buf lint/breaking`. Stub Go/Kotlin dari sumber sama (tidak edit manual), Go stubs di `backend/proto/gen`, Kotlin di `android/app/src/main/java/com/meme_kita/mobile/proto/`.

5. **Astryx** — Paket tidak di npm publik → fallback lokal `admin/src/components/Astryx.tsx` dengan API identik (App Shell, Table, Dialog, Field, Button, Badge, Tabs, Pagination, Card). Ganti import ke `@astryx/react` bila tersedia. Warna `#6C4CF1` aksen.

6. **Manrope** — Bundel dari `/home/demo/Downloads/imo_fonts` ke `android/app/src/main/res/font/`, fallback `sans-serif`.

7. **Trending** — Awal pakai `hot_score` sama, tuning pangkat nanti (ADR).

8. **Sawer penarikan** — Tidak termasuk fase awal.

9. **Ekspor data & hapus akun** — `ExportData` → signed URL sementara; hapus akun → anonim "Pengguna dihapus".

10. **Rate limit** — Interceptor in-memory bucket (10 vote/menit, 5 komentar/menit). Prod bisa ganti Redis.

11. **Kesalahan Go yang dicegah** — Semua error diperiksa, tidak ada `go func()` tanpa context, tidak ada data race (`go test -race` lolos), `defer Close()` tepat, tanpa N+1 (Preload), tanpa `fmt.Sprintf` SQL, tanpa `Save` untuk partial, tanpa `nil, nil`, tanpa global var `db` (inject), UTC untuk waktu, validasi di handler.

12. **Build 2GB** — Gradle `Xmx1024m`, workers 1, satu ABI; Docker `mem_limit`; Go test `-p 1`; Node 512MB.
