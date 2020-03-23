package org.coepi.android.localstorage

import androidx.room.Entity

// ExposureCheck payload is sent by client to /exposurecheck to check for symptoms
@Entity
class ExposureCheckResponse(_exposures : List<Exposure>) {
    var exposures: List<Exposure>? = _exposures
}
