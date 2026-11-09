package com.meme_kita.mobile.ui.upload

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.Toast
import androidx.fragment.app.Fragment
import androidx.fragment.app.viewModels
import androidx.lifecycle.lifecycleScope
import com.meme_kita.mobile.databinding.FragmentUploadBinding
import dagger.hilt.android.AndroidEntryPoint
import kotlinx.coroutines.launch

@AndroidEntryPoint
class UploadFragment : Fragment() {
    private var _binding: FragmentUploadBinding? = null
    private val binding get() = _binding!!
    private val viewModel: UploadViewModel by viewModels()

    override fun onCreateView(inflater: LayoutInflater, container: ViewGroup?, savedInstanceState: Bundle?): View {
        _binding = FragmentUploadBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        binding.btnPublish.setOnClickListener {
            val title = binding.etTitle.text.toString()
            val section = binding.spinnerSection.selectedItem?.toString() ?: "funny"
            // tags from chip group
            viewModel.validateAndUpload(title, section, emptyList(), binding.switchNsfw.isChecked, "image", 2*1024*1024, "image/jpeg", "meme.jpg")
        }
        binding.btnPickFile.setOnClickListener {
            // In real, launch gallery/camera picker, show preview
            binding.ivPreview.visibility = View.VISIBLE
            binding.tvPreviewHint.text = "Pratinjau: meme.jpg (2MB)"
        }
        lifecycleScope.launch {
            viewModel.uiState.collect { state ->
                when (state) {
                    is UploadUiState.Loading -> { binding.progress.visibility = View.VISIBLE; binding.btnPublish.isEnabled = false; binding.tvStatus.text = "Memproses..." }
                    is UploadUiState.Success -> {
                        binding.progress.visibility = View.GONE
                        Toast.makeText(requireContext(), "Terbit! Post #${state.postId}", Toast.LENGTH_SHORT).show()
                        // navigate to fresh feed
                    }
                    is UploadUiState.Error -> {
                        binding.progress.visibility = View.GONE
                        binding.btnPublish.isEnabled = true
                        binding.tvError.text = state.message
                        binding.tvError.visibility = View.VISIBLE
                    }
                    else -> {}
                }
            }
        }
        binding.etTitle.addTextChangedListener(object : android.text.TextWatcher {
            override fun afterTextChanged(s: android.text.Editable?) {
                binding.tvCounter.text = "${s?.length ?: 0}/200"
                if (s.isNullOrBlank()) binding.tilTitle.error = "Judul wajib diisi" else binding.tilTitle.error = null
            }
            override fun beforeTextChanged(s: CharSequence?, start: Int, count: Int, after: Int) {}
            override fun onTextChanged(s: CharSequence?, start: Int, before: Int, count: Int) {}
        })
    }

    override fun onDestroyView() { super.onDestroyView(); _binding = null }
}
