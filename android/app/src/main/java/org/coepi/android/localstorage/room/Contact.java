package org.coepi.android.localstorage.room;

import androidx.annotation.NonNull;
import androidx.room.ColumnInfo;
import androidx.room.Entity;
import androidx.room.PrimaryKey;

@Entity
public class Contact {
    @PrimaryKey @NonNull
    public String uid;

    @ColumnInfo(name = "timestamp")
    public String timestamp;

    @ColumnInfo(name = "geohash")
    public String geohash;
}
