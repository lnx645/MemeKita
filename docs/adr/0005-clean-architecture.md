# ADR 0005 — Clean Architecture

- Status: diterima
- Android: `ui → ViewModel → UseCase → Repository → DataSource` (Hilt, StateFlow, Result)
- Backend: Contexts sebagai port, adapter di batas (Repo, gRPC, worker, FCM, R2)
- Admin: Feature-Sliced, komponen Astryx, service terpusat
- Aturan: no circular dep (cek `mix xref`, Gradle boundaries, eslint boundaries), no god class, satu file ≤500 baris.
