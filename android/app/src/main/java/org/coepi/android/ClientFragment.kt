package org.coepi.android

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.databinding.DataBindingUtil
import androidx.fragment.app.Fragment
import androidx.recyclerview.widget.LinearLayoutManager
import org.coepi.android.databinding.FragmentClientBinding
import org.coepi.android.localstorage.*
import java.util.*

class ClientFragment : Fragment() {
    lateinit var model : CoEpiModel

    override fun onCreateView(
            inflater: LayoutInflater, container: ViewGroup?,
            savedInstanceState: Bundle?
    ): View? {
        model = CoEpiModel(context)

        val binding : FragmentClientBinding = DataBindingUtil.inflate(inflater, R.layout.fragment_client, container, false)

        val deviceAdapter = DeviceAdapter()
        val exposureAdapter = ExposureAdapter()

        // Specify the current activity as the lifecycle owner of the binding.
        // This is necessary so that the binding can observe LiveData updates.
        binding.lifecycleOwner = this

        val deviceManager = LinearLayoutManager(this.context)
        binding.deviceList.adapter = deviceAdapter
        binding.deviceList.layoutManager = deviceManager

        val exposureManager = LinearLayoutManager(this.context)
        binding.exposuresList.adapter = exposureAdapter
        binding.exposuresList.layoutManager = exposureManager

        binding.buttonExposureAndSymptoms.setOnClickListener { _ ->
            val uuIDs: List<Contact>? = model.listContacts(0, 99999999999)
            // get Symptoms from a Dialog box
            val symptomsString = binding.textSymptomReport.text.toString()
            val symptoms = Base64.getEncoder().encodeToString(symptomsString.toByteArray())
            val eas = ExposureAndSymptoms(symptoms, uuIDs)
            model.onExposureAndSymptoms(eas)
        }

        binding.buttonExposureCheck.setOnClickListener {
            val uuIDs: List<Contact>? = model.listContacts(0, 99999999999)
            val exposureCheck = ExposureCheck(uuIDs)
            model.onExposureCheck(exposureCheck)
        }
        return binding.root
    }

}