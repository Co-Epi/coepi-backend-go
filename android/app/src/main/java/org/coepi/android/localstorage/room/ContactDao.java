package org.coepi.android.localstorage.room;

import androidx.room.Dao;
import androidx.room.Delete;
import androidx.room.Insert;
import androidx.room.Query;

import org.coepi.android.localstorage.GeoTime;
import org.coepi.android.localstorage.room.Contact;

import java.util.List;

@Dao
public interface ContactDao {
    @Query("SELECT * FROM contact")
    List<Contact> getAll();

    @Query("SELECT * FROM contact WHERE uid IN (:uuIds)")
    List<Contact> loadAllByIds(String uuIds);

    @Query("SELECT * FROM contact WHERE :first <= timestamp AND " +
            "timestamp <= :last LIMIT 1")
    List<Contact> findByRange(String first, String last);

    @Insert
    void insertAll(Contact contact );

    @Delete
    void delete( Contact contact );
}