package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

var (
	cityFlag     = flag.String("city", "", "City name (e.g., tokyo, london, etc.)")
	latFlag      = flag.Float64("lat", 0, "Latitude in degrees")
	lonFlag      = flag.Float64("lon", 0, "Longitude in degrees")
	timeFlag     = flag.String("time", "", "Time in RFC3339 format (e.g., 2026-03-22T12:00:00Z), default is now")
	showFlag     = flag.String("show", "all", "What to show: all, time, moon, planets (comma-separated)")
)

func main() {
	flag.Parse()

	// Parse show options
	showOptions := make(map[string]bool)
	for _, opt := range strings.Split(*showFlag, ",") {
		showOptions[strings.TrimSpace(opt)] = true
	}
	if *showFlag == "all" {
		showOptions["time"] = true
		showOptions["moon"] = true
		showOptions["planets"] = true
	}

	// Determine location
	var lat, lon float64
	var cityName string
	var description string
	if *cityFlag != "" {
		city := strings.ToLower(*cityFlag)
		if loc, ok := cities[city]; ok {
			lat, lon = loc.Lat, loc.Lon
			description = loc.Description
			cityName = strings.Title(city)
		} else {
			fmt.Printf("Unknown city: %s\n", *cityFlag)
			return
		}
	} else if *latFlag != 0 || *lonFlag != 0 {
		lat, lon = *latFlag, *lonFlag
		cityName = fmt.Sprintf("Custom location (%.4f, %.4f)", lat, lon)
		description = ""
	} else {
		// Default to Marmaris
		lat, lon = cities["marmaris"].Lat, cities["marmaris"].Lon
		description = cities["marmaris"].Description
		cityName = "Marmaris"
	}

	// Determine time
	var t time.Time
	if *timeFlag != "" {
		var err error
		t, err = time.Parse(time.RFC3339, *timeFlag)
		if err != nil {
			fmt.Printf("Invalid time format: %s\n", *timeFlag)
			return
		}
	} else {
		t = time.Now()
	}

	// Show time
	if showOptions["time"] {
		showCurrentTime(os.Stdout, cityName, description, t)
	}

	// Show moon
	if showOptions["moon"] {
		showMoonInfo(os.Stdout, cityName, description, lat, lon, t)
	}

	// Show planets
	if showOptions["planets"] {
		showPlanetsInfo(os.Stdout, cityName, description, lat, lon, t)
	}
}

func showCurrentTime(w io.Writer, cityName string, description string, t time.Time) {
	var loc *time.Location
	var err error
	switch {
	case strings.Contains(cityName, "Marmaris"):
		loc, err = time.LoadLocation("Europe/Istanbul")
	case strings.Contains(cityName, "Moscow"):
		loc, err = time.LoadLocation("Europe/Moscow")
	default:
		loc = time.UTC
	}
	if err != nil {
		fmt.Fprintf(w, "Error loading location for %s: %v\n", cityName, err)
		return
	}
	localTime := t.In(loc)
	if description != "" {
		fmt.Fprintf(w, "Current time in %s (%s): %s\n", cityName, description, localTime.Format("2006-01-02 15:04:05 MST"))
	} else {
		fmt.Fprintf(w, "Current time in %s: %s\n", cityName, localTime.Format("2006-01-02 15:04:05 MST"))
	}
}

func showMoonInfo(w io.Writer, cityName string, description string, lat, lon float64, t time.Time) {
	phase, illum, alt := moonPhaseAndAltitude(t, lat, lon)
	if description != "" {
		fmt.Fprintf(w, "Moon phase in %s (%s): %s (illumination %.1f%%), altitude above horizon: %.2f°\n", cityName, description, phase, illum*100, alt)
	} else {
		fmt.Fprintf(w, "Moon phase in %s: %s (illumination %.1f%%), altitude above horizon: %.2f°\n", cityName, phase, illum*100, alt)
	}
}

func showPlanetsInfo(w io.Writer, cityName string, description string, lat, lon float64, t time.Time) {
	planets := planetsAboveHorizon(lat, lon, t)
	if len(planets) == 0 {
		if description != "" {
			fmt.Fprintf(w, "Над горизонтом в %s (%s) сейчас нет известных планет (по приближённым расчётам).\n", cityName, description)
		} else {
			fmt.Fprintf(w, "Над горизонтом в данной локации (%.4f, %.4f) сейчас нет известных планет (по приближённым расчётам).\n", lat, lon)
		}
		return
	}
	if description != "" {
		fmt.Fprintf(w, "Планеты над горизонтом в %s (%s):\n", cityName, description)
	} else {
		fmt.Fprintf(w, "Планеты над горизонтом в локации (%.4f, %.4f):\n", lat, lon)
	}
	for _, p := range planets {
		fmt.Fprintf(w, " - %s: %.2f°\n", p.Name, p.Altitude)
	}
}
