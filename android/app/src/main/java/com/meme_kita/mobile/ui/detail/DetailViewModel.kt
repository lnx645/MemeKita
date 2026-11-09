package com.meme_kita.mobile.ui.detail

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.meme_kita.mobile.data.repository.CoinRepository
import com.meme_kita.mobile.data.repository.PostRepository
import com.meme_kita.mobile.domain.model.Result
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

@HiltViewModel
class DetailViewModel @Inject constructor(
    private val postRepo: PostRepository,
    private val coinRepo: CoinRepository
) : ViewModel() {
    private val _uiState = MutableStateFlow<DetailUiState>(DetailUiState.Loading)
    val uiState: StateFlow<DetailUiState> = _uiState

    fun loadPost(postId: String) {
        viewModelScope.launch {
            _uiState.value = DetailUiState.Loading
            // In real, call PostService.GetPost via gRPC
            _uiState.value = DetailUiState.Error("Post tidak ditemukan")
        }
    }

    fun react(postId: String, reaction: String) {
        viewModelScope.launch {
            // call ReactionService.ReactPost, update percentage "87% orang relate"
            // Reaksi tidak memengaruhi hot_score — hanya UI
        }
    }

    fun sawer(postId: String, amount: Long) {
        viewModelScope.launch {
            _uiState.value = DetailUiState.SawerLoading
            when (val res = coinRepo.sendTip(postId, amount)) {
                is Result.Success -> _uiState.value = DetailUiState.SawerSuccess("Terima kasih! Saweranmu terkirim.")
                is Result.Error -> {
                    if (res.message.contains("Koin tidak cukup")) {
                        _uiState.value = DetailUiState.SawerError("Koin tidak cukup")
                    } else {
                        _uiState.value = DetailUiState.SawerError(res.message)
                    }
                }
                else -> {}
            }
        }
    }

    fun shareToWa(postId: String, mediaUrl: String) {
        // Uses ImageShare util — media diunduh dulu dari CDN lalu share intent
    }
}

sealed class DetailUiState {
    object Loading : DetailUiState()
    data class Success(val post: com.meme_kita.mobile.domain.model.Post) : DetailUiState()
    data class Error(val message: String) : DetailUiState()
    object SawerLoading : DetailUiState()
    data class SawerSuccess(val message: String) : DetailUiState()
    data class SawerError(val message: String) : DetailUiState()
}
