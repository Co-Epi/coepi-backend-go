package org.coepi.android;
import android.content.Context;
import android.os.StrictMode;
import android.widget.Toast;

import com.android.volley.Request;
import com.android.volley.RequestQueue;
import com.android.volley.Response;
import com.android.volley.VolleyError;
import com.android.volley.toolbox.StringRequest;
import com.android.volley.toolbox.Volley;
/*import com.squareup.okhttp.MediaType;
import com.squareup.okhttp.OkHttpClient;
import com.squareup.okhttp.Request;
import com.squareup.okhttp.RequestBody;
import com.squareup.okhttp.Response;*/

import org.coepi.android.localstorage.GeoTime;
import org.coepi.android.localstorage.room.Contact;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Vector;

public class FrontEndAPIClient {

    FrontEndAPIClient (Context ctx){
        this.ctx=ctx;
    }

    private Context ctx = null;

    private String url_api = "https://coepi.wolk.com:8081";//TEST "https://www.google.com";

    private String rs="";
    private String json ="";

    private void toastEcho( String rq, String rs ){
        String echo = "Sent:"+rq+"\n got:"+rs;
        Toast toast_echo = Toast.makeText( ctx, echo, Toast.LENGTH_LONG );
        toast_echo.show();
    }

    public void sendContactAndSymptoms(String symptom_id, List<Contact> contacts){//HashMap<String, GeoTime> uuIDs

        //allow direct call in main thread
        /*StrictMode.ThreadPolicy policy = new StrictMode.ThreadPolicy.Builder().permitAll().build();
        StrictMode.setThreadPolicy(policy);*/

        //Set the request post body
        json = "{\"Symptoms\":\""+symptom_id+"\", \n" +
                "\"Contacts\":[";
        String comma="";
        for  ( Contact c : contacts )  {
            json += comma+ "{\"UUID\":\""+c.uid+"\",Date:\"2020-03-24\"}\n";
            comma=",";
        }
        json += "]}\";";

        //Instantiate RequestQueue
        RequestQueue queue = Volley.newRequestQueue(ctx);
        StringRequest stringRequest = new StringRequest(Request.Method.POST, url_api/*+"/exposureandsymptoms"*/,
                new Response.Listener<String>() {
                    @Override
                    public void onResponse(String response) {
                        // Display the first 500 characters of the response string.
                        rs =response.substring(0,500);
                        toastEcho(json, rs);
                    }
                }, new Response.ErrorListener() {
            @Override
            public void onErrorResponse(VolleyError error) {
                rs ="That didn't work!"+error.toString();
                toastEcho(json, rs);
            }
        });
        //queue.start();
        //add queue (send)
        queue.add(stringRequest);//


        /*OKHttpClient
        OkHttpClient client = new OkHttpClient();
        MediaType mediaType = MediaType.parse("text/plain");
        RequestBody body = RequestBody.create(mediaType, json );
        Request request = new Request.Builder()
                .url(url_api)
                .method("POST", body)
                .build();
        try {
            Response response = client.newCall(request).execute();
            rs = response.body().toString();
        }catch( Exception e ){
            rs = e.getClass().getCanonicalName()+" "+e.getMessage();
        }*/


        /*With HttpClient
        HttpGet postSymptomsAndContactsRequest = new HttpPost(url_api);
        //Set the API media type in http content-type header
        postSymptomsAndContactsRequest.addHeader("content-type", "application/xml");
        StringEntity symptomAndContact = new StringEntity(body));
        postRequest.setEntity(userEntity);
        //SEND IT
        HttpResponse response = httpClient.execute(postRequest);*/

    }
}
