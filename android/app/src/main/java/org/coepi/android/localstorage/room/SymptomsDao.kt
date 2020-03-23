package org.coepi.android.localstorage.room

import androidx.room.Dao
import androidx.room.Delete
import androidx.room.Insert
import androidx.room.Query
import org.coepi.android.localstorage.Contact
import org.coepi.android.localstorage.Symptoms

@Dao
interface SymptomsDao {
    @get:Query("SELECT * FROM symptoms")
    val all: List<Symptoms?>?

    @Query("SELECT * FROM symptoms WHERE symptomID IN (:id)")
    fun loadAllById(id: String?): List<Symptoms?>?

    @Query("SELECT * FROM symptoms WHERE :first <= timeStamp AND timeStamp <= :last LIMIT 1")
    fun findByRange(first: String?, last: String?): List<Symptoms?>?

    @Insert
    fun insertAll(contact: Symptoms?)

    @Delete
    fun delete(contact: Symptoms?)
}