package main

import (
    "fmt"
    "os"
    "github.com/codingsince1985/geo-golang"
    "github.com/codingsince1985/geo-golang/bing"
    "github.com/codingsince1985/geo-golang/google"
    "github.com/codingsince1985/geo-golang/here"
    "github.com/codingsince1985/geo-golang/mapquest/nominatim"
    "github.com/codingsince1985/geo-golang/mapquest/open"
    "github.com/codingsince1985/geo-golang/opencage"
    "github.com/codingsince1985/geo-golang/mapbox"
    "github.com/codingsince1985/geo-golang/openstreetmap"
    "github.com/codingsince1985/geo-golang/chained"
    "github.com/codingsince1985/geo-golang/locationiq"
)

const (
    addr     = "Melbourne VIC"
    lat, lng = -37.813611, 144.963056
    RADIUS   = 50
    ZOOM     = 18
)

func main() {
    ExampleGeocoder()
}

func ExampleGeocoder() {
    fmt.Println("Google Geocoding API")
    try(google.Geocoder(os.Getenv("GOOGLE_API_KEY")))

    fmt.Println("Mapquest Nominatim")
    try(nominatim.Geocoder(os.Getenv("MAPQUEST_NOMINATIM_KEY")))

    fmt.Println("Mapquest Open streetmaps")
    try(open.Geocoder(os.Getenv("MAPQUEST_OPEN_KEY")))

    fmt.Println("OpenCage Data")
    try(opencage.Geocoder(os.Getenv("OPENCAGE_API_KEY")))

    fmt.Println("HERE API")
    try(here.Geocoder(os.Getenv("HERE_APP_ID"), os.Getenv("HERE_APP_CODE"), RADIUS))

    fmt.Println("Bing Geocoding API")
    try(bing.Geocoder(os.Getenv("BING_API_KEY")))

    fmt.Println("Mapbox API")
    try(mapbox.Geocoder(os.Getenv("MAPBOX_API_KEY")))

    fmt.Println("OpenStreetMap")
    try(openstreetmap.Geocoder())

    fmt.Println("LocationIQ")
    try(locationiq.Geocoder(os.Getenv("LOCATIONIQ_API_KEY"), ZOOM))

    // Chained geocoder will fallback to subsequent geocoders
    fmt.Println("ChainedAPI[OpenStreetmap -> Google]")
    try(chained.Geocoder(
        openstreetmap.Geocoder(),
        google.Geocoder(os.Getenv("GOOGLE_API_KEY")),
    ))
}

func try(geocoder geo.Geocoder) {
    location, _ := geocoder.Geocode(addr)
    if location != nil {
        fmt.Printf("%s location is (%.6f, %.6f)\n", addr, location.Lat, location.Lng)
    } else {
        fmt.Println("got <nil> location")
    }
    address, _ := geocoder.ReverseGeocode(lat, lng)
    if address != nil {
        fmt.Printf("Address of (%.6f,%.6f) is %s\n", lat, lng, address.FormattedAddress)
        fmt.Printf("Detailed address: %#v\n", address)
    } else {
        fmt.Println("got <nil> address")
    }
    fmt.Println("\n")
}
