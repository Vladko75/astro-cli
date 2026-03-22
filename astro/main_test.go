package main

import (
	"strings"
	"testing"
	"time"
)

func TestGetCurrentTimeData(t *testing.T) {
	// Disable colors for consistent test output
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