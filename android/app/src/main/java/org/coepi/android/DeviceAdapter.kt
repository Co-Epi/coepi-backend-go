package org.coepi.android

import java.util.*
import android.graphics.Color
import android.view.LayoutInflater
import android.view.ViewGroup
import androidx.recyclerview.widget.DiffUtil
import androidx.recyclerview.widget.ListAdapter
import androidx.recyclerview.widget.RecyclerView
import org.coepi.android.databinding.ViewDeviceBinding

class DeviceAdapter() : ListAdapter<UUID, DeviceAdapter.ViewHolder>(DeviceDiffCallback())  {

    class DeviceDiffCallback : DiffUtil.ItemCallback<UUID>() {
        override fun areItemsTheSame(oldItem: UUID, newItem: UUID): Boolean {
            return oldItem.equals(newItem)
        }
        override fun areContentsTheSame(oldItem: UUID, newItem: UUID): Boolean {
            return oldItem.equals(newItem)
        }
    }

    class ViewHolder private constructor(val binding : ViewDeviceBinding) : RecyclerView.ViewHolder(binding.root) {

        fun bind(item: UUID) {
            binding.uuid.text = item.toString()
            binding.executePendingBindings()
        }

        companion object {
            fun from(parent: ViewGroup): ViewHolder {
                val layoutInflater = LayoutInflater.from(parent.context)
                val binding = ViewDeviceBinding.inflate(layoutInflater, parent, false)
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

}


