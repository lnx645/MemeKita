package com.meme_kita.mobile.ui.challenge

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.meme_kita.mobile.domain.model.Challenge
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject
import java.time.Duration
import java.time.Instant
import kotlinx.coroutines.delay

@HiltViewModel
class ChallengeViewModel @Inject constructor() : ViewModel() {
    private val _uiState = MutableStateFlow<ChallengeUiState>(ChallengeUiState.Loading)
    val uiState: StateFlow<ChallengeUiState> = _uiState
    private val _countdown = MutableStateFlow("23j 59m")
    val countdown: StateFlow<String> = _countdown

    init { startCountdown() }

    fun loadChallenge() {
        viewModelScope.launch {
            _uiState.value = ChallengeUiState.Loading
            delay(300)
            _uiState.value = ChallengeUiState.Success(
                Challenge("1", "macet", "macet", Instant.now().toString(), Instant.now().plus(Duration.ofHours(24)).toString(), "open"),
                emptyList()
            )
        }
    }

    private fun startCountdown() {
        viewModelScope.launch {
            while (true) {
                val now = Instant.now()
                val end = now.plus(Duration.ofHours(24))
                val d = Duration.between(now, end)
                _countdown.value = "${d.toHours()}j ${d.toMinutes() % 60}m"
                delay(60_000)
            }
        }
    }
}
sealed class ChallengeUiState {
    object Loading : ChallengeUiState()
    data class Success(val challenge: Challenge, val entries: List<com.meme_kita.mobile.domain.model.Post>) : ChallengeUiState()
    data class Error(val message: String) : ChallengeUiState()
}
