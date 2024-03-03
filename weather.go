package main

import (
    "fmt"
    "flag"
    "net/http"
    "os"
    "time"
    "encoding/json"
)


const HH = "15"
const HHMM = "15:04"
const TAJM = "2006-01-02T15:04"
const colorRed = "\033[0;31m"
const colorNone = "\033[0m"

func get_long_lat(location string) [2]string {
  url := "https://geocoding-api.open-meteo.com/v1/search?name=" + location +
      "&count=1&language=en&format=json"
  resp,err := http.Get(url)
    if err != nil {
      fmt.Println("Error: invalid location")
      os.Exit(1)
    }
    // Todo: extract latitude and longtidue from resp
    defer resp.Body.Close()
    var coordData map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&coordData)
    var lat = coordData["results"].([]interface{})[0].(map[string]interface{})["latitude"]
    var long = coordData["results"].([]interface{})[0].(map[string]interface{})["longitude"]
    var longLat [2]string
    longLat[0] = fmt.Sprintf("%.4f", long)
    longLat[1] = fmt.Sprintf("%.4f", lat)
    return longLat
}

func get_forecast(longitude, latitude string) interface{}{
  url := "https://api.open-meteo.com/v1/forecast?latitude=" + latitude +
    "&longitude=" + longitude + 
    "&hourly=temperature_2m,rain,windspeed_10m&forecast_days=1&timezone=Europe%2FBerlin"

  resp, err := http.Get(url)
  if err != nil {
    fmt.Println("Error: invalid location")
    os.Exit(1)
  }
  var forecast map[string]interface{}
  json.NewDecoder(resp.Body).Decode(&forecast)
  hourly := forecast["hourly"].(map[string]interface{})
  return hourly
}

func print_forecast(forecast interface{}) interface{}{
  timev := forecast.(map[string]interface{})["time"].([]interface{})
  rain := forecast.(map[string]interface{})["rain"].([]interface{})
  temp := forecast.(map[string]interface{})["temperature_2m"].([]interface{})
  wind := forecast.(map[string]interface{})["windspeed_10m"].([]interface{})



  now := time.Now().Format(HH)

  fmt.Print("Time\tTemp\tWind\tRain\n")

  for key := range timev {
    timefmt, err := time.Parse(TAJM, timev[key].(string))
    if err != nil {
      os.Exit(1)
    }
    tempfmt := fmt.Sprintf("%.1f", temp[key].(float64))
    windfmt := fmt.Sprintf("%.1f", wind[key].(float64))
    rainfmt := fmt.Sprintf("%.1f", rain[key].(float64))
    if timefmt.Format(HH) == now {
      fmt.Print(">")
    }
    fmt.Print(timefmt.Format(HHMM), "\t", tempfmt, "\t", windfmt, "\t", rainfmt, "\n")
  }
  return(forecast)
}


func main() {
  location := flag.String("location", "gothenburg", "The forecast location")
  flag.Parse()
  longLat := get_long_lat(*location)
  forecast := get_forecast(longLat[0], longLat[1])
  print_forecast(forecast)
}
