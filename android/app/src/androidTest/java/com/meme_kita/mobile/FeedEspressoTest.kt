package com.meme_kita.mobile

import androidx.test.ext.junit.runners.AndroidJUnit4
import androidx.test.espresso.Espresso.onView
import androidx.test.espresso.assertion.ViewAssertions.matches
import androidx.test.espresso.matcher.ViewMatchers.*
import org.junit.Test
import org.junit.runner.RunWith

@RunWith(AndroidJUnit4::class)
class FeedEspressoTest {
    @Test
    fun feed_displays_tabs() {
        onView(withText("Hot")).check(matches(isDisplayed()))
    }
}
