# ADR 0001 — JWT + Argon2 untuk Auth

- Status: diterima
- Tanggal: 2026-11-09
- Konteks: butuh auth gRPC tanpa session, hash password aman.
- Keputusan: Joken (HS256) untuk JWT (access 1 jam, refresh 7 hari), Argon2 untuk password, EncryptedSharedPreferences (Keystore) di Android.
- Konsekuensi: secret dari env/file, tidak di-commit; refresh via `AuthService.RefreshToken`.
