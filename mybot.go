/*

mybot - Illustrative Slack bot in Go

Copyright (c) 2015 RapidLoop

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type Info struct {
	PlaceID string `json:"place_id"`
	Licence string `json:"licence"`
	OsmType string `json:"osm_type"`
	OsmID string `json:"osm_id"`
	Boundingbox []string `json:"boundingbox"`
	Lat string `json:"lat"`
	Lon string `json:"lon"`
	DisplayName string `json:"display_name"`
	Class string `json:"class"`
	Type string `json:"type"`
	Importance float64 `json:"importance"`
	Icon string `json:"icon"`
	Address struct {
		City string `json:"city"`
		County string `json:"county"`
		State string `json:"state"`
		Country string `json:"country"`
		CountryCode string `json:"country_code"`
	} `json:"address"`
}

type Meteo struct {
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone string `json:"timezone"`
	Offset int `json:"offset"`
	Daily struct {
		Data []struct {
			Time int `json:"time"`
			Summary string `json:"summary"`
			Icon string `json:"icon"`
			SunriseTime int `json:"sunriseTime"`
			SunsetTime int `json:"sunsetTime"`
			MoonPhase float64 `json:"moonPhase"`
			PrecipIntensity float64 `json:"precipIntensity"`
			PrecipIntensityMax float64 `json:"precipIntensityMax"`
			PrecipIntensityMaxTime int `json:"precipIntensityMaxTime"`
			PrecipProbability float64 `json:"precipProbability"`
			PrecipType string `json:"precipType"`
			TemperatureMin float64 `json:"temperatureMin"`
			TemperatureMinTime int `json:"temperatureMinTime"`
			TemperatureMax float64 `json:"temperatureMax"`
			TemperatureMaxTime int `json:"temperatureMaxTime"`
			ApparentTemperatureMin float64 `json:"apparentTemperatureMin"`
			ApparentTemperatureMinTime int `json:"apparentTemperatureMinTime"`
			ApparentTemperatureMax float64 `json:"apparentTemperatureMax"`
			ApparentTemperatureMaxTime int `json:"apparentTemperatureMaxTime"`
			DewPoint float64 `json:"dewPoint"`
			Humidity float64 `json:"humidity"`
			WindSpeed float64 `json:"windSpeed"`
			WindBearing int `json:"windBearing"`
			Visibility float64 `json:"visibility"`
			CloudCover float64 `json:"cloudCover"`
			Pressure float64 `json:"pressure"`
			Ozone float64 `json:"ozone"`
		} `json:"data"`
	} `json:"daily"`
}

func main() {
	fmt.Println(time.Now().Unix())
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: mybot slack-bot-token\n")
		os.Exit(1)
	}

	// start a websocket-based Real Time API session
	ws, id := slackConnect(os.Args[1])
	fmt.Println("mybot ready, ^C exits")

	for {
		// read each incoming message
		m, err := getMessage(ws)
		if err != nil {
			log.Fatal(err)
		}

		// see if we're mentioned
		if m.Type == "message" && strings.HasPrefix(m.Text, "<@"+id+">") {
			// if so try to parse if
			parts := strings.Fields(m.Text)
			if len(parts) == 3 && parts[1] == "stock" {
				// looks good, get the quote and reply with the result
				go func(m Message) {
					m.Text = getQuote(parts[2])
					postMessage(ws, m)
				}(m)
				// NOTE: the Message object is copied, this is intentional
			} else if len(parts) == 3 && parts[1] == "meteo" {
				go func(m Message) {
					m.Text = getMeteo(parts[2])
					postMessage(ws, m)
				}(m)
			} else {
				// huh?
				m.Text = fmt.Sprintf("sorry, that does not compute\n")
				postMessage(ws, m)
			}
		}
	}
}

func getMeteo(sym string) string {
	sym = strings.ToUpper(sym)
	url := fmt.Sprintf("http://nominatim.openstreetmap.org/search/%s?format=json&addressdetails=1&limit=1", sym)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	defer resp.Body.Close()
	htmlData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	bytes := []byte(htmlData)
	var infos []Info
	err = json.Unmarshal(bytes, &infos)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	if infos[0].Type != "city" &&  infos[0].Type != "administrative" {
		return fmt.Sprintf(":x: \"%s\" n'est pas une ville", sym)
	} else {
		//fmt.Printf("Lon: %s et Lat: %s\n", infos[0].Lon, infos[0].Lat)
		return getCoord(infos[0].Lon, infos[0].Lat)
	}
}

func icon(icon string) string {
/*
clear-day, clear-night, rain, snow, sleet, wind, fog, cloudy,
partly-cloudy-day partly-cloudy-night hail, thunderstorm, or tornado
*/
	switch icon {
		case "clear-day":
			return ":sunny:"//
		case "clear-night":
			return ":sunny:"
		case "rain":
			return ":rain_cloud:"//
		case "snow":
			return ":snow_cloud:"//
		case "sleet":
			return ":snowflake: :fire:"//
		case "wind":
			return ":wind_blowing_face:"//
		case "fog":
			return ":foggy:"//
		case "cloudy":
			return ":cloud:"//
		case "partly-cloudy-day":
			return ":mostly_sunny:"//
		case "hail":
			return "hail"//
		case "thunderstorm":
			return ":thunder_cloud_and_rain:"//
		case "tornado":
			return ":tornado:"//
		case "partly-cloudy-night":
			return ":mostly_sunny:"//
		default:
			return icon
	}
}

func getCoord(Lon string, Lat string) string {
	time := int32(time.Now().Unix())
	url := fmt.Sprintf("https://api.darksky.net/forecast/%s/%s,%s,%d?&units=si&exclude=currently,minutely,hourly&lang=fr", os.Getenv("API_FORECAST"), Lat, Lon, time)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	defer resp.Body.Close()
	htmlData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	bytes := []byte(htmlData)
	var meteos Meteo
	err = json.Unmarshal(bytes, &meteos)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return fmt.Sprintf("%s :arrow_right: %s || Min: %.1f°C et Max: %.1f°C", icon(meteos.Daily.Data[0].Icon), meteos.Daily.Data[0].Summary, meteos.Daily.Data[0].TemperatureMin, meteos.Daily.Data[0].TemperatureMax)
}

// Get the quote via Yahoo. You should replace this method to something
// relevant to your team!
func getQuote(sym string) string {
	sym = strings.ToUpper(sym)
	url := fmt.Sprintf("http://download.finance.yahoo.com/d/quotes.csv?s=%s&f=nsl1op&e=.csv", sym)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	rows, err := csv.NewReader(resp.Body).ReadAll()
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	if len(rows) >= 1 && len(rows[0]) == 5 {
		return fmt.Sprintf("%s (%s) is trading at $%s", rows[0][0], rows[0][1], rows[0][2])
	}
	return fmt.Sprintf("unknown response format (symbol was \"%s\")", sym)
}
