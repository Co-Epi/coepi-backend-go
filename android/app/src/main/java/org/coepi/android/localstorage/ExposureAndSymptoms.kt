package org.coepi.android.localstorage

import androidx.room.Entity

// ExposureAndSymptoms payload is sent by client to /exposureandsymptoms when user reports symptoms
@Entity
class ExposureAndSymptoms(_symptoms: String, _contacts: List<Contact>?) {
    var symptoms: String = _symptoms
    var contacts: List<Contact>? = _contacts
}


