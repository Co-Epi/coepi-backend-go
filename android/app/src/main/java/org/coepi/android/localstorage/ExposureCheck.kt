package org.coepi.android.localstorage

import androidx.room.Entity

// ExposureCheck payload is sent by client to /exposurecheck to check for symptoms
@Entity
class ExposureCheck(_contacts: List<Contact>?) {
    var contacts: List<Contact>? = _contacts
}


