package com.meme_kita.mobile.util

import android.content.Context
import android.content.Intent
import androidx.core.content.FileProvider
import java.io.File

object ImageShare {
    fun shareToWhatsAppStatus(context: Context, mediaUrl: String, postUrl: String) {
        val intent = Intent(Intent.ACTION_SEND).apply {
            type = "text/plain"
            putExtra(Intent.EXTRA_TEXT, "$postUrl\nLihat meme lucu di MemeKita!")
            `package` = "com.whatsapp"
        }
        try {
            context.startActivity(intent)
        } catch (_: Exception) {
            val generic = Intent.createChooser(Intent(Intent.ACTION_SEND).apply {
                type = "text/plain"
                putExtra(Intent.EXTRA_TEXT, "$mediaUrl\n$postUrl")
            }, "Bagikan")
            context.startActivity(generic)
        }
    }
}
