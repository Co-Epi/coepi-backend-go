package org.coepi.android

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.fragment.app.Fragment
import androidx.navigation.fragment.NavHostFragment
import org.coepi.android.localstorage.*
import java.util.*

class ClientTestFragment : Fragment() {
    lateinit var model : CoEpiModel

    override fun onCreateView(
            inflater: LayoutInflater, container: ViewGroup?,
            savedInstanceState: Bundle?
    ): View? {
        model = CoEpiModel(context)
        // Inflate the layout for this fragment
        return inflater.inflate(R.layout.fragment_clienttest, container, false)
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        view.findViewById<View>(R.id.buttonExposureAndSymptoms).setOnClickListener { view ->
            val uuIDs: List<Contact>? = model.listContacts(0, 99999999999)
            // get Symptoms from a Dialog box
            val symptomsString = "mild fever, hard to breathe"
            val symptoms = Base64.getEncoder().encodeToString(symptomsString.toByteArray())
            val eas = ExposureAndSymptoms(symptoms, uuIDs)
            model.onExposureAndSymptoms(eas)
        }

        view.findViewById<View>(R.id.buttonExposureCheck).setOnClickListener {
            val uuIDs: List<Contact>? = model.listContacts(0, 99999999999)
            val exposureCheck = ExposureCheck(uuIDs)
            model.onExposureCheck(exposureCheck)
        }
    }

}