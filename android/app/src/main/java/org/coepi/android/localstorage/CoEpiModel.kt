package org.coepi.android.localstorage

import android.content.Context
import android.location.Location
import android.util.Log
import android.widget.Toast
import androidx.lifecycle.ViewModel
import androidx.room.Room
import com.android.volley.Request
import com.android.volley.RequestQueue
import com.android.volley.Response
import com.android.volley.toolbox.JsonObjectRequest
import com.android.volley.toolbox.StringRequest
import com.android.volley.toolbox.Volley
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.Job
import kotlinx.coroutines.launch
import org.coepi.android.localstorage.room.AppDatabase
import org.json.JSONObject
import java.util.*
import kotlin.coroutines.Continuation
import kotlin.coroutines.resume
import kotlin.coroutines.resumeWithException
import kotlin.coroutines.suspendCoroutine

class CoEpiModel internal constructor(ctx: Context?) : ViewModel() {
    private val endpointURL = "https://coepi.wolk.com:8080"
    private var ctx: Context? = null
    private var rs = ""
    private var json = ""
    private var db: AppDatabase? = null
    private var queue: RequestQueue
    private var viewModelJob = Job()
    private val uiScope = CoroutineScope(Dispatchers.Main + viewModelJob)

    init {
        this.ctx = ctx
        queue = Volley.newRequestQueue(ctx)
        // https://medium.com/@ajaysaini.official/building-database-with-room-persistence-library-ecf7d0b8f3e9
        // https://developer.android.com/training/data-storage/room/defining-data
        db = Room.databaseBuilder(ctx!!, AppDatabase::class.java, "coepi-exposure-symptoms").build()
    }

    // ExposureAndSymptoms payload is sent by client to /exposureandsymptoms when user reports symptoms
    private suspend fun sendExposureAndSymptoms(eas: ExposureAndSymptoms) = suspendCoroutine<JSONObject> {
        // TODO: Use Kotlin Serialization instead
        json = """
            {"symptoms":"${eas.symptoms}", 
            "contacts":[
            """.trimIndent()
        var comma = ""
        eas.contacts?.let {
            for (c in it) {
                val dateStamp = getDateStamp(c.timeStamp)
                json += """
                $comma{"uuidHash":"${c.uuidHash}","dateStamp":"$dateStamp"}
                """.trimIndent()
                comma = ","
            }
        }
        json += "]}\";"
        Log.i("ExposureAndSymptoms", json)
        val stringRequest = StringRequest(Request.Method.POST, "$endpointURL/exposureandsymptoms",
                Response.Listener { response -> // Display the first 500 characters of the response string.
                    rs = response.substring(0, 500)
                    toastEcho(json, rs)
                }, Response.ErrorListener { error ->
            rs = "That didn't work!$error"
            toastEcho(json, rs)
        })
        queue.add(stringRequest)
    }

    fun onExposureAndSymptoms(eas: ExposureAndSymptoms) {
        // TODO: don't use uiScope, do in background
        uiScope.launch {
            sendExposureAndSymptoms(eas)
        }
    }

    // ExposureCheck payload is sent by client to /exposurecheck to check for symptoms
    // TODO: Set this up as a Background Service https://developer.android.com/training/run-background-service/create-service
    private suspend fun sendExposureCheck(ec: ExposureCheck) = suspendCoroutine<JSONObject> { cont ->
        // TODO: Use Kotlin Serialization instead
        json = """
            {"contacts":[
            """.trimIndent()
        var comma = ""
        ec.contacts?.let {
            for (c in it) {
                val dateStamp = getDateStamp(c.timeStamp)
                json += """
                $comma{"uuidHash":"${c.uuidHash}","date":"$dateStamp"}
                """.trimIndent()
                comma = ","
            }
        }
        json += "]}\";"

        //Instantiate RequestQueue
        val queue = Volley.newRequestQueue(ctx)
        val jsonRequest = RequestWithHeaders(Request.Method.POST, "$endpointURL/exposurecheck",
                json,
                cont, Response.ErrorListener { error ->
            rs = "That didn't work!$error"
            toastEcho(json, rs)
        })
        queue.add(jsonRequest)
    }

    fun onExposureCheck(ec: ExposureCheck) {
        // TODO: don't use uiScope, do in background
        uiScope.launch {
            sendExposureCheck(ec)
        }
    }

    fun listContacts(startTimeStamp: Long, endTimeStamp: Long): List<Contact>? {
        return db!!.contactDao().findByRange(startTimeStamp, endTimeStamp)
    }

    // TODO: Have new BLE UUIDs result in onNewContact
    @ExperimentalUnsignedTypes
    fun onNewContact(selfUUID: UUID, otherUUID: UUID, location : Location) {
        val c = Contact()
        c.uuidHash = computeUUIDHash(selfUUID, otherUUID).toString()

        // Note that while we store the Lat-Long of self and another device, we would never send this to a server, only a GeoHash, with user preferences
        c.latitude = location.latitude
        c.longitude = location.longitude
        c.timeStamp =  System.currentTimeMillis() / 1000L
        db!!.contactDao().insert(c)
    }

    // TODO: move this out, replace with observer pattern of Android Architecture Components architecture https://developer.android.com/topic/libraries/architecture
    private fun toastEcho(rq: String, rs: String) {
        val echo = "Sent:$rq\n got:$rs"
        val toast_echo = Toast.makeText(ctx, echo, Toast.LENGTH_LONG)
        toast_echo.show()
    }
}
