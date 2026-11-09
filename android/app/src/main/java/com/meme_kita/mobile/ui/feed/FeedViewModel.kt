package com.meme_kita.mobile.ui.feed

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import androidx.paging.*
import com.meme_kita.mobile.domain.model.Post
import com.meme_kita.mobile.domain.model.Result
import com.meme_kita.mobile.domain.usecase.GetFeedUseCase
import com.meme_kita.mobile.domain.usecase.VotePostUseCase
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.*
import kotlinx.coroutines.launch
import javax.inject.Inject

@HiltViewModel
class FeedViewModel @Inject constructor(
    private val getFeed: GetFeedUseCase,
    private val votePost: VotePostUseCase
) : ViewModel() {
    private val _uiState = MutableStateFlow<FeedUiState>(FeedUiState.Loading)
    val uiState: StateFlow<FeedUiState> = _uiState.asStateFlow()

    private var currentTab = "hot"
    private var currentCursor: String? = null

    fun loadFeed(tab: String, refresh: Boolean = false) {
        currentTab = tab
        if (refresh) currentCursor = null
        viewModelScope.launch {
            _uiState.value = FeedUiState.Loading
            when (val res = getFeed(tab, currentCursor)) {
                is Result.Success -> _uiState.value = FeedUiState.Success(res.data)
                is Result.Error -> _uiState.value = FeedUiState.Error(res.message)
                Result.Loading -> {}
            }
        }
    }

    fun vote(postId: String, value: Int) {
        viewModelScope.launch {
            // optimistic update
            val current = (_uiState.value as? FeedUiState.Success)?.posts ?: return@launch
            val updated = current.map { if (it.id == postId) it.copy(myVote = value, upvotes = it.upvotes + if (value==1) 1 else 0) else it }
            _uiState.value = FeedUiState.Success(updated)
            when (val res = votePost(postId, value)) {
                is Result.Success -> {
                    val (up, down, my) = res.data
                    val final = updated.map { if (it.id == postId) it.copy(upvotes = up, downvotes = down, myVote = my) else it }
                    _uiState.value = FeedUiState.Success(final)
                }
                is Result.Error -> _uiState.value = FeedUiState.Success(current) // revert
                else -> {}
            }
        }
    }

    fun doubleTapUpvote(postId: String) = vote(postId, 1)
}

sealed class FeedUiState {
    object Loading : FeedUiState()
    data class Success(val posts: List<Post>) : FeedUiState()
    data class Error(val message: String) : FeedUiState()
}
