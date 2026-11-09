package com.meme_kita.mobile.data.repository

import com.meme_kita.mobile.data.remote.GrpcClient
import com.meme_kita.mobile.domain.model.Post
import com.meme_kita.mobile.domain.model.Result
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class PostRepository @Inject constructor(private val grpc: GrpcClient) {
    // In real app, use generated gRPC code; here we simulate via REST fallback for demo
    // For production, generated code from proto via grpc-kotlin
    suspend fun getFeed(tab: String, cursor: String?, limit: Int, includeNsfw: Boolean): Result<List<Post>> = withContext(Dispatchers.IO) {
        try {
            // Simulate network; in real implementation call FeedServiceGrpc
            Result.Success(emptyList())
        } catch (e: Exception) {
            Result.Error(e.message ?: "gagal memuat feed", e)
        }
    }
    suspend fun votePost(postId: String, value: Int): Result<Triple<Long, Long, Int>> = withContext(Dispatchers.IO) {
        try { Result.Success(Triple(0,0,value)) } catch (e: Exception) { Result.Error("gagal vote", e) }
    }
    suspend fun createUpload(mediaType: String, fileSize: Long, contentType: String, fileName: String): Result<Pair<String,String>> = withContext(Dispatchers.IO) {
        try { Result.Success("https://example.com/upload" to "key123") } catch (e: Exception) { Result.Error("gagal minta URL", e) }
    }
    suspend fun completePost(mediaKey: String, title: String, sectionId: String, tags: List<String>, isNsfw: Boolean, mediaType: String): Result<Post> = withContext(Dispatchers.IO) {
        try { Result.Success(Post("", "", "", "", "", "", "", title, mediaType, "", "", 0,0,"published", isNsfw,0,0,0,0.0,tags,"",0,false)) } catch (e: Exception) { Result.Error("gagal terbitkan", e) }
    }
}

@Singleton
class AuthRepository @Inject constructor(private val grpc: GrpcClient, private val tokenStorage: com.meme_kita.mobile.data.local.TokenStorage) {
    suspend fun login(email: String, password: String): Result<com.meme_kita.mobile.domain.model.User> = withContext(Dispatchers.IO) {
        try { Result.Success(com.meme_kita.mobile.domain.model.User("1", "user", email, "", "", "user",0,0)) } catch (e: Exception) { Result.Error("email atau password salah", e) }
    }
    suspend fun register(username: String, email: String, password: String): Result<com.meme_kita.mobile.domain.model.User> = withContext(Dispatchers.IO) {
        try { Result.Success(com.meme_kita.mobile.domain.model.User("1", username, email, "", "", "user",0,0)) } catch (e: Exception) { Result.Error("gagal daftar", e) }
    }
    suspend fun loginWithGoogle(idToken: String): Result<com.meme_kita.mobile.domain.model.User> = withContext(Dispatchers.IO) {
        try { Result.Success(com.meme_kita.mobile.domain.model.User("1", "googleuser", "google@test.id", "", "", "user",0,0)) } catch (e: Exception) { Result.Error("gagal login Google", e) }
    }
    fun logout() { tokenStorage.clear() }
}

@Singleton
class CoinRepository @Inject constructor(private val grpc: GrpcClient) {
    suspend fun getBalance(): Long = 500
    suspend fun listPackages(): List<com.meme_kita.mobile.domain.model.CoinPackage> = listOf(
        com.meme_kita.mobile.domain.model.CoinPackage("1",100,15000,"coins_100",true),
        com.meme_kita.mobile.domain.model.CoinPackage("2",500,49000,"coins_500",true)
    )
    suspend fun sendTip(postId: String, amount: Long): Result<Long> = withContext(Dispatchers.IO) {
        if (amount > getBalance()) Result.Error("Koin tidak cukup") else Result.Success(getBalance() - amount)
    }
}
