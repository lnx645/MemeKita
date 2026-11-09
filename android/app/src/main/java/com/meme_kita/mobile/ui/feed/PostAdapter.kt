package com.meme_kita.mobile.ui.feed

import android.view.LayoutInflater
import android.view.ViewGroup
import androidx.recyclerview.widget.DiffUtil
import androidx.recyclerview.widget.ListAdapter
import androidx.recyclerview.widget.RecyclerView
import coil.load
import com.meme_kita.mobile.R
import com.meme_kita.mobile.databinding.ItemPostBinding
import com.meme_kita.mobile.domain.model.Post
import com.meme_kita.mobile.util.TimeAgo

class PostAdapter(
    private val onVoteUp: (String) -> Unit,
    private val onVoteDown: (String) -> Unit,
    private val onDoubleTap: (String) -> Unit,
    private val onComment: (String) -> Unit,
    private val onShare: (Post) -> Unit
) : ListAdapter<Post, PostAdapter.VH>(Diff) {
    object Diff : DiffUtil.ItemCallback<Post>() {
        override fun areItemsTheSame(a: Post, b: Post) = a.id == b.id
        override fun areContentsTheSame(a: Post, b: Post) = a == b
    }
    inner class VH(val binding: ItemPostBinding) : RecyclerView.ViewHolder(binding.root) {
        fun bind(post: Post) {
            binding.tvTitle.text = post.title
            binding.tvUsername.text = post.username
            binding.tvTime.text = TimeAgo.format(post.createdAt)
            binding.chipSection.text = post.sectionName
            binding.tvScore.text = TimeAgo.formatNumber(post.upvotes - post.downvotes)
            binding.ivMedia.load(post.thumbUrl.ifEmpty { post.mediaUrl }) {
                placeholder(R.drawable.placeholder)
                crossfade(true)
            }
            binding.btnUp.setColorFilter(if (post.myVote == 1) 0xFF6C4CF1.toInt() else 0xFF6B7280.toInt())
            binding.btnDown.setColorFilter(if (post.myVote == -1) 0xFF6C4CF1.toInt() else 0xFF6B7280.toInt())
            binding.btnUp.setOnClickListener { onVoteUp(post.id) }
            binding.btnDown.setOnClickListener { onVoteDown(post.id) }
            binding.ivMedia.setOnClickListener { onDoubleTap(post.id) }
            binding.btnComment.setOnClickListener { onComment(post.id) }
            binding.btnShare.setOnClickListener { onShare(post) }
            // zoom: coil already handles, video auto-play handled via ExoPlayer in detail
            binding.btnUp.contentDescription = "Vote naik"
            binding.btnDown.contentDescription = "Vote turun"
            binding.btnComment.contentDescription = "Komentar"
            binding.btnShare.contentDescription = "Bagikan"
        }
    }
    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): VH {
        val binding = ItemPostBinding.inflate(LayoutInflater.from(parent.context), parent, false)
        return VH(binding)
    }
    override fun onBindViewHolder(holder: VH, position: Int) = holder.bind(getItem(position))
}
