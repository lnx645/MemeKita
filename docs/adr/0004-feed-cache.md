# ADR 0004 — Feed Cursor + Cache

- Status: diterima
- Keputusan: paginasi cursor (`inserted_at`), Hot cache Cachex (`:feed_cache`), bust tiap vote.
- Alasan: scroll tanpa batas mulus, tanpa N+1, tanpa load semua.
