package com.meme_kita.mobile

import com.meme_kita.mobile.util.TimeAgo
import org.junit.Assert.*
import org.junit.Test

class TimeAgoTest {
    @Test fun formatNumber_ribu() { assertEquals("1.2 rb", TimeAgo.formatNumber(1200)) }
    @Test fun formatNumber_juta() { assertEquals("3.4 jt", TimeAgo.formatNumber(3400000)) }
    @Test fun formatNumber_small() { assertEquals("42", TimeAgo.formatNumber(42)) }
}
