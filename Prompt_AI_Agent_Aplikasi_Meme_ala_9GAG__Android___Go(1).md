# Prompt AI Agent: Aplikasi Meme ala 9GAG (Android + Go)

> **Keputusan 2026-11-09: backend memakai Go (bukan Elixir).** Alasan: developer Go lebih mudah dicari di Indonesia dan beban MemeKita (API + feed + vote + job latar) bukan beban jutaan koneksi persisten tempat Elixir unggul. Kontrak .proto dan semua fitur tidak berubah.

Nov 9, 2026

Tempel blok di bawah ini ke AI agent (Claude Code, Cursor, dan sejenisnya) sebagai instruksi awal; bagian berikutnya di dokumen ini adalah lampiran spesifikasinya.

## Prompt utama (siap tempel)

```text
PERAN
Kamu engineer senior Android dan Go. Bangun aplikasi komunitas meme dengan fitur dan alur pengguna yang setara dengan 9GAG: feed vertikal, vote naik/turun, komentar bersarang, unggah gambar/GIF/video, section (kategori), profil, dan moderasi.

TUJUAN
Empat hasil yang berjalan bersama: (1) backend Go yang bisa dijalankan lokal dengan docker compose up, (2) aplikasi Android native yang terhubung ke backend itu, (3) panel admin web (Vite + React) yang mengelola seluruh siklus aplikasi, (4) fitur unik pembeda dari 9GAG (Meme Duel, generator meme, kanal daerah, tantangan harian, sawer kreator, reaksi Indonesia, balas komentar bergambar, share ke Status WA). Proyek ini akan dipublikasikan perusahaan, jadi kode, tes, dan dokumentasi diperlakukan setara produksi.

ATURAN KERAS
1. Samakan fitur dan alur penggunaan dengan 9GAG. Jangan menyalin nama, logo, palet warna, teks, ikon, aset, atau kode milik 9GAG maupun Lahelu. Buat nama, merek, dan tampilan sendiri. Nama kerja: MemeKita (boleh diganti).
2. Kerjakan per fase sesuai dokumen. Selesaikan satu fase penuh (kode, tes, dokumentasi) sebelum lanjut.
3. Setiap fitur punya tes otomatis: Go testing (go test + testify) untuk backend, JUnit dan Espresso untuk Android. Jangan tandai selesai bila tes gagal.
4. Jangan menebak. Keputusan yang memengaruhi arsitektur dan tidak tercantum dicatat di docs/ASSUMPTIONS.md, lalu lanjut.
5. Kontrak API hanya lewat file .proto di folder proto/. Stub Go dan Kotlin dihasilkan dari file itu dan tidak diedit manual. Jalankan buf lint dan buf breaking di CI.
6. Keamanan dasar wajib: hash password argon2, TLS, validasi semua input di server, rate limit lewat interceptor gRPC, batas ukuran unggahan, verifikasi tipe file dari isinya, URL unggah bertanda tangan berumur pendek, token disimpan aman di Android Keystore.
7. Semua konten pengguna bisa dilaporkan, diantrekan, dan diturunkan lewat alur moderasi.
8. Bangun juga semua fitur unik di bagian 'Fitur unik (pembeda dari 9GAG)' — itu pembeda produk ini dari 9GAG, bukan tambahan opsional.
9. Notifikasi memakai FCM; login/register memakai Google Auth (Google Sign-In) selain email/password. Kredensial Firebase dan Google Auth ada di folder /home/demo/Downloads/takmungkin di komputer pemilik produk — baca dari sana, jangan hardcode, jangan commit ke repo.
10. Ikuti sistem desain di bagian 'Sistem desain' (ikon Material Symbols, font Manrope) dan pastikan semua skenario di bagian 'Skenario pengguna' bisa dijalankan. Checklist 'Detail kecil yang wajib ada' harus terpenuhi sebelum rilis.
11. Penyimpanan media memakai Cloudflare R2 (bukan MinIO di produksi). Kredensial R2 ada di folder /home/demo/Downloads/cloudflare di komputer pemilik produk — baca dari sana, jangan hardcode, jangan commit ke repo.
12. Ikuti referensi UI/UX di bagian 'Referensi UI/UX' — pola dan perilaku harus setara 9GAG (feed vertikal, kartu post, vote, tab Hot/Fresh/Trending) dengan pembeda di §Fitur unik dan identitas sendiri (nama, warna, font, ikon). Dilarang menyalin nama, logo, palet warna, teks antarmuka, ikon, atau aset 9GAG/Lahelu.
13. Sawer memakai koin (mata uang dalam aplikasi): pengguna membeli koin lewat top-up, lalu memakai koin untuk menyawer kreator. Panel admin berbasis web (Vite + React) dan mengelola seluruh siklus aplikasi (lihat §Panel Admin) — bukan LiveView. Frontend panel admin memakai komponen Astryx (https://astryx.atmeta.com/components).
14. DILARANG stub: tidak ada fungsi kosong, TODO, `notImplemented()`, mock, atau data palsu yang pura-pura berfungsi. Semua fitur diimplementasi penuh dan teruji; bila satu bagian belum bisa selesai, laporkan sebagai BELUM SELESAI — jangan ditutup dengan stub.
15. Semua kode WAJIB clean architecture agar mudah di-maintenance ke depan (lihat §Clean Architecture & Maintainability).
16. Jangan menulis kode berantakan: gunakan logika terbaik di aplikasi — pilih algoritma dan struktur data yang tepat, tangani edge case, dan tulis kode seolah-olah akan dibaca orang lain (lihat §Clean Architecture & Maintainability, "Logika terbaik").
17. Gunakan arsitektur terbaik dan modern; perkecil spaghetti code — dilarang circular dependency, god class, dan mencampur UI + logika + data dalam satu file; setiap lapisan punya batas yang dijaga tooling (lihat §Clean Architecture & Maintainability, "Arsitektur terbaik & modern").
18. Logo dan identitas merek harus REAL (bukan asal-asalan, bukan placeholder, bukan teks saja): logo berupa gambar maskot gaya emoji lucu yang digambar sebagai vektor; dipakai di splash screen, ikon aplikasi, dan Play Store; harus disetujui pemilik produk sebelum rilis (lihat §Identitas Merek & Logo).
19. Mesin dev hanya 2 GB RAM dan pernah crash/freeze saat develop mobile, test, dan running backend — agent WAJIB pintar memakai memori: jangan jalankan build + emulator + backend + test bersamaan; ikuti optimasi di §Mesin Dev (2 GB RAM).
20. Backend Go WAJIB clean architecture dan memakai GORM; dilarang kesalahan-kesalahan kecil yang sering dibuat AI — error diabaikan, goroutine bocor, data race, tanpa context, resource tidak ditutup, N+1, SQL injection, AutoMigrate di production, state global (lihat §Clean Architecture & Maintainability, "Kesalahan umum AI dalam Go").

STACK
- Android: Kotlin native; package `com.meme_kita.mobile`; layout XML dengan ViewBinding (tanpa Jetpack Compose); Fragment + Navigation Component; MVVM dengan ViewModel dan StateFlow; coroutines; Hilt; RecyclerView + Paging 3; Coil untuk gambar; Media3 ExoPlayer untuk video; grpc-kotlin dan protobuf untuk klien; Firebase Messaging (FCM)
- Backend: Go (versi stabil terbaru, minimal 1.22); server gRPC (grpc-go); PostgreSQL via GORM (driver postgres); River untuk tugas latar (antrean di PostgreSQL); AWS SDK for Go v2 untuk S3; JWT; Argon2
- Penyimpanan media: Cloudflare R2 (S3-compatible; MinIO hanya untuk dev lokal bila perlu); kredensial R2 di /home/demo/Downloads/cloudflare
- Panel admin: Vite + React (TypeScript; web; mengelola seluruh siklus aplikasi — lihat §Panel Admin)
- Kontrak API: Protobuf (proto3) dengan buf
- Auth & push: Google Auth (Google Sign-In) untuk login/register + FCM (Firebase Cloud Messaging) untuk notifikasi; kredensial Firebase dan Google Auth di folder /home/demo/Downloads/takmungkin (komputer pemilik produk)
- Infra: Docker Compose; CI perusahaan atau GitHub Actions
- Desain: ikon Material Symbols (https://m3.material.io/styles/icons); font Manrope (file di /home/demo/Downloads/imo_fonts, bundel ke res/font)

OUTPUT YANG DIMINTA
- Monorepo: proto/, backend/, android/, admin/ (Vite + React), docs/
- README berisi cara menjalankan backend dan membangun aplikasi
- docs/ARCHITECTURE.md, docs/ASSUMPTIONS.md, dan docs/adr/ (Architecture Decision Records)
- Migrasi (golang-migrate) dan seed (10 user, 50 post contoh yang orisinal)
- Tes unit dan integrasi untuk alur utama

CARA KERJA
Baca seluruh spesifikasi. Buat daftar tugas per fase dan tampilkan. Setelah tiap fase, laporkan: yang selesai, tes yang lulus, dan yang belum.
```

## Fitur dan alur pengguna

Produk harus punya tujuh modul inti ala 9GAG: feed, vote, komentar, unggah, section, profil, dan moderasi.

| Modul | Perilaku wajib | Aturan penting |
| --- | --- | --- |
| Feed | Tab Hot, Fresh, Trending; scroll tanpa batas; kartu post berisi media, judul, jumlah vote dan komentar | Paginasi cursor; video diputar otomatis tanpa suara |
| Vote | Upvote atau downvote per post, bisa dibatalkan | Satu vote per user per post; tampilan langsung berubah, lalu disinkronkan |
| Komentar | Balasan bersarang, vote komentar, urut Top atau Terbaru | Pemilik bisa menghapus komentarnya sendiri |
| Unggah | Gambar, GIF, video pendek; judul wajib; pilih section; tag; penanda NSFW | Batas ukuran dan durasi diatur di konfigurasi |
| Section | Kategori konten, mis. Funny, Gaming, Anime, Indonesia; bisa diikuti | Satu post masuk satu section |
| Profil | Avatar, post saya, upvote saya, komentar saya, pengikut | Profil publik bisa dilihat tanpa login |
| Pencarian | Cari judul, tag, dan user | Pakai indeks teks PostgreSQL pada fase awal |
| Notifikasi | Balasan komentar, vote pada post saya, pengikut baru; dikirim lewat FCM | Bisa dimatikan per jenis |
| Moderasi | Lapor post atau komentar, antrean moderator, sembunyikan atau hapus, blokir sementara | Setiap tindakan moderator tercatat |
| Akun | Daftar, masuk (email/password atau Google Auth), reset password | Pengunjung tanpa akun boleh membaca; vote, komentar, dan unggah butuh login; kredensial Firebase/Google Auth di /home/demo/Downloads/takmungkin (komputer pemilik produk) |

## Fitur unik (pembeda dari 9GAG)

Fitur di bawah ini adalah alasan orang mencoba aplikasi ini padahal 9GAG sudah ada. Semuanya wajib dibangun (lihat ATURAN KERAS 8). Fitur inti 9GAG tetap harus sama; fitur unik ini yang membedakan.

| Fitur unik | Perilaku wajib | Aturan penting |
| --- | --- | --- |
| Meme Duel | Dua post diadu satu lawan satu; pengguna vote salah satu; ada bracket dan juara mingguan | Satu vote per user per duel; duel berjalan 24 jam; pemenang masuk tab Duel dan dapat lencana |
| Generator meme | Bikin meme di dalam aplikasi: pilih template, isi teks atas/bawah, langsung unggah | Template resmi dikurasi; template kiriman pengguna butuh persetujuan moderator |
| Kanal daerah | Section khusus humor daerah: Sunda, Jawa, Betawi, dan lainnya; konten boleh berbahasa daerah | Nama dan deskripsi section berbahasa Indonesia; moderator per kanal |
| Tantangan harian | Tiap hari satu tema (mis. "macet"); pengguna unggah post bertema itu; post terbaik naik ke atas | Satu tema per hari, diganti otomatis tengah malam; entri memakai post biasa dengan tag tantangan |
| Sawer kreator (koin) | Pengguna menyawer kreator memakai koin; koin dibeli lewat top-up; ada leaderboard kreator mingguan | Koin dicatat di ledger (saldo, top-up, kirim, terima); top-up lewat Google Play Billing; harga paket koin diatur admin; tanpa top-up, tombol sawer nonaktif dan tidak boleh pura-pura berfungsi |
| Reaksi Indonesia | Selain upvote: Ngakak, Relate banget, Anjay, Sedih, Marah; tampil persentase "87% orang relate" | Satu reaksi per user per post; reaksi bisa diganti; reaksi tidak memengaruhi hot_score |
| Balas komentar bergambar | Komentar bisa dibalas dengan gambar atau GIF | Balasan gambar memakai alur unggah yang sama (batas ukuran, verifikasi tipe, moderasi) |
| Share ke Status WA | Dari post, satu ketuk untuk membagikan ke Status WhatsApp | Pakai Android share intent + tautan kembali ke aplikasi; media diunduh dulu dari CDN |

### Alur Meme Duel

1. Sistem memilih dua post terbit dari section yang sama (atau tema yang sama) dan membuka duel 24 jam.
2. Pengguna membuka tab Duel, melihat dua post berdampingan, dan vote salah satu.
3. Skor duel tampil langsung (persentase, bukan jumlah mentah).
4. Setelah 24 jam, duel ditutup; post pemenang mendapat lencana "Juara Duel" dan masuk daftar juara mingguan.
5. Juara mingguan dihitung dari jumlah kemenangan duel dalam satu minggu.

### Alur generator meme

1. Pengguna membuka "Buat Meme", memilih template dari daftar (template resmi + template yang disetujui).
2. Ia mengisi teks atas dan bawah, mengatur ukuran, lalu melihat pratinjau.
3. Ia menekan "Unggah sebagai post" — gambar hasil render dibuat di perangkat, lalu mengikuti alur unggah biasa (CreateUpload → CompletePost).

### Alur tantangan harian

1. Worker River (job berkala tengah malam) membuat tantangan baru dengan tema hari itu dan menutup tantangan kemarin.
2. Pengguna mengunggah post dan menandai "Ikut tantangan hari ini" (atau memakai tag tantangan).
3. Entri tantangan tampil di halaman tantangan, diurutkan pakai hot_score yang sama.
4. Saat tantangan ditutup, tiga entri teratas mendapat lencana dan tampil di arsip tantangan.

### Alur sawer kreator (koin)

1. Pengguna membeli koin dulu: Profil → "Beli Koin" → memilih paket koin (mis. 100 koin, 500 koin) → bayar lewat Google Play Billing → saldo koin bertambah; pembelian tercatat di ledger.
2. Pengguna menekan "Sawer" di post → memilih jumlah koin → konfirmasi.
3. Server memindahkan koin dari saldo pengirim ke saldo kreator dalam satu transaksi (saldo tidak boleh minus; bila saldo kurang, tampil "Koin tidak cukup" + tombol "Beli Koin").
4. Kreator menerima notifikasi "Kamu menerima saweran X koin"; riwayat sawer tampil di Profil kreator.
5. Leaderboard mingguan dihitung dari total koin yang diterima dalam satu minggu.
6. Penarikan koin oleh kreator (pencairan ke uang) TIDAK termasuk fase awal — diputuskan pemilik produk bila sudah ada (lihat docs/ASSUMPTIONS.md).

### Alur baca dan vote

1. Pengunjung membuka beranda dan feed Hot dimuat.
2. Saat scroll mendekati dasar, halaman berikutnya dimuat otomatis.
3. Pengunjung menekan panah vote. Tanpa login, muncul modal masuk; dengan login, vote tersimpan dan skor berubah.

### Alur masuk dengan Google

1. Pengguna menekan "Masuk dengan Google" di layar masuk.
2. Aplikasi membuka Google Sign-In; pengguna memilih akun Google-nya.
3. Aplikasi mengirim ID token Google ke AuthService.LoginWithGoogle.
4. Server memverifikasi ID token ke Google; bila email belum terdaftar, dibuat akun baru (username dibuat dari nama/email, bisa diganti setelahnya); bila sudah terdaftar, pengguna langsung masuk.
5. Server menerbitkan JWT seperti login biasa; seluruh alur setelahnya sama.
6. Kredensial Firebase/Google Auth dibaca server dari folder /home/demo/Downloads/takmungkin di komputer pemilik produk — bukan dari repo.

### Alur unggah

1. Pengguna yang sudah masuk menekan Unggah.
2. Ia memilih file, mengisi judul, memilih section, dan menambah tag.
3. Aplikasi memanggil RPC CreateUpload; server memeriksa tipe dan ukuran yang dinyatakan, lalu mengembalikan URL unggah bertanda tangan.
4. Aplikasi mengunggah file langsung ke penyimpanan media, lalu memanggil CompletePost.
5. Worker River memverifikasi isi file, membuat thumbnail, dan mengonversi video; post berstatus "memproses".
6. Setelah selesai, status menjadi "terbit" dan post masuk feed Fresh.

### Alur moderasi

1. Pengguna melaporkan post atau komentar dengan satu alasan dari daftar.
2. Laporan masuk antrean moderator.
3. Moderator memilih: abaikan, sembunyikan, atau hapus dan blokir penulis.
4. Penulis menerima notifikasi bila kontennya diturunkan.

## Arsitektur sistem

&#91;embedded content: arsitektur sistem · 2 klien, 1 backend, 2 penyimpanan, 1 worker\]

Aplikasi memanggil backend lewat gRPC, kecuali unggah file yang langsung ke penyimpanan memakai URL bertanda tangan. Worker River berjalan di aplikasi Go yang sama, mengambil tugas dari PostgreSQL, membuat thumbnail, dan memperbarui status post.

### Keputusan desain

- Kontrak API ditulis sekali dalam file .proto; stub Go dan Kotlin dihasilkan otomatis dari file yang sama, dan buf breaking di CI menjaga kompatibilitas.
- Transportasi memakai gRPC di atas HTTP/2 dan TLS; JWT dikirim lewat metadata dan diperiksa oleh interceptor di server.
- Unggah media memakai URL bertanda tangan ke penyimpanan, bukan streaming lewat gRPC.
- Tugas latar memakai River di atas PostgreSQL, jadi Redis tidak dibutuhkan pada tahap awal.
- Feed memakai paginasi cursor, dan feed Hot di-cache di memori (in-memory cache dengan TTL).
- Vote ditulis dalam satu transaksi database bersama penghitung di tabel posts.

### Peringkat Hot

Rumus awal; agent boleh menyetel pangkat dan konstanta setelah ada data nyata.

```latex
\text{hot\_score} = \frac{\text{upvotes} - \text{downvotes}}{(\text{umur\_jam} + 2)^{1.5}}
```

Trending memakai jendela waktu pendek: post yang vote-nya naik paling cepat dalam beberapa jam terakhir.

## Cara kerja sistem (end-to-end)

Bagian ini menjelaskan cara kerja aplikasi seperti 9GAG bekerja — dari post diunggah sampai muncul di feed, naik peringkat, dan dimoderasi.

### Diagram sistem

```text
Android (Kotlin) ──gRPC/TLS──> Go (gRPC) ──GORM──> PostgreSQL
      │                              │  River worker
      │                              ├─ thumbnail & transcode
      │                              ├─ hitung hot_score & trending
      │                              ├─ cron tantangan harian
      │                              └─ kirim notifikasi (FCM)
      └─ unggah langsung ──URL bertanda tangan──> Cloudflare R2 (media)
Panel admin (Vite + React) ──> antrean laporan ──> PostgreSQL
CDN ──> menyajikan media ke aplikasi
```

### Siklus hidup post (cara kerja feed ala 9GAG)

1. Post terbit masuk **Fresh** — urut kronologis terbaru dulu. Fresh adalah pintu masuk semua konten baru.
2. Setiap vote (upvote/downvote) memperbarui `hot_score` post itu (rumus di atas) dalam transaksi yang sama dengan pencatatan vote.
3. **Hot** menampilkan post terbit yang diurutkan `hot_score` tertinggi — post bagus naik, post lama turun sendiri karena umur di penyebut rumus.
4. **Trending** menampilkan post dengan kenaikan vote tercepat dalam jendela beberapa jam terakhir (bukan skor total) — post yang sedang viral, bukan yang sudah lama populer.
5. Post yang disembunyikan/dihapus moderator hilang dari semua feed. Post NSFW hanya tampil bila pengguna mengaktifkan "tampilkan NSFW".
6. Feed Hot di-cache di memori (in-memory cache dengan TTL) dan dihitung ulang berkala; Fresh dan Trending selalu baca dari database dengan paginasi cursor.

### Siklus hidup vote dan komentar

1. Vote disimpan optimistis di aplikasi (tampilan berubah langsung), lalu dikirim ke server; bila server menolak (mis. belum login), tampilan dikembalikan.
2. Satu vote per user per post dijamin unik di database; vote baru mengganti vote lama dalam satu transaksi beserta penghitungnya.
3. Komentar bersarang disimpan dengan `parent_id`; daftar komentar dimuat per post, urut Top (skor vote komentar) atau Terbaru.
4. Komentar yang dilaporkan dan disembunyikan moderator hilang dari daftar, tetapi balasan di bawahnya tetap terjaga strukturnya.

### Siklus hidup media

1. Aplikasi meminta URL unggah bertanda tangan (CreateUpload) — server hanya memeriksa tipe dan ukuran yang dinyatakan.
2. File diunggah langsung ke penyimpanan media (bukan lewat backend).
3. Worker River memverifikasi isi file (tipe dari isi, bukan dari nama), membuat thumbnail, mengonversi video ke format standar, lalu mengubah status post menjadi "terbit".
4. Media disajikan lewat CDN; aplikasi memakai Coil (gambar) dan ExoPlayer (video, putar otomatis tanpa suara).
5. Media disimpan di Cloudflare R2 (S3-compatible, egress gratis). URL unggah adalah presigned URL berumur pendek. Kredensial R2 dibaca dari folder /home/demo/Downloads/cloudflare di komputer pemilik produk; jangan hardcode, jangan commit ke repo.

### Siklus hidup moderasi

1. Laporan masuk antrean dengan status "menunggu".
2. Moderator di panel admin web melihat antrean, membuka konten, dan memutuskan: abaikan, sembunyikan, atau hapus + blokir penulis.
3. Setiap keputusan tercatat di `mod_actions` dan penulis menerima notifikasi bila kontennya diturunkan.
4. Konten yang disembunyikan/dihapus tidak muncul di feed, pencarian, maupun profil publik.

### Siklus hidup notifikasi

1. Kejadian (balasan komentar, vote pada post saya, pengikut baru, konten diturunkan) ditulis ke tabel `notifications`.
2. Worker mengirim push (FCM) hanya untuk jenis notifikasi yang diaktifkan pengguna. Kredensial FCM (service account) dibaca dari folder /home/demo/Downloads/takmungkin di komputer pemilik produk; jangan hardcode, jangan commit ke repo.
3. Notifikasi bisa dimatikan per jenis di pengaturan.

### Skala dan ketahanan

- Backend tanpa state: bisa dijalankan lebih dari satu node di belakang load balancer; state bersama hanya di PostgreSQL dan penyimpanan media.
- Feed Hot di-cache; bila cache kosong, feed tetap bisa dihitung dari database (lebih lambat, tetap benar).
- Rate limit lewat interceptor gRPC melindungi vote, komentar, dan unggah dari spam.
- Semua tugas berat (thumbnail, transcode, hitung ulang skor, kirim push) berjalan di worker River, bukan di request gRPC.

## Panel Admin (web — Vite + React)

Panel admin berbasis web (Vite + React, TypeScript) dan mengelola **SELURUH siklus aplikasi**. Bukan LiveView. Admin login memakai akun admin (email/password + JWT; akun admin dibuat lewat seed/CLI, bukan dari aplikasi).

### Komponen UI — Astryx (mengikat)

Frontend panel admin WAJIB memakai komponen **Astryx** (https://astryx.atmeta.com/components) — katalog komponen dengan contoh siap pakai untuk setiap varian, state, dan pola.

- **Pasang sesuai panduan resmi Astryx** (ikuti instalasi/setup dari situs Astryx; jangan tebak nama paket atau cara pasang).
- **Pakai komponen Astryx yang tersedia**: App Shell + Side Nav (layout + sidebar), Table + Table Selection Toolbar (daftar pengguna/konten/moderasi/koin), Dialog + Drawer (detail & konfirmasi), Field + Text Input + Selector + Date Input + Switch + Checkbox (form), Button + Button Group + Dropdown Menu (aksi), Badge + Toast + Skeleton + Spinner + Progress Bar (status & feedback), Tabs/Tab List + Pagination + Breadcrumbs (navigasi), Card + Empty State (kartu & layar kosong).
- **Dilarang membuat komponen custom bila Astryx sudah menyediakannya.** Komponen custom hanya bila Astryx tidak punya komponen yang setara — dan itu harus dicatat di docs/ASSUMPTIONS.md beserta alasannya.
- Tampilan mengikuti gaya Astryx; warna/identitas MemeKita (ungu #6C4CF1 + teal) dipakai sebagai aksen di atas komponen Astryx, bukan dengan mengganti komponennya.

### Modul admin

| Modul | Yang bisa dilakukan admin |
| --- | --- |
| Dashboard | Statistik: user baru, post per hari, vote, duel, tantangan, koin terjual, sawer; grafik; konten yang butuh perhatian |
| Pengguna | Cari user; lihat profil + riwayat; ubah role (user/moderator/admin); blokir/buka blokir sementara atau permanen; reset password; hapus akun |
| Konten (Post) | Cari post; lihat media + skor; sembunyikan/hapus; cabut dari Hot; tandai NSFW |
| Komentar | Cari komentar; hapus; sembunyikan |
| Moderasi | Antrean laporan (filter per status/alasan); putuskan Abaikan/Sembunyikan/Hapus + Blokir + catatan; log mod_actions |
| Section & Kanal | Buat/edit/hapus section; atur kanal daerah; tetapkan moderator per kanal |
| Duel | Lihat duel berjalan/berakhir; buat duel manual; batalkan duel; lihat juara mingguan |
| Tantangan | Buat/edit tema tantangan harian; lihat entri; tetapkan pemenang |
| Template Meme | Setujui/tolak template kiriman pengguna; tambah/hapus template resmi |
| Koin & Sawer | Atur paket koin (jumlah + harga + SKU Play); lihat ledger transaksi koin; refund top-up bila perlu; lihat leaderboard sawer |
| Notifikasi | Kirim pengumuman/push ke semua user atau per segmen |
| Pengaturan Aplikasi | Feature flag (nyalakan/matikan Duel, Tantangan, Sawer); batas ukuran unggahan; teks pengumuman; mode pemeliharaan |
| Audit | Log semua tindakan admin (siapa, apa, kapan) — tidak bisa dihapus |

### Aturan admin

- Semua tindakan admin tercatat di audit log dan tidak bisa dihapus.
- Admin tidak bisa melihat password user dan tidak bisa mengubah saldo koin tanpa transaksi ledger (penyesuaian saldo harus lewat tipe 'adjustment' dengan alasan).
- Panel admin memakai auth yang sama dengan backend (JWT; role admin/moderator diperiksa di server, bukan cuma di UI).

## Clean Architecture & Maintainability (mengikat)

Semua kode — Android, backend Go, dan admin Vite + React — WAJIB clean architecture agar mudah di-maintenance ke depan. Kode yang susah dibaca, susah diubah, atau menumpuk di satu tempat adalah bug arsitektur.

### Aturan umum (semua lapisan)

- **DILARANG stub**: tidak ada fungsi kosong, `TODO`, `notImplemented()`, mock, atau data palsu yang pura-pura berfungsi. Semua fitur diimplementasi penuh dan teruji. Bila satu bagian belum bisa selesai, laporkan sebagai BELUM SELESAI — jangan ditutup dengan stub.
- **DILARANG dead code**: tidak ada kode mati, kode yang dikomentari, atau file yang tidak dipakai.
- **Satu fitur = satu modul**: tiap fitur punya modul/paket sendiri; tidak ada god class / god module; satu file maksimal ~500 baris — bila lebih, pecah.
- **DILARANG copy-paste**: kode yang dipakai dua kali atau lebih harus jadi komponen/fungsi bersama.
- **Dependency mengarah ke dalam**: lapisan luar (UI, handler, controller) memanggil lapisan dalam (domain/use-case); lapisan dalam tidak boleh tahu soal framework, database, atau UI.
- **Konfigurasi dari environment/config**: tidak ada nilai hardcode (URL, batas, kunci, tarif) di kode.
- **Penamaan jelas dan konsisten**: nama menggambarkan apa yang dilakukan; Bahasa Inggris untuk kode, Bahasa Indonesia untuk teks UI.
- **Dokumentasi hidup**: tiap modul punya README singkat; `docs/ARCHITECTURE.md` selalu diperbarui bila struktur berubah; `.proto` selalu sinkron dengan kode.

### Android (Kotlin)

- Lapisan: UI (Fragment) → ViewModel → UseCase/Repository → DataSource (gRPC/network, Room/cache). ViewModel tidak tahu soal gRPC/Room; Fragment tidak tahu soal database.
- Repository pattern untuk semua data; satu repository per agregat (Post, User, Duel, Coin).
- Dependency injection memakai Hilt; tidak ada instansiasi manual untuk dependency yang bisa di-inject.
- Coroutines + StateFlow; error ditangani eksplisit (Result/sealed class), bukan crash atau diam.
- Komponen UI bersama di satu tempat (tidak ada gaya tombol/chip ad-hoc per layar).
- ktlint + detekt bersih di CI; unit test per ViewModel/UseCase/Repository.

### Backend (Go)

- Clean Architecture: logika bisnis tinggal di usecase (mis. `posts`, `coins`), BUKAN di handler gRPC. Lapisan: handler gRPC (tipis) → usecase → repository (interface) → data (GORM).
- Handler gRPC tipis: validasi input → panggil usecase → kembalikan response.
- Model domain hanya untuk data; tidak ada logika bisnis di model atau di query.
- **GORM**: pakai GORM (driver postgres) untuk semua akses database. Model persistence GORM terpisah dari entity domain — jangan menempel tag `gorm:` di entity domain; repository memetakan entity ↔ model. DILARANG `AutoMigrate` untuk mengubah skema (skema hanya lewat migrasi golang-migrate); GORM hanya untuk query. Aturan detail ada di "Kesalahan umum AI dalam Go" di bawah.
- Worker River tipis: ambil tugas → panggil usecase → selesai.
- Fungsi murni (pure) sebisa mungkin; efek samping (DB, HTTP, push) di batas luar.
- `gofmt` + `go vet` + golangci-lint bersih di CI; `go test -race` untuk package yang memakai goroutine; unit test per usecase/repository + contract test gRPC.

### Kesalahan umum AI dalam Go (DILARANG)

AI sering menulis Go yang "kelihatan jalan" tapi menyimpan kesalahan kecil. Semua di bawah ini DILARANG — menemukannya di review berarti tugas BELUM SELESAI:

- **Error diabaikan**: `err` WAJIB diperiksa di setiap pemanggilan yang mengembalikannya; dilarang `_ =` untuk membuang error, dilarang `panic` untuk error biasa (panic hanya untuk kesalahan pemrograman yang tidak mungkin pulih).
- **Goroutine bocor**: setiap goroutine WAJIB punya kondisi berhenti (context dibatalkan / channel ditutup); dilarang `go func()` tanpa mekanisme berhenti; dilarang `time.Sleep` untuk sinkronisasi (pakai channel, WaitGroup, atau context).
- **Data race**: dilarang menulis map/slice yang sama dari dua goroutine tanpa mutex atau channel; dilarang mengoper `sync.Mutex` by value (selalu pointer); `go test -race` di CI harus lolos.
- **Context diabaikan**: setiap operasi I/O (DB, HTTP, gRPC, R2) WAJIB memakai `ctx` dari request (bukan `context.Background()` di tengah handler); timeout/deadline diset di batas luar.
- **Resource tidak ditutup**: `rows.Close()`, `resp.Body.Close()`, file, dan koneksi WAJIB `defer Close()` segera setelah dibuka; dilarang `defer` di dalam loop (tutup di akhir iterasi).
- **N+1 query**: dilarang query di dalam loop; pakai `Preload` / `Joins` / satu query batch.
- **SQL injection**: dilarang merangkai query dari string (`fmt.Sprintf` ke dalam SQL); selalu pakai parameter binding GORM.
- **GORM**: dilarang `AutoMigrate` untuk mengubah skema di production; dilarang `Save` untuk update sebagian (Save menimpa zero value — pakai `Updates`); `ErrRecordNotFound` ditangani eksplisit (bukan dianggap error umum); semua query memakai `WithContext(ctx)`; dilarang `First` tanpa `Order` (hasil tidak deterministik); multi-tulis yang harus atomik (mis. vote + penghitung) WAJIB dalam satu transaksi.
- **Error bocor ke klien**: pesan error internal (detail DB, stack trace) DILARANG dikirim ke klien; klien hanya menerima pesan aman + kode error.
- **Nilai nol yang salah**: dilarang mengembalikan `(nil, nil)`; dilarang mengembalikan error nil bersama data nil.
- **State global**: dilarang variabel global yang bisa diubah (`var db *gorm.DB` yang dipakai langsung dari handler) — DB dan dependency lain di-inject lewat konstruktor; `main.go` hanya wiring (tidak ada logika bisnis).
- **Waktu salah**: simpan dan bandingkan waktu selalu dalam UTC; dilarang `time.Now()` tanpa `.UTC()` untuk data yang disimpan; parsing tanggal memakai layout yang benar, bukan tebak format.
- **Validasi di belakang**: input WAJIB divalidasi di handler (batas luar) sebelum masuk usecase; usecase tidak boleh percaya input mentah.
- **Magic number & `any` sembarangan**: tanpa magic number/string; dilarang `interface{}`/`any` bila tipenya bisa ditulis eksplisit.
- **Logika di `main` / fungsi raksasa**: `main.go` hanya wiring; satu fungsi = satu tugas (lihat §Logika terbaik).
- **Tes palsu**: tes yang tidak pernah gagal (tanpa assert), tes yang memanggil server asli, dan tes yang bergantung pada urutan — semuanya dilarang.

CI memeriksa sebagian otomatis (`go vet`, golangci-lint dengan `errcheck`, `gosec`, `govet`, `staticcheck`, `gocritic`, `go test -race`); sisanya diperiksa saat review kode.

### Admin (Vite + React)

- Struktur per fitur (`features/users`, `features/coins`, ...); komponen bersama di `components/`; lapisan API di `services/` (panggilan HTTP terpusat).
- Komponen UI memakai Astryx (https://astryx.atmeta.com/components) — jangan buat komponen custom bila Astryx sudah menyediakannya (lihat §Panel Admin, "Komponen UI — Astryx").
- Tidak ada logika bisnis di komponen; state diambil dari hooks/services.
- TypeScript strict; ESLint + Prettier bersih di CI.
- Satu halaman = satu tanggung jawab; tabel, form, dan dialog dipakai ulang.

### Arsitektur terbaik & modern (anti-spaghetti — mengikat)

Tujuan: orang baru yang bergabung harus bisa menambah fitur tanpa harus memahami seluruh codebase. Spaghetti code — logika bercampur, dependency melingkar, satu file melakukan segalanya — DILARANG.

**Pilihan arsitektur per lapisan:**

- **Android — MVVM + Clean Architecture, modular per fitur.** Lapisan: Presentation (Fragment + ViewModel, aliran data searah/UDF dengan StateFlow) → Domain (UseCase, model domain, interface Repository — murni Kotlin, tanpa framework Android) → Data (implementasi Repository, gRPC DataSource, Room/cache). Tiap fitur jadi modul Gradle sendiri (`:feature:feed`, `:feature:duel`, `:feature:coins`, ...) dengan dependency searah: app → feature → core (domain/data). Modul fitur tidak boleh saling mengimpor.
- **Backend Go — Clean Architecture + Ports & Adapters.** Usecase adalah port (logika bisnis); adapter (repository GORM, gRPC handler, River worker, FCM, R2, Play Billing) tinggal di batas luar dan hanya memanggil usecase. Service eksternal (FCM, R2, Play Billing) dibungkus interface + implementasi, supaya bisa diganti tanpa mengubah logika bisnis.
- **Admin — Feature-Sliced Design (React + TypeScript).** Lapisan: `app/` (routing, store) → `features/` (satu folder per fitur: ui, model, api) → `entities/` → `shared/` (komponen, lib). Feature tidak boleh mengimpor feature lain; yang bersama tinggal di `shared/`.

**Aturan anti-spaghetti (mengikat):**

- DILARANG circular dependency (A memanggil B, B memanggil A) — di semua lapisan; dijaga tooling (Android: aturan dependency modul Gradle; Go: `go list` / `go mod graph` + aturan dependency antar package internal; React: aturan import/eslint boundaries).
- DILARANG god class/module: satu kelas/modul hanya punya SATU tanggung jawab; bila berubah karena dua alasan berbeda, pecah.
- DILARANG mencampur UI + logika + data dalam satu file.
- DILARANG mengakses database/network langsung dari UI atau handler — selalu lewat Repository/Usecase.
- Interface (kontrak) antarlapisan ditulis eksplisit; lapisan dalam tidak tahu soal lapisan luar.
- Satu perubahan fitur hanya menyentuh modul fitur itu + kontraknya — bila harus mengubah lima tempat untuk satu fitur kecil, strukturnya salah, rapikan dulu.

**Cara menambah fitur baru (panduan untuk developer baru):**

1. Tulis/ubah kontrak `.proto` dulu → generate stub → `buf breaking` harus lolos.
2. Backend: tambah usecase → repository → handler gRPC tipis → worker bila perlu → test usecase.
3. Android: tambah UseCase → Repository/DataSource → ViewModel → UI → test.
4. Admin: tambah modul di `features/` → halaman + tabel/form → hubungkan ke API.
5. Perbarui `docs/ARCHITECTURE.md` bila struktur berubah; catat keputusan arsitektur penting di `docs/adr/NNNN-judul.md` (Architecture Decision Record).
6. Jalankan seluruh tes + lint sebelum tandai selesai.

### Logika terbaik & penulisan kode (mengikat)

- **Pilih algoritma dan struktur data yang tepat** untuk masalahnya: jangan O(n²) bila O(n) atau O(log n) tersedia; pakai indeks database, paginasi, dan cache — jangan memuat semua data sekaligus.
- **Satu fungsi = satu tugas**: fungsi pendek dan jelas; bila satu fungsi melakukan tiga hal, pecah jadi tiga fungsi. Dilarang fungsi raksasa dan nested `if` berlapis-lapis — pakai early return / guard clause.
- **Tangani edge case**: input kosong, null, jaringan gagal, saldo kurang, data tidak ditemukan — semuanya ditangani eksplisit; jangan biarkan crash, diam, atau menampilkan data salah.
- **Tanpa magic number/string**: semua angka dan teks yang bermakna diberi nama konstanta yang jelas.
- **Tanpa duplikasi logika (DRY)**: logika yang sama di dua tempat harus jadi satu fungsi bersama.
- **Error handling eksplisit**: setiap kegagalan menghasilkan pesan/log yang bermakna; jangan menelan error diam-diam.
- **Kode dibaca seperti cerita**: nama variabel dan fungsi menjelaskan maksudnya; dilarang singkatan misterius dan nama satu huruf (kecuali indeks loop).
- **Performa**: dilarang query N+1; dilarang pekerjaan berat di UI thread; dilarang polling boros — pakai cache, paginasi, dan worker latar.
- **Review diri sendiri**: sebelum menandai satu tugas selesai, baca ulang kode sendiri dan tanyakan: "apakah orang lain bisa paham ini dalam 1 menit?" Bila tidak, rapikan dulu.

### CI (penjaga aturan)

- CI gagal bila: lint tidak bersih, ada `TODO`/`FIXME`/`notImplemented()`/stub di kode produksi, ada file yang tidak dipakai, atau coverage turun di bawah ambang yang ditetapkan.
- Tidak ada pengecualian diam-diam: bila satu aturan harus dilanggar, catat alasannya di `docs/ASSUMPTIONS.md` dan minta persetujuan pemilik produk.

## Mesin Dev (2 GB RAM — mengikat)

Komputer yang dipakai untuk development hanya punya **2 GB RAM**. Kemarin komputer sempat **crash dan freeze** saat develop mobile, test, dan running backend bersamaan. Agent WAJIB pintar memakai memori — jangan ulangi kejadian itu.

### Aturan umum

- **Jangan jalankan semuanya bersamaan.** Satu pekerjaan berat dalam satu waktu: build Android ATAU running backend ATAU test ATAU dev server admin. Jangan pernah keempatnya sekaligus.
- **Jangan jalankan emulator.** Uji di HP fisik (armeabi-v7a) — itu juga HP uji utama.
- **Tutup aplikasi berat lain** selama build/test. Bila memakai Android Studio, aktifkan Power Save Mode; bila terasa berat, build lewat terminal/CLI (`./gradlew`) saja.
- **Siapkan swap** 2 GB di mesin dev sebagai pengaman (bila belum ada).
- Bila build tetap freeze/crash, turunkan heap lebih jauh dan laporkan — JANGAN menaikkan heap di mesin ini.
- Build berat (rilis AAB, full test suite) jalankan di GitHub Actions/CI, bukan di mesin dev.

### Android (Gradle) — `gradle.properties` (wajib dipakai persis)

```properties
org.gradle.jvmargs=-Xmx1024m -XX:MaxMetaspaceSize=256m
org.gradle.workers.max=1
org.gradle.parallel=false
org.gradle.caching=true
org.gradle.configuration-cache=true
kotlin.daemon.jvmargs=-Xmx512m
android.useAndroidX=true
```

- Build dev hanya **satu ABI (armeabi-v7a)**; build universal hanya untuk rilis.
- Hindari `clean` bila tidak perlu; andalkan incremental build + build cache.
- R8/minify hanya untuk build rilis, JANGAN untuk build dev harian.
- Lint boleh di-skip lokal saat iterasi cepat (`-x lint`), tetapi CI wajib lint penuh dan bersih.
- Test Android dijalankan per modul, bukan seluruh suite sekaligus.

### Backend (Go + PostgreSQL + MinIO)

- Jalankan lewat Docker Compose dengan **batas memori tiap service** (contoh: postgres 256m, minio 128m, backend 256m) — jangan biarkan satu service memakan semua RAM.
- Jangan jalankan backend + Android build + test dalam waktu yang sama.
- River workers dibatasi (1–2) di dev.
- Test backend dijalankan per package (`go test ./internal/posts/...`), bukan seluruh suite sekaligus; batasi paralelisme (`-p 1`).

### Admin (Vite + React)

- Dev server Vite dijalankan sendiri (tanpa backend + build Android bersamaan).
- Batasi memori Node: `NODE_OPTIONS=--max-old-space-size=512`.
- Untuk cek hasil, utamakan `vite build` sekali jalan daripada dev server menyala terus.

### Test

- Test dijalankan per modul/batch kecil; seluruh suite hanya di CI.
- Bila satu test batch membuat mesin freeze, bagi lagi menjadi batch yang lebih kecil dan laporkan.

## Sistem desain (ikon, font, warna)

Semua layar WAJIB memakai sistem desain ini. Jangan memakai ikon, font, atau warna lain tanpa perintah.

### Ikon — Material Symbols

- Sumber: https://m3.material.io/styles/icons (Material Symbols).
- Gaya: Outlined untuk keadaan biasa, Filled untuk keadaan aktif/terpilih.
- Ukuran: 24dp (ikon biasa), 20dp (di dalam chip), 32dp (ikon besar di empty state).
- Warna: mengikuti token warna di bawah; ikon nonaktif memakai warna sekunder.
- Ikon wajib per aksi:

| Aksi | Nama ikon Material Symbols |
| --- | --- |
| Beranda | home |
| Cari | search |
| Buat / Unggah | add_circle |
| Duel | swords |
| Tantangan | emoji_events |
| Notifikasi | notifications |
| Profil | person |
| Vote naik / turun | arrow_upward / arrow_downward |
| Komentar | chat_bubble |
| Bagikan | share |
| Simpan | bookmark |
| Lainnya | more_vert |
| Laporkan | flag |
| Sawer | payments |
| Pengaturan | settings |
| Kembali | arrow_back |
| Tutup | close |

### Font — Manrope

- Nama font: **Manrope**. File font ada di folder /home/demo/Downloads/imo_fonts di komputer pemilik produk — bundel file font ke `res/font` aplikasi; jangan unduh dari internet saat build.
- Bobot: Regular 400, Medium 500, SemiBold 600, Bold 700, ExtraBold 800.
- Ukuran: Judul layar 20sp Bold; Judul post 16sp SemiBold; Isi/komentar 14sp Regular; Caption/waktu 12sp Regular; Tombol 14sp SemiBold; Tab 14sp Medium.
- Fallback: bila Manrope gagal dimuat, pakai sans-serif sistem — aplikasi tetap harus jalan.
- Lisensi: Manrope berlisensi SIL Open Font License (bebas dipakai, termasuk komersial).

### Warna (identitas sendiri — bukan palet 9GAG/Lahelu)

- Tema terang: Primary #6C4CF1 (ungu), Primary gelap #5A3ED6, Aksen #00B8A9 (teal), Latar #FFFFFF, Permukaan #F6F5FA, Teks utama #1A1A1A, Teks sekunder #6B7280, Garis #E8E6F0, Sukses #22C55E, Bahaya #EF4444, Peringatan #F59E0B.
- Tema gelap: Latar #121212, Permukaan #1E1E24, Teks utama #F2F2F2, Teks sekunder #9CA3AF, Primary sama.
- Aturan: tombol utama memakai Primary; vote aktif memakai Primary; badge belum dibaca memakai Primary; aksi berbahaya (hapus/blokir) memakai Bahaya; jangan memakai warna lain untuk aksi utama.

### Jarak, radius, komponen

- Jarak memakai kelipatan 4dp (4, 8, 12, 16, 24). Padding layar 16dp.
- Radius: kartu 12dp, tombol 8dp, chip pill (radius penuh), avatar bulat.
- Bottom navigation: Beranda, Duel, Buat (+, tombol tengah menonjol), Tantangan, Profil. Cari lewat ikon di app bar; notifikasi lewat ikon lonceng di app bar.
- App bar: judul + ikon Cari + ikon lonceng (dengan badge jumlah belum dibaca).
- Kartu post: media di atas, judul di bawah media, bar aksi (vote, komentar, bagikan, simpan) di bawah judul.

## Identitas Merek & Logo (mengikat)

Logo dan identitas merek harus REAL — digambar dan dirancang sungguh-sungguh, bukan asal-asalan, bukan placeholder, bukan teks saja, bukan ikon generik.

### Logo — maskot emoji lucu

- Logo berupa **gambar maskot gaya emoji yang lucu** (wajah tertawa / ekspresi kocak), **digambar sebagai vektor asli** — bukan karakter emoji Unicode (mis. menempel "😂" sebagai teks), bukan stiker hasil tempel, bukan hasil crop dari gambar lain.
- Maskot harus orisinal (bukan meniru maskot 9GAG, Lahelu, atau merek lain) dan mudah dikenali.
- Harus terbaca di ukuran kecil (ikon 16dp) maupun besar (splash screen): bentuk sederhana, kontras jelas, tidak ada detail kecil yang hilang saat diperkecil.
- Sediakan varian: logo di atas latar terang dan di atas latar gelap (atau satu logo yang bekerja di keduanya).
- File master: SVG (vektor) di `docs/brand/` + ekspor PNG di semua ukuran yang dibutuhkan.

### Yang harus dibuat dari logo

| Pemakaian | Spesifikasi |
| --- | --- |
| Ikon aplikasi Android | Adaptive icon (foreground + background), 108x108dp; tampil benar di semua bentuk mask (bulat, squircle) |
| Splash screen | Logo maskot di tengah latar Primary #6C4CF1 + nama "MemeKita" (Manrope Bold) di bawahnya |
| App bar | Versi kecil logo (atau nama saja bila ruang sempit) |
| Play Store | Ikon 512x512 PNG + feature graphic (1024x500) memakai maskot + nama |
| Panel admin (web) | Logo kecil di sidebar + favicon |
| Empty state & onboarding | Maskot boleh dipakai ulang sebagai ilustrasi (dengan ekspresi berbeda bila perlu) |

### Aturan

- DILARANG memakai placeholder, teks saja ("MemeKita" tanpa gambar), ikon generik (mis. huruf M saja), atau gambar hasil comot dari internet.
- Logo dan splash screen HARUS ditunjukkan ke pemilik produk untuk persetujuan sebelum rilis — bila pemilik bilang jelek, buat ulang sampai disetujui.
- Nama aplikasi yang tampil: "MemeKita" (tulis persis begitu).

## Referensi UI/UX (teks — mengikat)

Tampilan harus **setara 9GAG dalam pola dan perilaku** — pengguna 9GAG langsung paham tanpa belajar: feed vertikal, kartu post, vote panah, tab Hot/Fresh/Trending. **Pembedanya**: identitas sendiri (nama MemeKita, warna ungu/teal, font Manrope, ikon Material Symbols) + fitur unik di §Fitur unik (Duel, Tantangan, reaksi Indonesia, sawer). Kualitas harus **production grade** — rapi, konsisten, dan lengkap, bukan tampilan ala tutorial. DILARANG menyalin nama, logo, palet warna, teks antarmuka, ikon, dan aset 9GAG/Lahelu.

Token warna, font, dan ikon mengikuti §Sistem desain. Yang di bawah ini adalah spesifikasi tiap layar.

### 1. Splash & Onboarding

- Splash: latar Primary #6C4CF1, logo maskot emoji lucu di tengah + nama "MemeKita" (Manrope Bold) di bawahnya (lihat §Identitas Merek & Logo), loading singkat → Beranda (tamu) atau langsung masuk bila sudah login.
- Onboarding (hanya pertama kali): 3 halaman geser — (1) "Meme terbaru setiap hari" (2) "Duel meme dan tantangan harian" (3) "Sawer kreator favoritmu" — tombol "Lewati" + "Mulai" → Beranda. Bisa dilewati.

### 2. Beranda / Feed

- App bar: logo "MemeKita" kiri; ikon Cari + lonceng (badge jumlah belum dibaca) kanan.
- Tab: Hot (aktif, garis bawah Primary), Trending, Fresh. Geser kiri/kanan antar tab.
- Kartu post: header (avatar + username + waktu relatif + chip nama section + ⋮); judul (16sp SemiBold, maks 2 baris); media (gambar/GIF/video, rasio asli, sudut 12dp, video putar otomatis tanpa suara); bar aksi (panah naik + skor + panah turun | komentar + jumlah | bagikan | simpan). Vote aktif: panah menyala Primary + haptic ringan.
- Ketuk dua kali media = upvote. Ketuk kartu = buka Detail Post.
- Scroll tanpa batas; pull-to-refresh; skeleton saat loading; empty state bila section kosong.
- Bottom navigation: Beranda, Duel, Buat (+ tengah menonjol, Primary), Tantangan, Profil.

### 3. Detail Post

- App bar: panah kembali + "Post" + ⋮ (Laporkan, Salin tautan, Sembunyikan).
- Media penuh di atas (gambar bisa di-zoom pinch; video ada putar/jeda, fullscreen, kontrol suara).
- Judul + info (username, section, waktu, jumlah vote).
- Bar aksi sama seperti kartu, ditambah tombol "Sawer" (ikon payments) bila post milik orang lain.
- Di bawah bar aksi: bar reaksi (Ngakak, Relate banget, Anjay, Sedih, Marah + jumlah); ketuk untuk memilih; reaksi terbanyak tampil sebagai "87% orang relate".
- Komentar di bawah: kolom tulis di bawah (avatar + "Tulis komentar..." + ikon gambar + kirim); urutan Top/Terbaru; komentar bersarang menjorok; vote komentar (panah kecil); tekan lama komentar sendiri → Hapus.

### 4. Tab Duel

- Judul tab "Duel". Kartu duel: dua post berdampingan, masing-masing media + judul; tombol vote di bawah tiap post; persentase skor tampil langsung (bar progres).
- Timer "Berakhir dalam 13j 24m" di atas kartu.
- Di bawah: "Juara minggu ini" — daftar 3 pemenang dengan lencana.
- Duel berakhir: kartu menampilkan pemenang dengan lencana "Juara Duel".

### 5. Tab Tantangan

- Header: tema hari ini (besar, "Tema hari ini: macet") + hitung mundur ke tengah malam + tombol "Ikut Tantangan" (Primary).
- Daftar entri: kartu post biasa dengan badge tantangan, diurutkan hot_score.
- Arsip: "Tantangan kemarin" — 3 pemenang + tautan lihat semua.

### 6. Buat / Unggah

- Layar "Buat": dua pilihan besar — "Unggah Meme" (ikon image) dan "Buat Meme" (ikon add_circle).
- Unggah: pilih file (galeri/kamera) → pratinjau → judul (wajib, counter karakter) → section (dropdown) → tag (chip, tambah) → toggle NSFW → tombol "Terbitkan" (Primary). Status "Memproses..." → sukses → buka post di Fresh.
- Validasi: judul kosong → pesan merah di bawah kolom; file terlalu besar → pesan jelas.

### 7. Buat Meme (generator)

- Grid template (gambar kecil + nama); kolom cari template.
- Editor: pratinjau besar; kolom teks atas + teks bawah; slider ukuran teks; tombol "Unggah sebagai post" (Primary).
- Tombol "Ajukan Template" untuk template kiriman pengguna (moderator menyetujui).

### 8. Cari

- Kolom cari + tombol batal; riwayat pencarian + "Tren saat ini" (tag populer).
- Hasil: tab Post | Tag | User. Kartu post ringkas; user: avatar + nama + tombol Ikuti.
- Tanpa hasil: empty state "Tidak ditemukan, coba kata lain".

### 9. Notifikasi

- App bar "Notifikasi"; tab "Semua" / "Belum dibaca".
- Baris: avatar pelaku + teks ("Budi mengomentari post kamu") + thumbnail post + waktu relatif; belum dibaca: latar permukaan + titik Primary.
- Ketuk → buka kontennya. Tekan lama → hapus / tandai dibaca.
- Empty state: "Belum ada notifikasi".

### 10. Profil

- Profil sendiri: avatar besar + tombol "Edit Profil"; username + bio; statistik (Post, Pengikut, Mengikuti); tab Post saya | Upvote saya | Komentar saya | Disimpan.
- Profil orang lain: tombol "Ikuti" (Primary) / "Mengikuti" (outline); tombol ⋮ (Blokir, Laporkan).
- Edit Profil: ganti avatar (kamera/galeri), nama tampilan, bio, tautan.

### 11. Pengaturan

- Daftar menu berkelompok: Akun (Edit Profil, Ganti Password, Hapus Akun), Preferensi (Tema, Bahasa, Hemat Data), Notifikasi, Penyimpanan (Bersihkan Cache + ukuran), Bantuan (Pusat Bantuan, Kebijakan Privasi, Syarat Layanan, Tentang + versi).
- Tiap baris: ikon + label + panah ›; toggle untuk on/off.

### 12. Masuk / Daftar

- Layar Masuk: logo + tombol "Lanjutkan dengan Google" (putih, ikon G) di atas; pemisah "atau"; kolom email + password; tombol "Masuk" (Primary); tautan "Lupa password?" + "Belum punya akun? Daftar".
- Daftar: username + email + password + konfirmasi; checkbox syarat & usia; tombol "Buat Akun".
- Error: pesan merah di bawah kolom yang salah; loading di tombol saat proses.

### 13. Section / Kanal daerah

- Daftar section: kartu (ikon + nama + jumlah post + tombol Ikuti).
- Halaman section: header (nama + deskripsi + tombol Ikuti), tab Hot/Fresh di dalam section, feed seperti Beranda.
- Kanal daerah: section Sunda/Jawa/Betawi dengan badge "Daerah".

### 14. Panel Admin (web — Vite + React)

- Layout: sidebar kiri (Dashboard, Pengguna, Konten, Moderasi, Duel, Tantangan, Template, Koin & Sawer, Notifikasi, Pengaturan, Audit) + area kerja kanan.
- Dashboard: kartu statistik + grafik + daftar hal yang butuh perhatian.
- Tabel di semua modul: cari, filter, sortir, paginasi; klik baris → halaman detail.
- Detail laporan: pratinjau konten + riwayat penulis; tombol Abaikan / Sembunyikan / Hapus + Blokir + kolom catatan.
- Koin & Sawer: tabel paket koin (jumlah, harga, SKU, aktif/nonaktif) + ledger transaksi koin + leaderboard.
- Log audit bisa dilihat dan dicari; tidak bisa dihapus.

### Aturan UX production-grade (mengikat)

- Animasi: transisi layar 200ms; vote ada animasi skala kecil + haptic ringan; skeleton loading; pull-to-refresh di semua feed.
- Gesture: geser kiri/kanan antar tab; ketuk dua kali media = upvote.
- Semua teks Bahasa Indonesia; tidak ada teks Inggris di layar pengguna.
- Dark theme penuh (token gelap di §Sistem desain); tidak ada layar yang tertinggal terang.
- Aksesibilitas: content description semua ikon, target sentuh minimal 48dp, kontras cukup.
- DILARANG: UI palsu/hardcode yang pura-pura berfungsi; data contoh yang tampil sebagai data nyata; menyalin aset 9GAG/Lahelu.

## Skenario pengguna (user scenarios)

Skenario ini adalah cara pengguna memakai aplikasi dari awal sampai akhir. Agent WAJIB memastikan setiap skenario bisa dijalankan tanpa pengguna bingung — bila satu langkah tidak jelas di layar, itu bug UX.

Persona: **Tamu** (tanpa login), **Pengguna baru**, **Pengguna aktif**, **Kreator**, **Moderator**, **Admin**.

### Skenario 1 — Tamu membuka aplikasi pertama kali

Tujuan: tamu langsung paham aplikasi tanpa harus daftar.
1. Tamu membuka aplikasi → tampil splash singkat → layar Beranda (tab Hot) langsung terlihat, tanpa dipaksa login.
2. Tamu scroll feed → post berikutnya dimuat otomatis.
3. Tamu menekan vote → muncul modal "Masuk untuk vote" dengan tombol Masuk dan Daftar.
Hasil: tamu bisa membaca semua feed tanpa akun; semua aksi (vote, komentar, unggah) meminta login dulu.

### Skenario 2 — Daftar akun baru (email)

1. Tamu menekan "Daftar" → isi username, email, password → tekan "Buat Akun".
2. Aplikasi mengirim email verifikasi; pengguna menekan tautan di email.
3. Pengguna kembali ke aplikasi → sudah masuk → tampil sapa "Selamat datang, [nama]!" dan tombol "Jelajahi Feed".
Hasil: akun aktif, pengguna masuk ke Beranda.

### Skenario 3 — Masuk dengan Google

1. Pengguna menekan "Masuk dengan Google" → memilih akun Google.
2. Bila email belum terdaftar, akun baru dibuat otomatis (username dari nama, bisa diganti di Profil).
3. Pengguna masuk ke Beranda.
Hasil: masuk dalam 2 ketukan, tanpa password.

### Skenario 4 — Vote dan batal vote

1. Pengguna menekan panah naik pada post → panah menyala, skor naik 1.
2. Pengguna menekan panah yang sama lagi → vote batal, skor kembali.
3. Pengguna menekan panah turun → vote berubah jadi downvote.
Hasil: skor selalu sesuai dengan vote terakhir; tidak bisa vote ganda.

### Skenario 5 — Komentar dan balas komentar

1. Pengguna membuka detail post → menekan ikon komentar → menulis komentar → kirim.
2. Komentar muncul di urutan Terbaru; pengguna beralih ke urutan Top.
3. Pengguna menekan "Balas" pada komentar orang → balasan tampil menjorok di bawah komentar itu.
4. Pengguna menekan ikon gambar di kolom komentar → memilih GIF → terkirim sebagai balasan bergambar.
Hasil: komentar dan balasan tampil di tingkat yang benar.

### Skenario 6 — Unggah meme pertama

1. Pengguna menekan tombol Buat (+) → memilih gambar dari galeri.
2. Mengisi judul (wajib), memilih section, menambah tag.
3. Menekan "Terbitkan" → tampil status "Memproses..." → beberapa detik kemudian post muncul di Fresh.
4. Pengguna membuka Profil → post-nya ada di "Post saya".
Hasil: post terbit dan bisa ditemukan di Fresh.

### Skenario 7 — Buat meme dari template

1. Pengguna menekan Buat → "Buat Meme" → memilih template.
2. Mengisi teks atas dan bawah → melihat pratinjau langsung.
3. Menekan "Unggah sebagai post" → mengikuti alur unggah biasa.
Hasil: meme buatan sendiri terbit sebagai post.

### Skenario 8 — Ikut Meme Duel

1. Pengguna membuka tab Duel → melihat dua post berdampingan.
2. Menekan salah satu post untuk vote → persentase skor tampil langsung.
3. Setelah duel berakhir, pemenang mendapat lencana; pengguna bisa melihat daftar juara mingguan.
Hasil: duel bisa diikuti tanpa penjelasan tambahan.

### Skenario 9 — Ikut tantangan harian

1. Pengguna membuka tab Tantangan → melihat tema hari ini ("macet").
2. Menekan "Ikut Tantangan" → mengunggah post dengan tag tantangan otomatis.
3. Entri tampil di halaman tantangan, diurutkan seperti feed.
Hasil: entri tercatat sebagai peserta tantangan hari itu.

### Skenario 10 — Beli koin dan sawer kreator

1. Pengguna membuka Profil → "Beli Koin" → memilih paket (mis. 100 koin) → bayar lewat Google Play → saldo koin bertambah, tampil "Pembelian berhasil".
2. Pengguna membuka post → menekan "Sawer" → memilih jumlah koin → konfirmasi → tampil "Terima kasih! Saweranmu terkirim." dan saldo berkurang.
3. Kreator menerima notifikasi "Kamu menerima saweran X koin".
4. Bila saldo kurang: tampil "Koin tidak cukup" + tombol "Beli Koin".
Hasil: koin berpindah dan tercatat di ledger; kreator bisa melihat riwayat di Profil.

### Skenario 11 — Bagikan ke Status WA

1. Pengguna membuka post → menekan Bagikan → memilih "Status WhatsApp".
2. WhatsApp terbuka dengan media + tautan post; pengguna memposting ke Status.
3. Orang yang mengetuk tautan di Status masuk ke aplikasi (atau toko aplikasi bila belum terpasang).
Hasil: post tersebar ke WhatsApp dalam 2 ketukan.

### Skenario 12 — Cari meme dan user

1. Pengguna menekan ikon Cari → mengetik "kucing" → hasil post dan tag tampil.
2. Beralih ke tab "User" → menemukan kreator → membuka Profil → menekan "Ikuti".
Hasil: pencarian menemukan post, tag, dan user; follow langsung berfungsi.

### Skenario 13 — Laporkan konten

1. Pengguna membuka post yang melanggar → menekan ⋮ → "Laporkan" → memilih alasan ("Spam") → konfirmasi.
2. Tampil "Terima kasih, laporanmu sudah kami terima."
3. Post itu hilang dari feed pengguna tersebut.
Hasil: laporan tercatat dan konten tidak terlihat lagi oleh pelapor.

### Skenario 14 — Moderator menangani laporan

1. Moderator membuka panel admin web → antrean laporan tampil (terbaru dulu).
2. Membuka satu laporan → melihat konten + alasan + riwayat penulis.
3. Memilih "Sembunyikan" + menulis catatan → konten hilang dari semua feed; penulis menerima notifikasi.
Hasil: keputusan tercatat di mod_actions dan bisa diaudit.

### Skenario 15 — Atur notifikasi

1. Pengguna membuka Pengaturan → Notifikasi.
2. Mematikan "Vote pada post saya", membiarkan "Balasan komentar" menyala.
3. Keluar dari pengaturan; tidak ada push vote lagi, push balasan tetap masuk.
Hasil: pengaturan per jenis langsung berlaku.

### Skenario 16 — Hapus akun

1. Pengguna membuka Pengaturan → Akun → "Hapus Akun Saya".
2. Aplikasi meminta konfirmasi + password (atau konfirmasi Google) → mengetik alasan (opsional).
3. Akun dan data pribadi dihapus; post yang sudah terbit tetap tampil sebagai "Pengguna dihapus" atau ikut terhapus sesuai kebijakan.
Hasil: akun benar-benar terhapus, sesuai UU PDP.

### Skenario 17 — Admin mengelola aplikasi

1. Admin login ke panel admin web (Vite + React) → Dashboard menampilkan statistik hari ini.
2. Admin membuka Moderasi → antrean laporan → memutuskan satu laporan ("Sembunyikan" + catatan).
3. Admin membuka Koin & Sawer → menambah paket koin baru (500 koin, harga Rp49.000) → paket tampil di aplikasi.
4. Admin membuka Pengaturan Aplikasi → mematikan sementara fitur Duel untuk perbaikan → aplikasi menampilkan Duel sebagai nonaktif.
Hasil: seluruh siklus aplikasi bisa dikelola dari panel admin; semua tindakan tercatat di audit log.

## Detail kecil yang wajib ada (checklist)

Hal-hal kecil ini yang paling sering terlewat. Semuanya wajib sebelum rilis — bila satu poin tidak ada, rilis belum boleh jalan.

### Status layar

- [ ] Empty state di semua layar daftar (feed kosong, komentar kosong, notifikasi kosong, pencarian tanpa hasil): ada ikon + teks penjelasan + tombol aksi, bukan layar kosong.
- [ ] Loading memakai skeleton (bukan spinner di tengah layar kosong) untuk feed dan detail.
- [ ] Error memakai pesan Bahasa Indonesia + tombol "Coba Lagi" (bukan kode error mentah).
- [ ] Offline: tampil banner "Tidak ada internet" + tampilkan konten cache terakhir; aksi (vote/komentar) yang gagal diberi tahu dan bisa diulang.
- [ ] Pull-to-refresh di semua feed.

### Navigasi

- [ ] Tombol kembali selalu ke layar sebelumnya; tidak ada jalan buntu.
- [ ] Deep link: tautan post, profil, duel, dan tantangan membuka layar yang benar (dari notifikasi, dari share WA, dari luar aplikasi).
- [ ] Mengetuk notifikasi push membuka kontennya, bukan cuma Beranda.
- [ ] Bottom navigation tidak membuat tumpukan layar yang membingungkan (satu ketuk = satu layar).

### Media

- [ ] Gambar bisa di-zoom (pinch) di detail post.
- [ ] Video: putar otomatis tanpa suara di feed; di detail ada putar/jeda, fullscreen, dan kontrol suara.
- [ ] GIF berputar otomatis; thumbnail buram saat loading.
- [ ] Post NSFW diburamkan sampai pengguna mengetuk "Lihat"; NSFW tidak tampil bila pengaturan "tampilkan NSFW" mati.
- [ ] Ukuran dan durasi unggahan dibatasi dan pesannya jelas bila melebihi batas.

### Teks dan angka

- [ ] Waktu relatif: "baru saja", "5 menit lalu", "2 jam lalu", "kemarin", lalu tanggal.
- [ ] Angka disingkat: 1,2 rb; 3,4 jt (bukan 1200 atau 3400000).
- [ ] Judul panjang dipotong 2 baris + "..."; teks panjang memakai "Baca selengkapnya".
- [ ] Semua teks UI Bahasa Indonesia; tidak ada teks Inggris yang tertinggal di layar pengguna (kecuali nama fitur yang memang diputuskan).

### Aksesibilitas

- [ ] Semua ikon punya content description untuk TalkBack.
- [ ] Target sentuh minimal 48dp.
- [ ] Kontras teks cukup di tema terang dan gelap.
- [ ] Ukuran font mengikuti pengaturan sistem (jangan kunci ukuran font).

### Akun dan data (UU PDP)

- [ ] Hapus akun dari dalam aplikasi (lihat Skenario 16).
- [ ] Ekspor data saya (post, komentar, vote) dari Pengaturan.
- [ ] Ganti password; keluar dari semua perangkat.
- [ ] Kebijakan Privasi dan Syarat Layanan bisa dibuka dari dalam aplikasi.

### Pengaturan

- [ ] Tema: Terang / Gelap / Ikuti sistem — berlaku di semua layar, tidak ada layar yang tertinggal terang.
- [ ] Bahasa (bila lebih dari satu bahasa didukung).
- [ ] Hemat data: jangan putar video otomatis saat memakai data seluler.
- [ ] Bersihkan cache; tampilkan ukuran cache.
- [ ] Notifikasi per jenis (lihat Skenario 15).

### Moderasi dan keamanan

- [ ] Blokir user dan bisukan user dari profil/komentar.
- [ ] Alasan laporan lengkap (spam, SARA, pornografi, kekerasan, hak cipta, lainnya).
- [ ] Batas usia minimum tampil saat daftar; konten NSFW butuh konfirmasi usia.
- [ ] Rate limit: bila pengguna terlalu cepat vote/komentar, tampil pesan "Terlalu cepat, coba lagi sebentar" (bukan error mentah).

### Performa

- [ ] Feed scroll mulus (target 60fps); gambar dimuat bertahap (thumbnail dulu, lalu penuh).
- [ ] Video tidak diunduh penuh sebelum terlihat di layar.
- [ ] Ukuran APK rilis wajar; cache dibatasi dan bisa dibersihkan.

### Rilis

- [ ] Ikon aplikasi (adaptive icon), splash screen (logo maskot + nama), dan nama aplikasi benar di semua tempat; logo sudah disetujui pemilik produk.
- [ ] Screenshot dan deskripsi Play Store disiapkan (Bahasa Indonesia).
- [ ] Force update: versi lama yang tidak didukung menampilkan layar update, bukan crash.
- [ ] Crash reporting dan analytics aktif tanpa mengirim data pribadi.
- [ ] Build rilis diuji di HP nyata (bukan cuma emulator).

## Skema database dan API

Skema memakai 21 tabel PostgreSQL yang dikelola migrasi (golang-migrate); kontrak API ditulis dalam Protobuf dan dipanggil lewat gRPC.

### Tabel

| Tabel | Kolom kunci |
| --- | --- |
| users | id, username (unik), email (unik), password\_hash (boleh kosong bila hanya Google Auth), google\_id (unik, boleh kosong), avatar\_url, role (user, moderator, admin), status, created\_at |
| sections | id, slug (unik), name, is\_nsfw |
| posts | id, user\_id, section\_id, title, media\_type (image, gif, video), media\_url, thumb\_url, width, height, status (processing, published, hidden, removed), is\_nsfw, upvotes, downvotes, comment\_count, hot\_score, created\_at |
| post\_tags | post\_id, tag |
| votes | user\_id, post\_id, value (+1 atau -1), created\_at; unik pada (user\_id, post\_id) |
| comments | id, post\_id, user\_id, parent\_id, body, upvotes, downvotes, status, created\_at |
| comment\_votes | user\_id, comment\_id, value; unik pada (user\_id, comment\_id) |
| follows | follower\_id, target\_type (user atau section), target\_id |
| notifications | id, user\_id, type, ref\_id, read\_at, created\_at |
| reports | id, reporter\_id, target\_type, target\_id, reason, status, handled\_by, created\_at |
| mod\_actions | id, moderator\_id, target\_type, target\_id, action, note, created\_at |
| meme\_templates | id, name, image\_url, width, height, text\_boxes (JSON), is\_official, status, created\_at |
| duels | id, post\_a\_id, post\_b\_id, starts\_at, ends\_at, status (open, closed), winner\_post\_id |
| duel\_votes | user\_id, duel\_id, post\_id; unik pada (user\_id, duel\_id) |
| challenges | id, title, theme, starts\_at, ends\_at, status (open, closed) |
| challenge\_entries | challenge\_id, post\_id; unik pada (challenge\_id, post\_id) |
| tips | id, sender\_id, receiver\_id, post\_id, amount\_coins, status (completed, failed), created\_at |
| user\_coins | user\_id (unik), balance\_coins, updated\_at |
| coin\_packages | id, coins, price\_idr, play\_sku (unik), is\_active, created\_at |
| coin\_transactions | id, user\_id, type (topup, tip\_sent, tip\_received, refund, adjustment), amount\_coins, balance\_after, ref\_id, created\_at |
| post\_reactions | user\_id, post\_id, reaction (ngakak, relate, anjay, sedih, marah); unik pada (user\_id, post\_id) |

Indeks wajib: posts (status, created\_at), posts (section\_id, hot\_score), votes (post\_id), comments (post\_id, parent\_id).

### RPC (gRPC)

Semua panggilan bersifat unary dan paket proto diberi versi v1.

| Layanan.RPC | Fungsi | Perlu login |
| --- | --- | --- |
| AuthService.Register, Login, RefreshToken | Daftar, masuk, perbarui token | Tidak |
| AuthService.LoginWithGoogle | Masuk/daftar lewat Google Auth (ID token Google diverifikasi server) | Tidak |
| FeedService.GetFeed | Daftar post per tab (Hot, Fresh, Trending), section, dan cursor | Tidak |
| PostService.GetPost | Detail post | Tidak |
| PostService.CreateUpload, CompletePost | Minta URL unggah, lalu terbitkan post | Ya |
| PostService.DeletePost | Hapus post | Ya (pemilik atau moderator) |
| VoteService.VotePost | Nilai 1, -1, atau 0 untuk batal | Ya |
| CommentService.ListComments | Komentar bersarang, urut Top atau Baru | Tidak |
| CommentService.CreateComment, VoteComment | Tulis dan vote komentar | Ya |
| SectionService.ListSections, FollowSection | Daftar dan ikuti section | Ikuti: ya |
| UserService.GetProfile, ListUserPosts | Profil dan post pengguna | Tidak |
| SearchService.Search | Cari post, tag, user | Tidak |
| NotificationService.ListNotifications, MarkRead | Notifikasi | Ya |
| ReportService.CreateReport | Lapor konten | Ya |
| ModerationService.ListReports, ResolveReport | Antrean dan keputusan moderator | Ya (moderator) |
| DuelService.ListDuels, VoteDuel, GetDuelResult | Daftar duel, vote duel, hasil duel | Vote: ya |
| ChallengeService.ListChallenges, GetChallenge | Tantangan harian dan entrinya | Tidak |
| TemplateService.ListTemplates | Daftar template generator meme | Tidak |
| ReactionService.ReactPost | Reaksi Indonesia pada post | Ya |
| TipService.SendTip, ListTips | Kirim sawer (koin), riwayat sawer | Ya |
| CoinService.ListPackages, ConfirmTopUp, GetBalance, ListTransactions | Paket koin, konfirmasi top-up (verifikasi Play Billing), saldo, riwayat transaksi koin | Ya (kecuali ListPackages) |

Contoh bentuk kontrak (proto/memekita/v1/vote.proto):

```proto
syntax = "proto3";
package memekita.v1;

service VoteService {
  rpc VotePost(VotePostRequest) returns (VotePostResponse);
}

message VotePostRequest {
  string post_id = 1;
  sint32 value = 2; // 1, -1, atau 0 untuk batal
}

message VotePostResponse {
  int64 upvotes = 1;
  int64 downvotes = 2;
}
```

## Fase pembangunan dan kriteria selesai

Agent mengerjakan delapan fase berurutan; satu fase baru dimulai setelah kriteria fase sebelumnya terpenuhi.

| Fase | Isi | Selesai bila |
| --- | --- | --- |
| 0. Fondasi | Monorepo (proto, backend, android), Docker Compose (PostgreSQL, MinIO), kontrak proto awal, buf di CI, CI backend dan Android | docker compose up berjalan, CI hijau, stub Go dan Kotlin dihasilkan dari proto |
| 1. Akun dan data | Skema SQL, migrasi (golang-migrate), seed, AuthService, layar daftar dan masuk | Tes Register, Login, dan Refresh lulus; aplikasi bisa masuk ke backend lokal |
| 2. Posting dan feed | CreateUpload, CompletePost, worker River, FeedService, layar feed dan detail | Gambar, GIF, dan video muncul di feed aplikasi setelah diproses |
| 3. Interaksi | Vote, komentar bersarang, section, profil: backend dan layar | Vote tidak ganda; balasan komentar tampil di tingkat yang benar |
| 4. Moderasi dan notifikasi | Lapor, antrean, panel admin web (Vite + React), notifikasi | Laporan masuk antrean dan konten bisa diturunkan dari panel |
| 5. Pencarian dan ketahanan | Pencarian, cache, rate limit lewat interceptor, tinjauan keamanan | Tes beban dan daftar periksa keamanan lulus |
| 6. Persiapan rilis | Build rilis yang ditandatangani (AAB), konfigurasi produksi, monitoring, kebijakan privasi | Build rilis lulus uji di perangkat nyata, checklist 'Detail kecil yang wajib ada' terpenuhi, dan siap diunggah ke toko aplikasi |
| 7. Fitur unik | Meme Duel, generator meme, kanal daerah, tantangan harian, sawer kreator, reaksi Indonesia, balas komentar bergambar, share ke Status WA | Duel berjalan penuh satu putaran; tantangan harian berganti tema; sawer tercatat di ledger; semua tes lulus |

Sebelum tiap fase agent menulis rencana tugas; setelah tiap fase ia menjalankan seluruh tes, memperbarui docs, dan melapor.

### Kriteria penerimaan akhir

- [ ] Pengunjung bisa membaca feed Hot, Fresh, dan Trending di aplikasi tanpa login
- [ ] Pengguna bisa vote, komentar bersarang, dan unggah gambar, GIF, serta video
- [ ] Moderator bisa menurunkan konten dari antrean laporan
- [ ] Tes Go, JUnit, dan Espresso lulus di CI
- [ ] File .proto lulus buf lint dan buf breaking
- [ ] Tidak ada nama, logo, aset, atau teks milik 9GAG atau Lahelu di repo
- [ ] README cukup bagi orang baru untuk menjalankan backend dan membangun aplikasi dalam satu sesi
- [ ] Meme Duel berjalan penuh: duel dibuka, divote, ditutup, dan pemenang mendapat lencana
- [ ] Generator meme bisa membuat post dari template di dalam aplikasi
- [ ] Tantangan harian berganti tema otomatis dan entrinya tampil di halaman tantangan
- [ ] Sawer memakai koin; top-up dan perpindahan koin tercatat di ledger; saldo tidak bisa minus
- [ ] Panel admin (Vite + React) bisa mengelola seluruh siklus aplikasi (pengguna, konten, moderasi, duel, tantangan, koin, audit)
- [ ] Panel admin memakai komponen Astryx (https://astryx.atmeta.com/components), bukan komponen custom bila Astryx sudah menyediakannya
- [ ] Tidak ada stub, TODO, atau data palsu di repo (diperiksa otomatis di CI)
- [ ] Struktur kode mengikuti §Clean Architecture & Maintainability (lint bersih di semua lapisan)
- [ ] Kode ditulis dengan logika terbaik (tidak ada kode berantakan, fungsi raksasa, atau nested dalam) — dinilai saat review
- [ ] Tidak ada spaghetti code: tidak ada circular dependency, tidak ada god class/module, batas antarlapisan dijaga tooling
- [ ] Backend Go mengikuti Clean Architecture + GORM tanpa kesalahan umum AI (lolos go vet, golangci-lint, dan go test -race)
- [ ] Logo maskot emoji lucu sudah dibuat (vektor asli, bukan placeholder) dan disetujui pemilik produk; splash screen, ikon aplikasi, dan Play Store memakai logo yang sama
- [ ] Build dan test bisa jalan di mesin dev 2 GB RAM tanpa crash/freeze (ikuti §Mesin Dev)
- [ ] Reaksi Indonesia tampil beserta persentase relate; balas komentar bisa memakai gambar/GIF
- [ ] Post bisa dibagikan ke Status WhatsApp dalam satu ketuk
- [ ] Sistem desain dipakai di semua layar (ikon Material Symbols, font Manrope)
- [ ] Semua skenario pengguna di bagian 'Skenario pengguna' bisa dijalankan dari awal sampai akhir tanpa bingung
- [ ] Semua layar mengikuti 'Referensi UI/UX' (pola setara 9GAG, identitas sendiri, production grade)
- [ ] Checklist 'Detail kecil yang wajib ada' terpenuhi semua

## Batasan: orisinalitas, hak cipta, moderasi, privasi

Fungsi umum boleh ditiru, tetapi identitas dan aset produk lain tidak; konten pengguna perlu aturan sejak hari pertama.

- **Orisinalitas.** Fitur seperti feed, vote, dan komentar bebas dipakai. Nama, logo, warna khas, teks antarmuka, ikon, dan kode milik 9GAG atau Lahelu tidak boleh disalin.
- **Hak cipta.** Syarat layanan menyatakan pengguna menjamin hak atas unggahannya. Sediakan formulir takedown dan prosedur penghapusan.
- **Moderasi.** Laporan wajib bisa ditindak sebelum rilis publik. Konten NSFW disembunyikan secara bawaan. Sediakan jalur cepat untuk konten ilegal.
- **Toko aplikasi.** Google Play punya kebijakan khusus untuk aplikasi dengan konten buatan pengguna, antara lain soal moderasi dan pelaporan. Periksa kebijakan terbaru sebelum rilis.
- **Privasi.** Simpan data seminimal mungkin, hapus data pribadi saat akun dihapus, dan pelajari kewajiban UU Pelindungan Data Pribadi (UU PDP) sebelum rilis.
- **Usia.** Tetapkan batas usia minimum di syarat layanan.

Saya bukan penasihat hukum; konsultasikan syarat layanan dan kebijakan privasi ke ahli hukum sebelum rilis publik.
