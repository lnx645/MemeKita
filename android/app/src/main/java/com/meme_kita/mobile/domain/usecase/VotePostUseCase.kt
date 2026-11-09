package com.meme_kita.mobile.domain.usecase

import com.meme_kita.mobile.data.repository.PostRepository
import com.meme_kita.mobile.domain.model.Result
import javax.inject.Inject

class VotePostUseCase @Inject constructor(private val repo: PostRepository) {
    suspend operator fun invoke(postId: String, value: Int): Result<Triple<Long, Long, Int>> {
        require(value in listOf(-1,0,1)) { "value harus -1,0,1" }
        return repo.votePost(postId, value)
    }
}

class GetFeedUseCase @Inject constructor(private val repo: PostRepository) {
    suspend operator fun invoke(tab: String, cursor: String?, limit: Int = 20, includeNsfw: Boolean = false) =
        repo.getFeed(tab, cursor, limit, includeNsfw)
}

class SendTipUseCase @Inject constructor(private val coinRepo: com.meme_kita.mobile.data.repository.CoinRepository) {
    suspend operator fun invoke(postId: String, amount: Long): Result<Long> = coinRepo.sendTip(postId, amount)
}
