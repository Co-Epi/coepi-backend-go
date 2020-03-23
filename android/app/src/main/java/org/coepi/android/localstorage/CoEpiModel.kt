package org.coepi.android.localstorage

import android.content.Context
import android.location.Location
import android.os.Build
import androidx.annotation.RequiresApi
import androidx.room.Room
import org.coepi.android.localstorage.room.AppDatabase
import android.util.Log
import android.widget.Toast
import com.android.volley.Request
import com.android.volley.Response
import com.android.volley.toolbox.StringRequest
import com.android.volley.toolbox.Volley

class CoEpiModel internal constructor(ctx: Context?) {
    private var ctx: Context? = null
    private val url_api = "https://coepi.wolk.com:8080"
    private var rs = ""
    private var json = ""
    private var db: AppDatabase? = null

    /**
     * https://medium.com/@ajaysaini.official/building-database-with-room-persistence-library-ecf7d0b8f3e9
     * https://developer.android.com/training/data-storage/room/defining-data
     */
    init {
        this.ctx = ctx
        // allow queries on the main thread.
        // Don't do this on a real app! See PersistenceBasicSample for an example.
        db = Room.databaseBuilder(ctx!!,
                AppDatabase::class.java, "coepi-exposure-symptoms").allowMainThreadQueries
        ().build()
    }


    private fun toastEcho(rq: String, rs: String) {
        val echo = "Sent:$rq\n got:$rs"
        val toast_echo = Toast.makeText(ctx, echo, Toast.LENGTH_LONG)
        toast_echo.show()
    }

    // ExposureAndSymptoms payload is sent by client to /exposureandsymptoms when user reports symptoms
    fun sendExposureAndSymptoms(eas: ExposureAndSymptoms) {
        //allow direct call in main thread

        // TODO: Use Kotlin Serialization
        //Set the request post body
        json = """
            {"symptoms":"${eas.symptoms}", 
            "contacts":[
            """.trimIndent()
        var comma = ""
        eas.contacts?.let {
            for (c in it) {
                val sdf = java.text.SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss'Z'")
                val date = java.util.Date( c.timeStamp * 1000L)
                val dateStamp = sdf.format(date)
                json += """
                $comma{"uuidHash":"${c.uuidHash}","dateStamp":"${dateStamp}"}
                """.trimIndent()
                comma = ","
            }
        }
        json += "]}\";"
        Log.i("ExposureAndSymptoms", json)

        //Instantiate RequestQueue
        val queue = Volley.newRequestQueue(ctx)
        val stringRequest = StringRequest(Request.Method.POST, url_api +"/exposureandsymptoms",
                Response.Listener { response -> // Display the first 500 characters of the response string.
                    rs = response.substring(0, 500)
                    toastEcho(json, rs)
                }, Response.ErrorListener { error ->
            rs = "That didn't work!$error"
            toastEcho(json, rs)
        })
        //queue.start();
        //add queue (send)
        queue.add(stringRequest)
    }

    // ExposureCheck payload is sent by client to /exposurecheck to check for symptoms
    fun sendExposureCheck(ec: ExposureCheck) {
        // Set the request post body
        // TODO: Use Kotlin Serialization
        json = """
            {"contacts":[
            """.trimIndent()
        var comma = ""
        ec.contacts?.let {
            for (c in it) {
                json += """
                $comma{"uuidHash":"${c.uuidHash}",Date:"2020-03-24"}
                """.trimIndent()
                comma = ","
            }
        }
        json += "]}\";"

        //Instantiate RequestQueue
        val queue = Volley.newRequestQueue(ctx)
        val stringRequest = StringRequest(Request.Method.POST, url_api +"/exposurecheck",
                Response.Listener { response -> // Display the first 500 characters of the response string.
                    rs = response.substring(0, 500)
                    toastEcho(json, rs)
                }, Response.ErrorListener { error ->
            rs = "That didn't work!$error"
            toastEcho(json, rs)
        })
        queue.add(stringRequest)
    }

    @RequiresApi(Build.VERSION_CODES.O)
    fun storeContact(uuidHash: String?, location : Location) {

        val c = Contact()
        c.uuidHash = uuidHash!!
        // https://github.com/drfonfon/android-kotlin-geohash
        c.latitude = location.latitude
        c.longitude = location.longitude
        c.timeStamp =  System.currentTimeMillis() / 1000L;

        db!!.contactDao().insert(c)
    }

    fun listContacts(startTimeStamp: Long, endTimeStamp: Long): List<Contact>? {
        return db!!.contactDao().findByRange(startTimeStamp, endTimeStamp)
    }

}