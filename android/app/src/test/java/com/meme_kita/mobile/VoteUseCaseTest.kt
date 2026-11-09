package com.meme_kita.mobile

import com.meme_kita.mobile.domain.usecase.VotePostUseCase
import com.meme_kita.mobile.data.repository.PostRepository
import com.meme_kita.mobile.data.remote.GrpcClient
import kotlinx.coroutines.test.runTest
import org.junit.Assert.*
import org.junit.Test

class VoteUseCaseTest {
    @Test
    fun vote_valid_values() = runTest {
        val repo = PostRepository(GrpcClient())
        val useCase = VotePostUseCase(repo)
        val result = useCase("post1", 1)
        assertTrue(result is com.meme_kita.mobile.domain.model.Result.Success)
    }
    @Test(expected = IllegalArgumentException::class)
    fun vote_invalid_throws() = runTest {
        val repo = PostRepository(GrpcClient())
        val useCase = VotePostUseCase(repo)
        useCase("post1", 2)
    }
}
