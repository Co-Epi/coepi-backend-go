package org.coepi.android.localstorage

import com.android.volley.Response
import com.android.volley.toolbox.JsonObjectRequest
import org.json.JSONObject
import java.nio.ByteBuffer
import java.security.MessageDigest
import java.util.*
import kotlin.coroutines.Continuation
import kotlin.coroutines.resume
import kotlin.coroutines.resumeWithException

@ExperimentalUnsignedTypes
fun compareByteArray(a : UByteArray, b : UByteArray) : Int {
    for (i in 0..a.size) {
        if (a[i] != b[i]) {
            return a[i].toInt() - b[i].toInt()
        }
    }
    return 0
}

@ExperimentalUnsignedTypes
fun computeUUIDHash(selfUUID: UUID, otherUUID: UUID) : ByteArray {
    val md = MessageDigest.getInstance("SHA-256")
    val selfBytes = getBytesFromUUID(selfUUID).toUByteArray()
    val otherBytes = getBytesFromUUID(otherUUID).toUByteArray()
    var bytes: UByteArray
    if ( compareByteArray(selfBytes, otherBytes) < 0 ) {
        bytes = selfBytes + otherBytes
    } else {
        bytes = otherBytes + selfBytes
    }
    return md.digest(bytes.toByteArray()).toUByteArray().asByteArray()
}


fun getBytesFromUUID(uuid: UUID): ByteArray {
    val bb = ByteBuffer.wrap(ByteArray(16))
    bb.putLong(uuid.mostSignificantBits)
    bb.putLong(uuid.leastSignificantBits)
    return bb.array()
}

fun getUUIDFromBytes(bytes: ByteArray): UUID {
    val byteBuffer = ByteBuffer.wrap(bytes)
    val high = byteBuffer.long
    val low = byteBuffer.long
    return UUID(high, low)
}

fun getDateStamp(timeStamp : Long) : String {
    val sdf = java.text.SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss'Z'")
    val date = Date( timeStamp * 1000L)
    return  sdf.format(date)
}

class RequestWithHeaders(method: Int, url: String, private val content: String, continuation: Continuation<JSONObject>, errorListener: Response.ErrorListener) : JsonObjectRequest(method, url, null, Response.Listener<JSONObject> { response ->
    println(response)
    continuation.resume(response)
}, Response.ErrorListener {
    continuation.resumeWithException(Exception(it.cause))
}) {
    override fun getBody(): ByteArray {
        println(content)
        return content.toByteArray()
    }
}
