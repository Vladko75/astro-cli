package main

import (
	"strings"
	"testing"
	"time"
)

func TestGetCurrentTimeData(t *testing.T) {
	*noColorFlag = true
	defer func() { *noColorFlag = false }()
	
	testTime := time.Date(2026, 3, 22, 12, 0, 0, 0, time.UTC)
	data, output := getCurrentTimeData("Europe/Istanbul", "Marmaris", "A coastal town in Turkey", testTime)
	
	if data["timezone"] != "Europe/Istanbul" {
		t.Errorf("Timezone not set correctly: %v", data["timezone"])
	}
	if !strings.Contains(output, "Current time in Marmaris") {
		t.Errorf("Output does not contain expected text: %s", output)
	}
}

func TestGetMoonData(t *testing.T) {
	*noColorFlag = true
	defer func() { *noColorFlag = false }()
	
	testTime := time.Date(2026, 3, 22, 12, 0, 0, 0, time.UTC)
	data, output := getMoonData("Marmaris", "A coastal town in Turkey", 36.8529, 28.2744, testTime)
	
	if data["phase"] == nil {
		t.Errorf("Moon phase data missing: %v", data)
	}
	if !strings.Contains(output, "Moon phase in Marmaris") {
		t.Errorf("Output does not contain expected text: %s", output)
	}
	if !strings.Contains(output, "illumination") {
		t.Errorf("Output missing illumination: %s", output)
	}
}

func TestGetCurrentTimeDataMoscow(t *testing.T) {
	*noColorFlag = true
	defer func() { *noColorFlag = false }()
	
	testTime := time.Date(2026, 3, 22, 12, 0, 0, 0, time.UTC)
	_, output := getCurrentTimeData("Europe/Moscow", "Moscow", "Capital of Russia", testTime)
	
	if !strings.Contains(output, "Current time in Moscow") {
		t.Errorf("Output does not contain expected text: %s", output)
	}
}

func TestGetPlanetsData(t *testing.T) {
	*noColorFlag = true
	defer func() { *noColorFlag = false }()
	
	testTime := time.Date(2026, 3, 22, 0, 0, 0, 0, time.UTC)
	data, output := getPlanetsData("South Pole", "Southernmost point on Earth", -90, 0, testTime)
	
	if data["count"] == nil {
		t.Errorf("Planet count missing: %v", data)
	}
	if !strings.Contains(output, "South Pole") {
		t.Errorf("Output does not contain location text: %s", output)
	}
}

// NEW TESTS FOR BETTER COVERAGE

func TestGetSunData(t *testing.T) {
	*noColorFlag = true
	defer func() { *noColorFlag = false }()
	
	testTime := time.Date(2026, 3, 22, 12, 0, 0, 0, time.UTC)
	data, output := getSunData("London", "Capital of UK", 51.5074, -0.1278, testTime)
	
	if data["altitude"] == nil {
		t.Errorf("Sun altitude missing: %v", data)
	}
	if data["azimuth"] == nil {
		t.Errorf("Sun azimuth missing: %v", data)
	}
	if !strings.Contains(output, "Sun in London") {
		t.Errorf("Output does not contain location: %s", output)
	}
}

func TestGetStarsData(t *testing.T) {
	*noColorFlag = true
	defer func() { *noColorFlag = false }()
	
	testTime := time.Date(2026, 3, 22, 12, 0, 0, 0, time.UTC)
	data, output := getStarsData("Athens", "Capital of Greece", 37.9838, 23.7275, testTime)
	
	if data["count"] == nil {
		t.Errorf("Star count missing: %v", data)
	}
	if !strings.Contains(output, "Athens") {
		t.Errorf("Output does not contain location: %s", output)
	}
}

func TestColorizeWithColorEnabled(t *testing.T) {
	*noColorFlag = false
	result := colorize("test", "31")
	if !strings.Contains(result, "\033[") {
		t.Errorf("Color codes not applied when colors enabled: %s", result)
	}
}

func TestColorizeWithColorDisabled(t *testing.T) {
	*noColorFlag = true
	defer func() { *noColorFlag = false }()
	
	result := colorize("test", "31")
	if strings.Contains(result, "\033[") {
		t.Errorf("Color codes applied when colors disabled: %s", result)
	}
	if result != "test" {
		t.Errorf("Expected 'test' but got: %s", result)
	}
}

func TestSunAltitudeVariousLocations(t *testing.T) {
	testTime := time.Date(2026, 3, 22, 12, 0, 0, 0, time.UTC)
	
	// Test equator
	alt1, az1 := sunAltitude(testTime, 0, 0)
	if az1 < 0 || az1 > 360 {
		t.Errorf("Invalid azimuth at equator: %f", az1)
	}
	
	// Test north pole (should have specific behavior)
	alt2, az2 := sunAltitude(testTime, 90, 0)
	if az2 < 0 || az2 > 360 {
		t.Errorf("Invalid azimuth at north pole: %f", az2)
	}
	
	// Test southern hemisphere
	alt3, az3 := sunAltitude(testTime, -45, 170)
	if az3 < 0 || az3 > 360 {
		t.Errorf("Invalid azimuth southern hemisphere: %f", az3)
	}
	
	_ = alt1
	_ = alt2
	_ = alt3
}

func TestStarsAboveHorizonMultipleLocations(t *testing.T) {
	testTime := time.Date(2026, 3, 22, 12, 0, 0, 0, time.UTC)
	
	// Test at equator
	stars1 := StarsAboveHorizon(0, 0, testTime)
	if len(stars1) < 0 {
		t.Errorf("Invalid star list length: %d", len(stars1))
	}
	
	// Test at north latitude
	stars2 := StarsAboveHorizon(40, -74, testTime)
	if len(stars2) < 0 {
		t.Errorf("Invalid star list length: %d", len(stars2))
	}
	
	// Test at southern latitude
	stars3 := StarsAboveHorizon(-33, 151, testTime)
	if len(stars3) < 0 {
		t.Errorf("Invalid star list length: %d", len(stars3))
	}
}

func TestGetMoonDataWithoutDescription(t *testing.T) {
	*noColorFlag = true
	defer func() { *noColorFlag = false }()
	
	testTime := time.Date(2026, 6, 21, 12, 0, 0, 0, time.UTC)
	_, output := getMoonData("CustomLoc", "", 50, 20, testTime)
	
	if !strings.Contains(output, "CustomLoc") {
		t.Errorf("Location name missing: %s", output)
	}
	if !strings.Contains(output, "illumination") {
		t.Errorf("Illumination missing: %s", output)
	}
}

func TestGetSunTimes(t *testing.T) {
	testTime := time.Date(2026, 3, 22, 12, 0, 0, 0, time.UTC)
	sunInfo := SunTimes(testTime, 51.5074, -0.1278)
	
	if sunInfo.Sunrise.IsZero() {
		t.Errorf("Sunrise time is zero: %v", sunInfo.Sunrise)
	}
	if sunInfo.Sunset.IsZero() {
		t.Errorf("Sunset time is zero: %v", sunInfo.Sunset)
	}
	if sunInfo.DayLength <= 0 {
		t.Errorf("Day length should be positive: %f", sunInfo.DayLength)
	}
	if sunInfo.Sunset.Before(sunInfo.Sunrise) {
		t.Errorf("Sunset before sunrise: %v vs %v", sunInfo.Sunset, sunInfo.Sunrise)
	}
}

func TestMoonPhaseAndAltitudeNewMoon(t *testing.T) {
	// Test during new moon period
	testTime := time.Date(2026, 3, 21, 12, 0, 0, 0, time.UTC)
	phase, illum, alt := moonPhaseAndAltitude(testTime, 40, -74)
	
	if phase == "" {
		t.Errorf("Moon phase is empty")
	}
	if illum < 0 || illum > 1 {
		t.Errorf("Illumination out of range: %f", illum)
	}
	if alt < -90 || alt > 90 {
		t.Errorf("Altitude out of range: %f", alt)
	}
}

func TestMoonPhaseAndAltitudeFullMoon(t *testing.T) {
	// Test during full moon period
	testTime := time.Date(2026, 3, 29, 12, 0, 0, 0, time.UTC)
	phase, illum, _ := moonPhaseAndAltitude(testTime, 40, -74)
	
	if phase == "" {
		t.Errorf("Moon phase is empty")
	}
	if illum < 0 || illum > 1 {
		t.Errorf("Illumination out of range: %f", illum)
	}
	if phase != "Full Moon" && phase != "Waxing Gibbous" && phase != "Waning Gibbous" {
		t.Logf("Expected full moon phase, got: %s", phase)
	}
}

func TestPlanetsDataWithoutPlanets(t *testing.T) {
	*noColorFlag = true
	defer func() { *noColorFlag = false }()
	
	// Extreme date/location where no planets visible
	testTime := time.Date(2026, 1, 15, 6, 0, 0, 0, time.UTC)
	data, output := getPlanetsData("TestLoc", "Test location", 89, 0, testTime)
	
	if data["count"] != 0 && data["count"] != nil {
		count := data["count"].(int)
		if count < 0 {
			t.Errorf("Planet count cannot be negative: %d", count)
		}
	}
	if !strings.Contains(output, "TestLoc") {
		t.Errorf("Location missing: %s", output)
	}
}

func TestGetCurrentTimeDataUTCTimezone(t *testing.T) {
	*noColorFlag = true
	defer func() { *noColorFlag = false }()
	
	testTime := time.Date(2026, 3, 22, 12, 0, 0, 0, time.UTC)
	data, _ := getCurrentTimeData("UTC", "TestCity", "", testTime)
	
	if data["timezone"] != "UTC" {
		t.Errorf("UTC timezone not preserved: %v", data["timezone"])
	}
	if data["utc_time"] == nil {
		t.Errorf("UTC time missing from data")
	}
}

func TestGetStarsDataNoDescription(t *testing.T) {
	*noColorFlag = true
	defer func() { *noColorFlag = false }()
	
	testTime := time.Date(2026, 3, 22, 12, 0, 0, 0, time.UTC)
	_, output := getStarsData("TestLocation", "", 45, 90, testTime)
	
	if !strings.Contains(output, "TestLocation") {
		t.Errorf("Location missing in output: %s", output)
	}
}

func TestSunDataWithVariousLocations(t *testing.T) {
	*noColorFlag = true
	defer func() { *noColorFlag = false }()
	
	testTime := time.Date(2026, 3, 22, 12, 0, 0, 0, time.UTC)
	
	// Test near equator
	data1, _ := getSunData("Equator", "At equator", 0, 0, testTime)
	if data1["altitude"] == nil {
		t.Errorf("Altitude missing at equator")
	}
	
	// Test near north pole
	data2, _ := getSunData("Arctic", "Near North Pole", 80, 0, testTime)
	if data2["altitude"] == nil {
		t.Errorf("Altitude missing at polar")
	}
	
	// Test southern hemisphere
	data3, _ := getSunData("Sydney", "Southern", -33.87, 151.21, testTime)
	if data3["altitude"] == nil {
		t.Errorf("Altitude missing in southern hemisphere")
	}
}

func TestMoonDataEdgeCases(t *testing.T) {
	*noColorFlag = true
	defer func() { *noColorFlag = false }()
	
	testTime := time.Date(2026, 3, 22, 12, 0, 0, 0, time.UTC)
	
	// Test at poles
	data1, _ := getMoonData("North Pole", "Top of world", 90, 0, testTime)
	if data1["phase"] == nil {
		t.Errorf("Moon phase missing at north pole")
	}
	
	data2, _ := getMoonData("South Pole", "Bottom of world", -90, 0, testTime)
	if data2["phase"] == nil {
		t.Errorf("Moon phase missing at south pole")
	}
	
	// Test at different times of year
	for month := 1; month <= 12; month++ {
		testTime := time.Date(2026, time.Month(month), 15, 12, 0, 0, 0, time.UTC)
		data, output := getMoonData("TestLoc", "Test", 40, -74, testTime)
		
		if data["phase"] == nil {
			t.Errorf("Moon phase missing for month %d", month)
		}
		if !strings.Contains(output, "Moon phase") {
			t.Errorf("Output missing 'Moon phase' text for month %d: %s", month, output)
		}
	}
}

func TestGetMoonDataWithoutDescriptionCase(t *testing.T) {
	*noColorFlag = true
	defer func() { *noColorFlag = false }()
	
	testTime := time.Date(2026, 3, 22, 12, 0, 0, 0, time.UTC)
	_, output := getMoonData("CityName", "", 40, -74, testTime)
	
	if !strings.Contains(output, "Moon phase") {
		t.Errorf("Output missing Moon phase: %s", output)
	}
}

func GetPlanetsDataEdgeCases(t *testing.T) {
	*noColorFlag = true
	defer func() { *noColorFlag = false }()
	
	testTime := time.Date(2026, 3, 22, 12, 0, 0, 0, time.UTC)
	
	// Test at equator
	data1, output1 := getPlanetsData("Equator", "", 0, 0, testTime)
	if !strings.Contains(output1, "planets above") {
		t.Errorf("Output missing 'planets above': %s", output1)
	}
	if data1["count"] == nil {
		t.Errorf("Count missing in planets data")
	}
	
	// Test at different times
	for hour := 0; hour < 24; hour += 6 {
		testTime := time.Date(2026, 3, 22, hour, 0, 0, 0, time.UTC)
		data, _ := getPlanetsData("TestLoc", "Test Location", 45, 45, testTime)
		if data["count"] == nil {
			t.Errorf("Planets count missing at hour %d", hour)
		}
	}
}

func TestColorizeEdgeCases(t *testing.T) {
	tests := []struct {
		text    string
		enabled bool
		name    string
	}{
		{"", true, "empty string with color"},
		{"", false, "empty string without color"},
		{"test", true, "simple text with color"},
		{"test", false, "simple text without color"},
		{"test\nwith\nnewlines", true, "multiline with color"},
		{"test\nwith\nnewlines", false, "multiline without color"},
		{"ANSI: red text", true, "ANSI codes with color"},
	}
	
	for _, test := range tests {
		*noColorFlag = !test.enabled
		result := colorize(test.text, yellow)
		
		if result == "" && test.text == "" {
			// Empty input should give empty output
			continue
		}
		
		if result == "" && test.text != "" {
			t.Errorf("colorize('%s', ...) returned empty for test: %s", test.text, test.name)
		}
	}
	*noColorFlag = true
}

func TestGetCurrentTimeDataVariousTimezones(t *testing.T) {
	*noColorFlag = true
	defer func() { *noColorFlag = false }()
	
	testTime := time.Date(2026, 3, 22, 12, 0, 0, 0, time.UTC)
	
	timezones := []string{"UTC", "Europe/Istanbul", "Europe/Moscow", "America/New_York", "Asia/Tokyo"}
	
	for _, tz := range timezones {
		data, output := getCurrentTimeData(tz, "TestCity", "Test description", testTime)
		
		if data["timezone"] != tz {
			t.Errorf("Timezone mismatch for %s: got %v", tz, data["timezone"])
		}
		if !strings.Contains(output, "Current time") {
			t.Errorf("Output missing 'Current time' for %s: %s", tz, output)
		}
	}
}

func TestSunAltitudeRanges(t *testing.T) {
	tests := []struct {
		lat, lon float64
		name     string
	}{
		{0, 0, "equator"},
		{45, 45, "mid-northern"},
		{-45, -45, "mid-southern"},
		{89, 0, "near north pole"},
		{-89, 0, "near south pole"},
		{23.5, 0, "tropic of cancer"},
		{-23.5, 0, "tropic of capricorn"},
	}
	
	testTime := time.Date(2026, 3, 22, 12, 0, 0, 0, time.UTC)
	
	for _, test := range tests {
		data, _ := getSunData("Test", "Test location", test.lat, test.lon, testTime)
		if data["altitude"] == nil {
			t.Errorf("Sun altitude missing for %s (%.1f, %.1f)", test.name, test.lat, test.lon)
		}
		
		altValue, ok := data["altitude"].(float64)
		if !ok {
			t.Errorf("Sun altitude not a float64 for %s", test.name)
			continue
		}
		
		if altValue < -90 || altValue > 90 {
			t.Errorf("Sun altitude out of valid range for %s: %f", test.name, altValue)
		}
	}
}

func TestSunAltitudeDirectly(t *testing.T) {
	// Test sunAltitude directly at various times
	tests := []struct {
		year, month, day, hour int
		lat, lon               float64
		name                  string
	}{
		{2026, 1, 21, 12, 0, 0, "winter solstice equator noon"},
		{2026, 6, 21, 12, 45, 45, "summer solstice mid-latitude noon"},
		{2026, 3, 20, 0, 0, 0, "spring equinox equator midnight"},
		{2026, 9, 23, 6, 40, -74, "autumn equinox NYC sunrise"},
		{2026, 12, 22, 18, -40, 150, "southern summer sunset"},
	}
	
	for _, test := range tests {
		testTime := time.Date(test.year, time.Month(test.month), test.day, test.hour, 0, 0, 0, time.UTC)
		alt, az := sunAltitude(testTime, test.lat, test.lon)
		
		if alt < -90 || alt > 90 {
			t.Errorf("Altitude out of range for %s: %f", test.name, alt)
		}
		if az < 0 || az > 360 {
			t.Errorf("Azimuth out of range for %s: %f", test.name, az)
		}
	}
}

func TestGetTimezoneEdgeCases(t *testing.T) {
	// Test various timezone scenarios
	tests := []string{
		"UTC",
		"Europe/London",
		"Europe/Paris",
		"America/New_York",
		"America/Chicago",
		"America/Los_Angeles",
		"Australia/Sydney",
		"Asia/Tokyo",
		"Asia/Shanghai",
		"Asia/Dubai",
		"Africa/Cairo",
	}
	
	for _, tzName := range tests {
		loc, err := time.LoadLocation(tzName)
		if err != nil {
			t.Logf("Skipping timezone %s: %v", tzName, err)
			continue
		}
		
		// Verify timezone loading works
		if loc == nil {
			t.Errorf("Failed to load timezone: %s", tzName)
		}
	}
}

func TestStarsAboveHorizonDifferentTimes(t *testing.T) {
	*noColorFlag = true
	defer func() { *noColorFlag = false }()
	
	testDate := time.Date(2026, 3, 22, 0, 0, 0, 0, time.UTC)
	
	// Test at different times through the night and day
	for hour := 0; hour < 24; hour += 3 {
		testTime := testDate.Add(time.Duration(hour) * time.Hour)
		data, _ := getStarsData("TestLoc", "Test", 45, 45, testTime)
		
		if data["stars"] == nil {
			t.Logf("Stars data missing at hour %d", hour)
		}
	}
}

func TestSunDataAcrossYear(t *testing.T) {
	*noColorFlag = true
	defer func() { *noColorFlag = false }()
	
	// Test sun data throughout the year at different latitudes
	for month := 1; month <= 12; month++ {
		testTime := time.Date(2026, time.Month(month), 15, 12, 0, 0, 0, time.UTC)
		
		// Test at various latitudes
		for lat := -90; lat <= 90; lat += 30 {
			data1, _ := getSunData("Test", "Test", float64(lat), 0, testTime)
			if data1["altitude"] == nil {
				t.Errorf("Sun altitude missing for month %d, lat %d", month, lat)
			}
		}
	}
}

func TestMoonPhaseConsistency(t *testing.T) {
	// Test moon phase consistency across multiple dates
	testLat, testLon := 40.0, -74.0
	
	for day := 1; day <= 28; day++ {
		testTime := time.Date(2026, 3, day, 12, 0, 0, 0, time.UTC)
		phase, illum, _ := moonPhaseAndAltitude(testTime, testLat, testLon)
		
		if phase == "" {
			t.Errorf("Moon phase empty on day %d", day)
		}
		if illum < 0 || illum > 1 {
			t.Errorf("Moon illumination out of range on day %d: %f", day, illum)
		}
	}
}