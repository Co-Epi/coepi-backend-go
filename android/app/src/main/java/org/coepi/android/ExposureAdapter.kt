package org.coepi.android


import android.view.LayoutInflater
import android.view.ViewGroup
import androidx.recyclerview.widget.DiffUtil
import androidx.recyclerview.widget.ListAdapter
import androidx.recyclerview.widget.RecyclerView
import org.coepi.android.databinding.ViewExposureBinding
import org.coepi.android.localstorage.Exposure

class ExposureAdapter() : ListAdapter<Exposure, ExposureAdapter.ViewHolder>(ExposureDiffCallback())  {

    class ViewHolder private constructor(val binding : ViewExposureBinding) : RecyclerView.ViewHolder(binding.root) {

        fun bind(item: Exposure) {
            binding.exposure = item
            binding.executePendingBindings()
        }

        companion object {
            fun from(parent: ViewGroup): ViewHolder {
                val layoutInflater = LayoutInflater.from(parent.context)
                val binding = ViewExposureBinding.inflate(layoutInflater, parent, false)
                return ViewHolder(binding)
            }
        }
    }

    // Create new views (invoked by the layout manager)
    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): ViewHolder {
        // create a new view
        return ViewHolder.from(parent)
    }

    // Replace the contents of a view (invoked by the layout manager)
    override fun onBindViewHolder(holder: ViewHolder, position: Int) {
        // get action at this position
        val fl = getItem(position)!!
        holder.bind(fl)
        holder.binding.executePendingBindings();
    }


    class ExposureDiffCallback : DiffUtil.ItemCallback<Exposure>() {
        override fun areItemsTheSame(oldItem: Exposure, newItem: Exposure): Boolean {
            return oldItem.symptoms == newItem.symptoms
        }
        override fun areContentsTheSame(oldItem: Exposure, newItem: Exposure): Boolean {
            return oldItem.symptoms == newItem.symptoms
        }
    }


}
