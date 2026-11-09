package com.meme_kita.mobile.data.remote

import com.google.firebase.messaging.FirebaseMessagingService
import com.google.firebase.messaging.RemoteMessage
import android.app.NotificationChannel
import android.app.NotificationManager
import android.os.Build
import androidx.core.app.NotificationCompat
import com.meme_kita.mobile.R

class FcmService : FirebaseMessagingService() {
    override fun onNewToken(token: String) {
        // kirim token ke backend via gRPC (di AuthRepository)
    }
    override fun onMessageReceived(message: RemoteMessage) {
        val title = message.notification?.title ?: "MemeKita"
        val body = message.notification?.body ?: message.data["message"] ?: ""
        showNotification(title, body)
    }
    private fun showNotification(title: String, body: String) {
        val manager = getSystemService(NotificationManager::class.java)
        val channelId = "memekita_default"
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            manager.createNotificationChannel(NotificationChannel(channelId, "MemeKita", NotificationManager.IMPORTANCE_DEFAULT))
        }
        val notification = NotificationCompat.Builder(this, channelId)
            .setSmallIcon(R.mipmap.ic_launcher)
            .setContentTitle(title)
            .setContentText(body)
            .setAutoCancel(true)
            .build()
        manager.notify(System.currentTimeMillis().toInt(), notification)
    }
}
