package com.meme_kita.mobile.data.remote

import io.grpc.ManagedChannelBuilder
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class GrpcClient @Inject constructor() {
    private val host = "10.0.2.2"
    private val port = 50051
    val channel by lazy {
        ManagedChannelBuilder.forAddress(host, port).usePlaintext().build()
    }
}
