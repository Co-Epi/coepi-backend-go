package org.coepi.android.localstorage.room

import androidx.room.Dao
import androidx.room.Delete
import androidx.room.Insert
import androidx.room.Query
import org.coepi.android.localstorage.Contact

@Dao
interface ContactDao {
    @get:Query("SELECT * FROM contact")
    val all: List<Contact?>?

    @Query("SELECT * FROM contact WHERE uuidHash in (:ids)")
    fun loadAllById(ids: String?): List<Contact?>?

    @Query("SELECT * FROM contact WHERE :first <= timeStamp AND timeStamp <= :last LIMIT 1")
    fun findByRange(first: Long?, last: Long?): List<Contact>?

    @Insert
    fun insert(contact: Contact?)

    @Delete
    fun delete(contact: Contact?)
}