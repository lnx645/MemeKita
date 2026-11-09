package com.meme_kita.mobile.ui.generator

import android.graphics.Bitmap
import android.graphics.Canvas
import android.graphics.Color
import android.graphics.Paint
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.meme_kita.mobile.data.repository.PostRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

@HiltViewModel
class GeneratorViewModel @Inject constructor(private val repo: PostRepository) : ViewModel() {
    private val _preview = MutableStateFlow<Bitmap?>(null)
    val preview: StateFlow<Bitmap?> = _preview

    fun render(template: Bitmap, topText: String, bottomText: String, textSize: Float): Bitmap {
        val result = template.copy(Bitmap.Config.ARGB_8888, true)
        val canvas = Canvas(result)
        val paint = Paint().apply {
            color = Color.WHITE
            this.textSize = textSize
            isFakeBoldText = true
            setShadowLayer(4f, 2f, 2f, Color.BLACK)
            textAlign = Paint.Align.CENTER
        }
        if (topText.isNotBlank()) {
            canvas.drawText(topText.uppercase(), result.width / 2f, 60f, paint)
        }
        if (bottomText.isNotBlank()) {
            canvas.drawText(bottomText.uppercase(), result.width / 2f, result.height - 20f, paint)
        }
        _preview.value = result
        return result
    }

    fun uploadAsPost(bitmap: Bitmap, title: String, sectionId: String) {
        viewModelScope.launch {
            // compress bitmap to JPEG, then follow CreateUpload -> CompletePost flow
            // same as UploadViewModel
        }
    }
}
