package com.meme_kita.mobile.ui.feed

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.Toast
import androidx.fragment.app.Fragment
import androidx.fragment.app.viewModels
import androidx.lifecycle.lifecycleScope
import androidx.recyclerview.widget.LinearLayoutManager
import com.google.android.material.tabs.TabLayout
import com.meme_kita.mobile.databinding.FragmentFeedBinding
import dagger.hilt.android.AndroidEntryPoint
import kotlinx.coroutines.launch

@AndroidEntryPoint
class FeedFragment : Fragment() {
    private var _binding: FragmentFeedBinding? = null
    private val binding get() = _binding!!
    private val viewModel: FeedViewModel by viewModels()
    private lateinit var adapter: PostAdapter

    override fun onCreateView(inflater: LayoutInflater, container: ViewGroup?, savedInstanceState: Bundle?): View {
        _binding = FragmentFeedBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        adapter = PostAdapter(
            onVoteUp = { viewModel.vote(it, 1) },
            onVoteDown = { viewModel.vote(it, -1) },
            onDoubleTap = { viewModel.doubleTapUpvote(it) },
            onComment = { /* nav to detail */ },
            onShare = { post -> sharePost(post) }
        )
        binding.recycler.layoutManager = LinearLayoutManager(requireContext())
        binding.recycler.adapter = adapter
        binding.swipeRefresh.setOnRefreshListener { viewModel.loadFeed(currentTab(), true) }
        binding.tabLayout.addOnTabSelectedListener(object : TabLayout.OnTabSelectedListener {
            override fun onTabSelected(tab: TabLayout.Tab) { viewModel.loadFeed(currentTab()) }
            override fun onTabUnselected(tab: TabLayout.Tab) {}
            override fun onTabReselected(tab: TabLayout.Tab) {}
        })
        lifecycleScope.launch {
            viewModel.uiState.collect { state ->
                when (state) {
                    is FeedUiState.Loading -> { binding.skeleton.visibility = View.VISIBLE; binding.recycler.visibility = View.GONE; binding.emptyState.visibility = View.GONE }
                    is FeedUiState.Success -> {
                        binding.swipeRefresh.isRefreshing = false
                        binding.skeleton.visibility = View.GONE
                        if (state.posts.isEmpty()) { binding.emptyState.visibility = View.VISIBLE; binding.recycler.visibility = View.GONE }
                        else { binding.emptyState.visibility = View.GONE; binding.recycler.visibility = View.VISIBLE; adapter.submitList(state.posts) }
                    }
                    is FeedUiState.Error -> {
                        binding.swipeRefresh.isRefreshing = false
                        binding.skeleton.visibility = View.GONE
                        Toast.makeText(requireContext(), state.message, Toast.LENGTH_SHORT).show()
                    }
                }
            }
        }
        viewModel.loadFeed("hot")
    }

    private fun currentTab(): String = when (binding.tabLayout.selectedTabPosition) {
        1 -> "trending"
        2 -> "fresh"
        else -> "hot"
    }
    private fun sharePost(post: com.meme_kita.mobile.domain.model.Post) {
        com.meme_kita.mobile.util.ImageShare.shareToWhatsAppStatus(requireContext(), post.mediaUrl, "https://memekita.id/post/${post.id}")
    }
    override fun onDestroyView() { super.onDestroyView(); _binding = null }
}
