package com.meme_kita.mobile.ui.upload

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.meme_kita.mobile.data.repository.PostRepository
import com.meme_kita.mobile.domain.model.Result
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

@HiltViewModel
class UploadViewModel @Inject constructor(private val repo: PostRepository) : ViewModel() {
    private val _uiState = MutableStateFlow<UploadUiState>(UploadUiState.Idle)
    val uiState: StateFlow<UploadUiState> = _uiState

    fun validateAndUpload(title: String, sectionId: String, tags: List<String>, isNsfw: Boolean, mediaType: String, fileSize: Long, contentType: String, fileName: String) {
        if (title.isBlank()) {
            _uiState.value = UploadUiState.Error("Judul wajib diisi")
            return
        }
        if (title.length > 200) {
            _uiState.value = UploadUiState.Error("Judul maksimal 200 karakter")
            return
        }
        if (fileSize > 50*1024*1024) {
            _uiState.value = UploadUiState.Error("File terlalu besar. Maksimal 50MB")
            return
        }
        viewModelScope.launch {
            _uiState.value = UploadUiState.Loading
            when (val upload = repo.createUpload(mediaType, fileSize, contentType, fileName)) {
                is Result.Success -> {
                    val (url, key) = upload.data
                    // Simulate upload to R2 via presigned URL (in real, use OkHttp PUT)
                    when (val post = repo.completePost(key, title, sectionId, tags, isNsfw, mediaType)) {
                        is Result.Success -> _uiState.value = UploadUiState.Success(post.data.id)
                        is Result.Error -> _uiState.value = UploadUiState.Error(post.message)
                        else -> {}
                    }
                }
                is Result.Error -> _uiState.value = UploadUiState.Error(upload.message)
                else -> {}
            }
        }
    }
}

sealed class UploadUiState {
    object Idle : UploadUiState()
    object Loading : UploadUiState()
    data class Success(val postId: String) : UploadUiState()
    data class Error(val message: String) : UploadUiState()
}
