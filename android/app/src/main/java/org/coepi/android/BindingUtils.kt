package org.coepi.android

import android.widget.TextView
import androidx.databinding.BindingAdapter
import org.coepi.android.localstorage.Exposure
import java.util.*

@BindingAdapter("uuid")
fun TextView.setUUID(a: UUID?) {
    a?.let {
        text = a.toString()
    }
}

@BindingAdapter("symptoms")
fun TextView.setSymptoms(a: Exposure?) {
    a?.let {
        text = a.symptoms
    }
}