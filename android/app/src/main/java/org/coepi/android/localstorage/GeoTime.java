package org.coepi.android.localstorage;

import org.jetbrains.annotations.NotNull;

public class GeoTime {
    public String timestamp = "";
    public String geohash ="";

    @NotNull
    public String toString() { return "["+timestamp+","+geohash+"]";}
}
