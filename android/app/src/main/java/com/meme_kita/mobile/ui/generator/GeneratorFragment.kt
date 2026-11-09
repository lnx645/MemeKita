package com.meme_kita.mobile.ui.generator

import android.graphics.BitmapFactory
import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.fragment.app.Fragment
import androidx.fragment.app.viewModels
import com.meme_kita.mobile.databinding.FragmentGeneratorBinding
import dagger.hilt.android.AndroidEntryPoint

@AndroidEntryPoint
class GeneratorFragment : Fragment() {
    private var _binding: FragmentGeneratorBinding? = null
    private val binding get() = _binding!!
    private val viewModel: GeneratorViewModel by viewModels()

    override fun onCreateView(inflater: LayoutInflater, container: ViewGroup?, savedInstanceState: Bundle?): View {
        _binding = FragmentGeneratorBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        // Grid template: official + approved user templates, searchable
        // Editor: preview big, top/bottom text, slider size, upload button
        binding.etTop.addTextChangedListener(object : android.text.TextWatcher {
            override fun afterTextChanged(s: android.text.Editable?) { updatePreview() }
            override fun beforeTextChanged(s: CharSequence?, start: Int, count: Int, after: Int) {}
            override fun onTextChanged(s: CharSequence?, start: Int, before: Int, count: Int) {}
        })
        binding.etBottom.addTextChangedListener(object : android.text.TextWatcher {
            override fun afterTextChanged(s: android.text.Editable?) { updatePreview() }
            override fun beforeTextChanged(s: CharSequence?, start: Int, count: Int, after: Int) {}
            override fun onTextChanged(s: CharSequence?, start: Int, before: Int, count: Int) {}
        })
        binding.sliderTextSize.addOnChangeListener { _, value, _ -> updatePreview(value) }
        binding.btnUpload.setOnClickListener {
            val bmp = binding.ivPreview.tag as? android.graphics.Bitmap ?: return@setOnClickListener
            viewModel.uploadAsPost(bmp, "Meme buatan", "funny")
        }
        binding.btnSubmitTemplate.setOnClickListener {
            // Ajukan template kiriman pengguna -> TemplateService.SubmitTemplate (moderator approve)
        }
    }

    private fun updatePreview(size: Float = 48f) {
        val template = BitmapFactory.decodeResource(resources, android.R.drawable.ic_menu_gallery)
        val bmp = viewModel.render(template, binding.etTop.text.toString(), binding.etBottom.text.toString(), size)
        binding.ivPreview.setImageBitmap(bmp)
        binding.ivPreview.tag = bmp
    }

    override fun onDestroyView() { super.onDestroyView(); _binding = null }
}
