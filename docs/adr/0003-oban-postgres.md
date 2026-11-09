# ADR 0003 — Oban di PostgreSQL

- Status: diterima
- Keputusan: tugas latar (thumbnail, transcode, hot_score, tantangan harian, FCM) pakai Oban di Postgres, tanpa Redis.
- Konsekuensi: queues `media:2, notifications:3, challenges:1` (dev 1–2).
