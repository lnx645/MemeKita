package com.meme_kita.mobile.ui.generator

import android.graphics.Bitmap
import android.graphics.Canvas
import android.graphics.Color
import android.graphics.Paint
import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.fragment.app.Fragment

class GeneratorFragment : Fragment() {
    // Real meme generation: draw text on bitmap
    fun renderMeme(template: Bitmap, topText: String, bottomText: String): Bitmap {
        val result = template.copy(Bitmap.Config.ARGB_8888, true)
        val canvas = Canvas(result)
        val paint = Paint().apply { color = Color.WHITE; textSize = 48f; isFakeBoldText = true; setShadowLayer(4f,2f,2f, Color.BLACK) }
        if (topText.isNotEmpty()) canvas.drawText(topText.uppercase(), 20f, 60f, paint)
        if (bottomText.isNotEmpty()) canvas.drawText(bottomText.uppercase(), 20f, result.height - 20f, paint)
        return result
    }
    override fun onCreateView(inflater: LayoutInflater, container: ViewGroup?, savedInstanceState: Bundle?): View {
        return inflater.inflate(android.R.layout.simple_list_item_1, container, false)
    }
}
