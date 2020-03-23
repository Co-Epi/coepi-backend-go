package org.coepi.android.localstorage;

import android.content.Context;

import androidx.room.Room;

import org.coepi.android.localstorage.room.AppDatabase;
import org.coepi.android.localstorage.room.Contact;

import java.util.HashMap;
import java.util.List;

public class RoomLocalStorage {

    private AppDatabase db = null;

    /**
     * https://medium.com/@ajaysaini.official/building-database-with-room-persistence-library-ecf7d0b8f3e9
     * https://developer.android.com/training/data-storage/room/defining-data
     */
    public RoomLocalStorage( Context context ){
       db = Room.databaseBuilder(context,
                AppDatabase.class, "coepi-exposure-symptoms").
               // allow queries on the main thread.
               // Don't do this on a real app! See PersistenceBasicSample for an example.
               allowMainThreadQueries().build();;
    }

    public void putContactInfo( String uuid, GeoTime gt ){
        Contact c = new Contact();
        c.uid = uuid;
        c.timestamp = gt.timestamp;
        c.geohash = gt.geohash;
        db.contactDao().insertAll( c );
    }

    public List<Contact > listContacts(GeoTime a, GeoTime b ){
        List<Contact > ret = db.contactDao().findByRange("000000000000000000","999999999999999999");
        return ret;
    }


}
