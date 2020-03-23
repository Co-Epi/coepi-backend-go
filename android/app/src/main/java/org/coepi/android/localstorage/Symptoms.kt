package org.coepi.android.localstorage

import androidx.room.ColumnInfo
import androidx.room.Entity
import androidx.room.PrimaryKey

// Contact represents a BLE pairing between 2 devices
@Entity
class Symptoms {
    // TODO: autoincrement this, add caseID
    @PrimaryKey
    var symptomID: Long = 0L

    @ColumnInfo(name = "symptoms")
    var symptoms: String = ""

    @ColumnInfo(name = "timestamp")
    var timestamp: Long = 0L
}