package com.meme_kita.mobile.ui.detail

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.Toast
import androidx.fragment.app.Fragment
import androidx.fragment.app.viewModels
import androidx.lifecycle.lifecycleScope
import com.meme_kita.mobile.databinding.FragmentDetailBinding
import com.meme_kita.mobile.util.ImageShare
import dagger.hilt.android.AndroidEntryPoint
import kotlinx.coroutines.launch

@AndroidEntryPoint
class DetailFragment : Fragment() {
    private var _binding: FragmentDetailBinding? = null
    private val binding get() = _binding!!
    private val viewModel: DetailViewModel by viewModels()

    override fun onCreateView(inflater: LayoutInflater, container: ViewGroup?, savedInstanceState: Bundle?): View {
        _binding = FragmentDetailBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        val postId = arguments?.getString("postId") ?: ""
        viewModel.loadPost(postId)

        // Reaksi Indonesia: Ngakak, Relate banget, Anjay, Sedih, Marah
        binding.btnNgakak.setOnClickListener { viewModel.react(postId, "ngakak") }
        binding.btnRelate.setOnClickListener { viewModel.react(postId, "relate") }
        binding.btnAnjay.setOnClickListener { viewModel.react(postId, "anjay") }
        binding.btnSedih.setOnClickListener { viewModel.react(postId, "sedih") }
        binding.btnMarah.setOnClickListener { viewModel.react(postId, "marah") }

        // Sawer koin
        binding.btnSawer.setOnClickListener {
            // Show dialog pilih jumlah koin, konfirmasi, cek saldo
            viewModel.sawer(postId, 50)
        }

        // Share WA
        binding.btnShare.setOnClickListener {
            ImageShare.shareToWhatsAppStatus(requireContext(), "https://cdn.memekita.id/img.jpg", "https://memekita.id/post/$postId")
        }

        // Komentar bergambar: ikon gambar di kolom komentar -> pilih GIF/gambar -> kirim via CreateComment dengan media_url
        binding.btnImageComment.setOnClickListener {
            // launch picker, then viewModel.createCommentWithImage
        }

        // Media zoom: pinch via PhotoView, video controls via ExoPlayer (autoplay muted in feed, controls in detail)
        // NSFW blur: if post.isNsfw, blur until tap "Lihat"

        lifecycleScope.launch {
            viewModel.uiState.collect { state ->
                when (state) {
                    is DetailUiState.Loading -> { binding.skeleton.visibility = View.VISIBLE; binding.content.visibility = View.GONE }
                    is DetailUiState.Success -> {
                        binding.skeleton.visibility = View.GONE
                        binding.content.visibility = View.VISIBLE
                        binding.tvTitle.text = state.post.title
                        // 87% orang relate -> show top reaction percentage
                    }
                    is DetailUiState.Error -> {
                        binding.skeleton.visibility = View.GONE
                        Toast.makeText(requireContext(), state.message, Toast.LENGTH_SHORT).show()
                    }
                    is DetailUiState.SawerSuccess -> Toast.makeText(requireContext(), state.message, Toast.LENGTH_SHORT).show()
                    is DetailUiState.SawerError -> {
                        if (state.message.contains("Koin tidak cukup")) {
                            // show dialog "Koin tidak cukup" + tombol "Beli Koin"
                            Toast.makeText(requireContext(), "Koin tidak cukup", Toast.LENGTH_SHORT).show()
                        }
                    }
                    else -> {}
                }
            }
        }
    }

    override fun onDestroyView() { super.onDestroyView(); _binding = null }
}
