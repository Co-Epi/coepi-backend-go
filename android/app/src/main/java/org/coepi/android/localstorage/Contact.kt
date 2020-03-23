package org.coepi.android.localstorage

import androidx.room.ColumnInfo
import androidx.room.Entity
import androidx.room.PrimaryKey

// Contact represents a BLE pairing between 2 devices
@Entity
class Contact {
    @PrimaryKey
    var uuidHash: String = ""

    @ColumnInfo(name = "timeStamp")
    var timeStamp: Long = 0L

    // NOTE: This should NEVER be transmitted to the server
    @ColumnInfo(name = "latitude")
    var latitude: Double? = null

    @ColumnInfo(name = "longitude")
    var longitude: Double? = null
}