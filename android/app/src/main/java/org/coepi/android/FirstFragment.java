package org.coepi.android;

import android.os.Bundle;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.Toast;

import androidx.annotation.NonNull;
import androidx.fragment.app.Fragment;
import androidx.navigation.fragment.NavHostFragment;

import org.coepi.android.localstorage.GeoTime;
import org.coepi.android.localstorage.RoomLocalStorage;
import org.coepi.android.localstorage.room.Contact;

import java.util.ArrayList;
import java.util.Date;
import java.util.HashMap;
import java.util.List;
import java.util.Vector;

public class FirstFragment extends Fragment {

    @Override
    public View onCreateView(
            LayoutInflater inflater, ViewGroup container,
            Bundle savedInstanceState
    ) {
        // Inflate the layout for this fragment
        RoomLocalStorage storage = new RoomLocalStorage(getContext());
        /* sample
        GeoTime gt = new GeoTime();
        gt.timestamp ="2020032223";
        gt.geohash="01230123";
        storage.putContactInfo( "Ax", gt);
        */
        return inflater.inflate(R.layout.fragment_first, container, false);
    }

    private static HashMap<String,GeoTime> uuIDstaic = new HashMap<String,GeoTime>();

    public void appendUUID( String uuid, GeoTime geotime ){
        uuIDstaic.put(uuid,geotime);
    }

    public void onViewCreated(@NonNull View view, Bundle savedInstanceState) {
        super.onViewCreated(view, savedInstanceState);

        view.findViewById(R.id.button_new_contact).setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View view) {
                String rnd = String.valueOf(Math.random());

                GeoTime gt = new GeoTime();
                gt.timestamp = new Date().toString();
                gt.geohash="";

                //uuIDs.put(date, rnd);
                RoomLocalStorage storage = new RoomLocalStorage(getContext());
                storage.putContactInfo(rnd, gt);

                String echo = "added:"+rnd+","+gt;
                Toast toast_echo = Toast.makeText( view.getContext(), echo, Toast.LENGTH_SHORT);
                toast_echo.show();
            }
        });

        view.findViewById(R.id.button_new_symptom).setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View view) {
                FrontEndAPIClient client = new FrontEndAPIClient(view.getContext());

                RoomLocalStorage storage = new RoomLocalStorage(getContext());
                List<Contact> uuIDs = storage.listContacts(null, null);
                client.sendContactAndSymptoms("SlNPTkJMT0I6c2V2ZXJlIGZldmVyLGNvdWdoaW5n", uuIDs );
            }
        });

        view.findViewById(R.id.button_first).setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View view) {
                NavHostFragment.findNavController(FirstFragment.this)
                        .navigate(R.id.action_FirstFragment_to_SecondFragment);
            }
        });
    }
}
