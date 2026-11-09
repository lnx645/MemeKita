package com.meme_kita.mobile.ui.duel

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.fragment.app.Fragment
import androidx.fragment.app.viewModels
import androidx.lifecycle.lifecycleScope
import dagger.hilt.android.AndroidEntryPoint
import com.meme_kita.mobile.databinding.FragmentDuelBinding
import kotlinx.coroutines.launch

@AndroidEntryPoint
class DuelFragment : Fragment() {
    private var _binding: FragmentDuelBinding? = null
    private val binding get() = _binding!!
    private val viewModel: DuelViewModel by viewModels()

    override fun onCreateView(inflater: LayoutInflater, container: ViewGroup?, savedInstanceState: Bundle?): View {
        _binding = FragmentDuelBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        binding.swipeRefresh.setOnRefreshListener { viewModel.loadDuels() }
        viewModel.loadDuels()
        lifecycleScope.launch {
            viewModel.uiState.collect { state ->
                when (state) {
                    is DuelUiState.Loading -> {
                        binding.skeleton.visibility = View.VISIBLE
                        binding.emptyState.visibility = View.GONE
                        binding.recycler.visibility = View.GONE
                    }
                    is DuelUiState.Success -> {
                        binding.swipeRefresh.isRefreshing = false
                        binding.skeleton.visibility = View.GONE
                        if (state.duels.isEmpty()) {
                            binding.emptyState.visibility = View.VISIBLE
                            binding.recycler.visibility = View.GONE
                            binding.tvEmpty.text = "Belum ada duel aktif — juara minggu ini akan tampil di sini"
                        } else {
                            binding.emptyState.visibility = View.GONE
                            binding.recycler.visibility = View.VISIBLE
                            // adapter with two posts side by side, vote buttons, percentage bar, timer
                        }
                    }
                    is DuelUiState.Error -> {
                        binding.swipeRefresh.isRefreshing = false
                        binding.skeleton.visibility = View.GONE
                    }
                }
            }
        }
    }

    override fun onDestroyView() { super.onDestroyView(); _binding = null }
}
