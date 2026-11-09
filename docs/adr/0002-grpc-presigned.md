# ADR 0002 — gRPC + Presigned Upload

- Status: diterima
- Keputusan: API via gRPC (HTTP2+TLS, JWT di metadata, interceptor), media via presigned R2 (bukan stream gRPC).
- Alasan: hemat memori backend (2GB), CDN langsung, verifikasi worker.
