package main

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestShowCurrentTime(t *testing.T) {
	var buf bytes.Buffer
	testTime := time.Date(2026, 3, 22, 12, 0, 0, 0, time.UTC)
	showCurrentTime(&buf, "Marmaris, Turkey", testTime)
	output := buf.String()
	if !strings.Contains(output, "Current time in Marmaris, Turkey") {
		t.Errorf("Output does not contain expected text: %s", output)
	}
}

func TestShowMoonInfo(t *testing.T) {
	var buf bytes.Buffer
	testTime := time.Date(2026, 3, 22, 12, 0, 0, 0, time.UTC)
	showMoonInfo(&buf, "Marmaris, Turkey", 36.8529, 28.2744, testTime)
	output := buf.String()
	if !strings.Contains(output, "Moon phase in Marmaris, Turkey") {
		t.Errorf("Output does not contain expected text: %s", output)
	}
	if !strings.Contains(output, "illumination") {
		t.Errorf("Output missing illumination: %s", output)
	}
}

func TestShowCurrentTimeMoscow(t *testing.T) {
	var buf bytes.Buffer
	testTime := time.Date(2026, 3, 22, 12, 0, 0, 0, time.UTC)
	showCurrentTime(&buf, "Moscow, Russia", testTime)
	output := buf.String()
	if !strings.Contains(output, "Current time in Moscow, Russia") {
		t.Errorf("Output does not contain expected text: %s", output)
	}
}

func TestShowPlanetsInfoNoPlanets(t *testing.T) {
	var buf bytes.Buffer
	// Use a time/location where no planets are visible, e.g., night time or polar
	testTime := time.Date(2026, 3, 22, 0, 0, 0, 0, time.UTC) // Midnight
	showPlanetsInfo(&buf, -90, 0, testTime) // South pole, might have no planets
	output := buf.String()
	// Depending on calculation, may or may not have planets, but test the output format
	if !strings.Contains(output, "локации") {
		t.Errorf("Output does not contain location text: %s", output)
	}
}