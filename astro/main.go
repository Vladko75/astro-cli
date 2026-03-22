package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

const Version = "1.2.0"

var (
	cityFlag     = flag.String("city", "", "City name (e.g., tokyo, london, etc.)")
	latFlag      = flag.Float64("lat", 0, "Latitude in degrees (-90 to 90)")
	lonFlag      = flag.Float64("lon", 0, "Longitude in degrees (-180 to 180)")
	timeFlag     = flag.String("time", "", "Time in RFC3339 format (e.g., 2026-03-22T12:00:00Z), default is now")
	showFlag     = flag.String("show", "all", "What to show: all, time, moon, planets, stars, sun (comma-separated)")
	jsonFlag     = flag.Bool("json", false, "Output in JSON format")
	noColorFlag  = flag.Bool("no-color", false, "Disable colored output")
	helpFlag     = flag.Bool("help", false, "Show help message")
	versionFlag  = flag.Bool("version", false, "Show version")
)

func main() {
	flag.Parse()

	// Handle help and version
	if *helpFlag {
		showHelp()
		return
	}
	if *versionFlag {
		fmt.Printf("astro version %s\n", Version)
		return
	}

	// Validate coordinates
	if *latFlag != 0 || *lonFlag != 0 {
		if *latFlag < -90 || *latFlag > 90 {
			fmt.Printf("Error: Latitude must be between -90 and 90 degrees\n")
			os.Exit(1)
		}
		if *lonFlag < -180 || *lonFlag > 180 {
			fmt.Printf("Error: Longitude must be between -180 and 180 degrees\n")
			os.Exit(1)
		}
	}

	// Parse show options
	showOptions := make(map[string]bool)
	validShowOptions := map[string]bool{
		"time": true, "moon": true, "planets": true, "sun": true, "stars": true,
	}
	if *showFlag == "all" {
		showOptions["time"] = true
		showOptions["moon"] = true
		showOptions["planets"] = true
		showOptions["sun"] = true
		showOptions["stars"] = false // stars not shown by default
	} else {
		for _, opt := range strings.Split(*showFlag, ",") {
			trimmed := strings.TrimSpace(strings.ToLower(opt))
			if trimmed == "" {
				continue
			}
			if !validShowOptions[trimmed] {
				fmt.Printf("Error: Invalid show option '%s'. Valid options: time, moon, planets, sun, stars\n", trimmed)
				os.Exit(1)
			}
			showOptions[trimmed] = true
		}
	}

	// Determine location and get timezone
	var lat, lon float64
	var cityName string
	var description string
	var tz string

	if *cityFlag != "" {
		city := strings.ToLower(*cityFlag)
		if loc, ok := cities[city]; ok {
			lat, lon = loc.Lat, loc.Lon
			description = loc.Description
			tz = loc.Timezone
			cityName = strings.Title(city)
		} else {
			fmt.Printf("Error: Unknown city: %s\n", *cityFlag)
			os.Exit(1)
		}
	} else if *latFlag != 0 || *lonFlag != 0 {
		lat, lon = *latFlag, *lonFlag
		cityName = fmt.Sprintf("Custom location (%.4f, %.4f)", lat, lon)
		description = ""
		tz = "UTC"
	} else {
		// Default to Marmaris
		lat, lon = cities["marmaris"].Lat, cities["marmaris"].Lon
		description = cities["marmaris"].Description
		tz = cities["marmaris"].Timezone
		cityName = "Marmaris"
	}

	// Determine time
	var t time.Time
	if *timeFlag != "" {
		var err error
		t, err = time.Parse(time.RFC3339, *timeFlag)
		if err != nil {
			fmt.Printf("Error: Invalid time format: %s\n", *timeFlag)
			os.Exit(1)
		}
	} else {
		t = time.Now()
	}

	// Prepare output data
	output := make(map[string]interface{})
	output["location"] = cityName
	if description != "" {
		output["description"] = description
	}
	output["coordinates"] = map[string]float64{"lat": lat, "lon": lon}
	output["timezone"] = tz
	output["time"] = t.Format(time.RFC3339)

	// Collect data for each requested option
	if showOptions["time"] {
		timeData, timeStr := getCurrentTimeData(tz, cityName, description, t)
		output["current_time"] = timeData
		if !*jsonFlag {
			fmt.Fprint(os.Stdout, timeStr)
		}
	}

	if showOptions["sun"] {
		sunData, sunStr := getSunData(cityName, description, lat, lon, t)
		output["sun"] = sunData
		if !*jsonFlag {
			fmt.Fprint(os.Stdout, sunStr)
		}
	}

	if showOptions["moon"] {
		moonData, moonStr := getMoonData(cityName, description, lat, lon, t)
		output["moon"] = moonData
		if !*jsonFlag {
			fmt.Fprint(os.Stdout, moonStr)
		}
	}

	if showOptions["planets"] {
		planetsData, planetsStr := getPlanetsData(cityName, description, lat, lon, t)
		output["planets"] = planetsData
		if !*jsonFlag {
			fmt.Fprint(os.Stdout, planetsStr)
		}
	}

	if showOptions["stars"] {
		starsData, starsStr := getStarsData(cityName, description, lat, lon, t)
		output["stars"] = starsData
		if !*jsonFlag {
			fmt.Fprint(os.Stdout, starsStr)
		}
	}

	// Output JSON if requested
	if *jsonFlag {
		jsonOutput, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			fmt.Printf("Error encoding JSON: %v\n", err)
			return
		}
		fmt.Println(string(jsonOutput))
	}
}

func showHelp() {
	help := `Astro CLI - Astronomical Information Tool

Usage: astro [options]

Options:
  -city string      City name (e.g., tokyo, london). Default: marmaris
  -lat float        Latitude in degrees (-90 to 90)
  -lon float        Longitude in degrees (-180 to 180)
  -time string      Time in RFC3339 format (e.g., 2026-03-22T12:00:00Z). Default: current time
  -show string      What to show: all, time, sun, moon, planets, stars (comma-separated). Default: all
  -json             Output in JSON format
  -no-color         Disable colored output
  -help             Show this help message
  -version          Show version number

Examples:
  astro --city=tokyo --show=time
  astro --lat=40.7128 --lon=-74.0060 --show=moon,planets
  astro --city=paris --show=all --json
  astro --city=london --show=stars
`
	fmt.Print(help)
}

// Helper functions for color output
func colorize(text, colorCode string) string {
	if *noColorFlag {
		return text
	}
	return fmt.Sprintf("\033[%sm%s\033[0m", colorCode, text)
}

const (
	// ANSI color codes
	redbold   = "1;31"
	yellow    = "33"
	cyan      = "36"
	green     = "32"
	magenta   = "35"
)

// Planet icons
var planetIcons = map[string]string{
	"Mercury": "☿",
	"Venus":   "♀",
	"Earth":   "⊕",
	"Mars":    "♂",
	"Jupiter": "♃",
	"Saturn":  "♄",
	"Uranus":  "♅",
	"Neptune": "♆",
}

// Data extraction functions
func getCurrentTimeData(tzStr, cityName, description string, t time.Time) (map[string]interface{}, string) {
	tz := getTimezone(tzStr)
	localTime := t.In(tz)

	data := map[string]interface{}{
		"utc_time":   t.Format(time.RFC3339),
		"local_time": localTime.Format("2006-01-02 15:04:05"),
		"timezone":   tzStr,
	}

	var output string
	timeStr := localTime.Format("2006-01-02 15:04:05 MST")
	if description != "" {
		output = fmt.Sprintf("%s Current time in %s (%s): %s\n",
			colorize("📍", cyan), colorize(cityName, cyan), description, colorize(timeStr, green))
	} else {
		output = fmt.Sprintf("%s Current time in %s: %s\n",
			colorize("📍", cyan), colorize(cityName, cyan), colorize(timeStr, green))
	}

	return data, output
}

func getSunData(cityName, description string, lat, lon float64, t time.Time) (map[string]interface{}, string) {
	sunInfo := SunTimes(t, lat, lon)
	alt, azimuth := sunAltitude(t, lat, lon)

	data := map[string]interface{}{
		"altitude":    alt,
		"azimuth":     azimuth,
		"rise_time":   sunInfo.Sunrise.Format("15:04:05"),
		"set_time":    sunInfo.Sunset.Format("15:04:05"),
		"day_length":  fmt.Sprintf("%.1f hours", sunInfo.DayLength),
	}

	status := "below"
	if alt > 0 {
		status = colorize(fmt.Sprintf("%.2f°", alt), yellow)
	}

	var output string
	if description != "" {
		output = fmt.Sprintf("%s Sun in %s (%s): altitude %s, azimuth %.0f°, sunrise %s, sunset %s\n",
			colorize("☀", yellow), colorize(cityName, yellow), description, status, azimuth,
			sunInfo.Sunrise.Format("15:04"), sunInfo.Sunset.Format("15:04"))
	} else {
		output = fmt.Sprintf("%s Sun at %s: altitude %s, azimuth %.0f°, sunrise %s, sunset %s\n",
			colorize("☀", yellow), colorize(cityName, yellow), status, azimuth,
			sunInfo.Sunrise.Format("15:04"), sunInfo.Sunset.Format("15:04"))
	}

	return data, output
}

func getMoonData(cityName, description string, lat, lon float64, t time.Time) (map[string]interface{}, string) {
	phase, illum, alt := moonPhaseAndAltitude(t, lat, lon)

	data := map[string]interface{}{
		"phase":         phase,
		"illumination":  fmt.Sprintf("%.1f%%", illum*100),
		"altitude":      alt,
	}

	phaseEmoji := "🌙"
	if illum < 0.1 {
		phaseEmoji = "🌑"
	} else if illum < 0.3 {
		phaseEmoji = "🌘"
	} else if illum < 0.5 {
		phaseEmoji = "🌗"
	} else if illum < 0.7 {
		phaseEmoji = "🌖"
	} else if illum < 0.9 {
		phaseEmoji = "🌕"
	} else {
		phaseEmoji = "🌑"
	}

	var output string
	if description != "" {
		output = fmt.Sprintf("%s Moon phase in %s (%s): %s (illumination %.1f%%), altitude: %.2f°\n",
			colorize(phaseEmoji, magenta), colorize(cityName, magenta), description,
			colorize(phase, magenta), illum*100, alt)
	} else {
		output = fmt.Sprintf("%s Moon phase in %s: %s (illumination %.1f%%), altitude: %.2f°\n",
			colorize(phaseEmoji, magenta), colorize(cityName, magenta),
			colorize(phase, magenta), illum*100, alt)
	}

	return data, output
}

func getPlanetsData(cityName, description string, lat, lon float64, t time.Time) (map[string]interface{}, string) {
	planets := planetsAboveHorizon(lat, lon, t)

	var planetsList []map[string]interface{}
	for _, p := range planets {
		planetsList = append(planetsList, map[string]interface{}{
			"name":     p.Name,
			"altitude": p.Altitude,
		})
	}

	data := map[string]interface{}{
		"planets": planetsList,
		"count":   len(planets),
	}

	var output string
	if len(planets) == 0 {
		if description != "" {
			output = fmt.Sprintf("%s No known planets above the horizon in %s (%s).\n",
				colorize("🪐", green), colorize(cityName, green), description)
		} else {
			output = fmt.Sprintf("%s No known planets above the horizon at %s.\n",
				colorize("🪐", green), colorize(cityName, green))
		}
	} else {
		header := fmt.Sprintf("%s Planets above the horizon in %s", colorize("🪐", green), colorize(cityName, green))
		if description != "" {
			header += fmt.Sprintf(" (%s)", description)
		}
		output = header + ":\n"

		for _, p := range planets {
			icon := planetIcons[p.Name]
			if icon == "" {
				icon = "•"
			}
			output += fmt.Sprintf("  %s %s: %.2f°\n", colorize(icon, green), p.Name, p.Altitude)
		}
	}

	return data, output
}

func getStarsData(cityName, description string, lat, lon float64, t time.Time) (map[string]interface{}, string) {
	stars := StarsAboveHorizon(lat, lon, t)

	var starsList []map[string]interface{}
	for _, s := range stars {
		starsList = append(starsList, map[string]interface{}{
			"name":      s.Name,
			"altitude":  s.Altitude,
			"magnitude": s.Magnitude,
		})
	}

	data := map[string]interface{}{
		"stars": starsList,
		"count": len(stars),
	}

	var output string
	if len(stars) == 0 {
		if description != "" {
			output = fmt.Sprintf("%s No bright stars above the horizon in %s (%s).\n",
				colorize("★", cyan), colorize(cityName, cyan), description)
		} else {
			output = fmt.Sprintf("%s No bright stars above the horizon at %s.\n",
				colorize("★", cyan), colorize(cityName, cyan))
		}
	} else {
		header := fmt.Sprintf("%s Bright stars above the horizon in %s", colorize("★", cyan), colorize(cityName, cyan))
		if description != "" {
			header += fmt.Sprintf(" (%s)", description)
		}
		output = header + ":\n"

		for _, s := range stars {
			output += fmt.Sprintf("  ✦ %s (Mag %.2f): %.2f°\n", colorize(s.Name, cyan), s.Magnitude, s.Altitude)
		}
	}

	return data, output
}
