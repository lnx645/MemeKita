package com.meme_kita.mobile.ui.challenge

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.fragment.app.Fragment
import androidx.fragment.app.viewModels
import androidx.lifecycle.lifecycleScope
import dagger.hilt.android.AndroidEntryPoint
import com.meme_kita.mobile.databinding.FragmentChallengeBinding
import kotlinx.coroutines.launch

@AndroidEntryPoint
class ChallengeFragment : Fragment() {
    private var _binding: FragmentChallengeBinding? = null
    private val binding get() = _binding!!
    private val viewModel: ChallengeViewModel by viewModels()

    override fun onCreateView(inflater: LayoutInflater, container: ViewGroup?, savedInstanceState: Bundle?): View {
        _binding = FragmentChallengeBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        viewModel.loadChallenge()
        lifecycleScope.launch {
            viewModel.uiState.collect { state ->
                when (state) {
                    is ChallengeUiState.Loading -> { binding.skeleton.visibility = View.VISIBLE; binding.content.visibility = View.GONE }
                    is ChallengeUiState.Success -> {
                        binding.skeleton.visibility = View.GONE
                        binding.content.visibility = View.VISIBLE
                        binding.tvTheme.text = "Tema hari ini: ${state.challenge.theme}"
                        binding.tvEntries.text = "${state.entries.size} entri"
                    }
                    is ChallengeUiState.Error -> { binding.skeleton.visibility = View.GONE }
                }
            }
        }
        lifecycleScope.launch {
            viewModel.countdown.collect { binding.tvCountdown.text = "Berakhir dalam $it" }
        }
        binding.btnJoin.setOnClickListener { /* navigate to upload with challenge tag */ }
    }

    override fun onDestroyView() { super.onDestroyView(); _binding = null }
}
