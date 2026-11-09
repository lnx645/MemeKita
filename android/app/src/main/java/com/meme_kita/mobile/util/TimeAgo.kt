package com.meme_kita.mobile.util

import java.time.Duration
import java.time.Instant

object TimeAgo {
    fun format(iso: String): String {
        return try {
            val instant = Instant.parse(iso)
            val now = Instant.now()
            val d = Duration.between(instant, now)
            when {
                d.toMinutes() < 1 -> "baru saja"
                d.toMinutes() < 60 -> "${d.toMinutes()} menit lalu"
                d.toHours() < 24 -> "${d.toHours()} jam lalu"
                d.toDays() == 1L -> "kemarin"
                else -> iso.substring(0,10)
            }
        } catch (_: Exception) { iso }
    }
    fun formatNumber(n: Long): String = when {
        n >= 1_000_000 -> String.format("%.1f jt", n/1_000_000.0)
        n >= 1_000 -> String.format("%.1f rb", n/1_000.0)
        else -> n.toString()
    }
}
