package com.meme_kita.mobile.domain.model

data class Post(
    val id: String,
    val userId: String,
    val username: String,
    val avatarUrl: String,
    val sectionId: String,
    val sectionName: String,
    val sectionSlug: String,
    val title: String,
    val mediaType: String,
    val mediaUrl: String,
    val thumbUrl: String,
    val width: Int,
    val height: Int,
    val status: String,
    val isNsfw: Boolean,
    val upvotes: Long,
    val downvotes: Long,
    val commentCount: Long,
    val hotScore: Double,
    val tags: List<String>,
    val createdAt: String,
    val myVote: Int,
    val isBookmarked: Boolean
)

data class User(
    val id: String,
    val username: String,
    val email: String,
    val avatarUrl: String,
    val bio: String,
    val role: String,
    val followerCount: Long,
    val followingCount: Long
)

data class Comment(
    val id: String,
    val postId: String,
    val userId: String,
    val username: String,
    val avatarUrl: String,
    val parentId: String?,
    val body: String,
    val mediaUrl: String?,
    val upvotes: Long,
    val downvotes: Long,
    val myVote: Int,
    val createdAt: String,
    val replies: List<Comment> = emptyList()
)

enum class ReactionType { NGAKAK, RELATE, ANJAY, SEDIH, MARAH }

data class Duel(
    val id: String,
    val postA: Post,
    val postB: Post,
    val startsAt: String,
    val endsAt: String,
    val status: String,
    val winnerPostId: String?,
    val votesA: Long,
    val votesB: Long,
    val percentA: Double,
    val percentB: Double
)

data class Challenge(
    val id: String,
    val title: String,
    val theme: String,
    val startsAt: String,
    val endsAt: String,
    val status: String
)

data class CoinPackage(
    val id: String,
    val coins: Long,
    val priceIdr: Long,
    val playSku: String,
    val isActive: Boolean
)
