package main

import (
	"math"
	"testing"
	"time"
)

func TestNormalizeAngle(t *testing.T) {
	tests := []struct {
		input    float64
		expected float64
	}{
		{0, 0},
		{180, 180},
		{360, 0},
		{450, 90},
		{-90, 270},
		{720, 0},
		{-450, 270},
		{359.9, 359.9},
		{360.1, 0.1},
	}
	for _, test := range tests {
		result := normalizeAngle(test.input)
		if math.Abs(result-test.expected) > 1e-10 {
			t.Errorf("normalizeAngle(%f) = %f; want %f", test.input, result, test.expected)
		}
	}
}

func TestJulianDay(t *testing.T) {
	tests := []struct {
		year, month, day, hour, min, sec int
		expected                        float64
	}{
		{2000, 1, 1, 12, 0, 0, 2451545.0},
		{2000, 1, 1, 0, 0, 0, 2451544.5}, // 12 hours before
		{1999, 12, 31, 0, 0, 0, 2451543.5},
	}
	for _, test := range tests {
		testTime := time.Date(test.year, time.Month(test.month), test.day, test.hour, test.min, test.sec, 0, time.UTC)
		result := julianDay(testTime)
		if math.Abs(result-test.expected) > 1e-3 { // Allow small error
			t.Errorf("julianDay(%d-%d-%d %d:%d:%d) = %f; want ~%f", test.year, test.month, test.day, test.hour, test.min, test.sec, result, test.expected)
		}
	}
}

func TestMoonPhaseAndAltitude(t *testing.T) {
	tests := []struct {
		year, month, day int
		lat, lon        float64
	}{
		{2026, 3, 22, 36.8529, 28.2744},
		{2026, 3, 14, 36.8529, 28.2744},
	}
	for _, test := range tests {
		testTime := time.Date(test.year, time.Month(test.month), test.day, 12, 0, 0, 0, time.UTC)
		phase, illum, alt := moonPhaseAndAltitude(testTime, test.lat, test.lon)
		if phase == "" {
			t.Errorf("Phase should not be empty for %d-%d-%d", test.year, test.month, test.day)
		}
		if illum < 0 || illum > 1 {
			t.Errorf("Illumination %f out of range [0,1] for %d-%d-%d", illum, test.year, test.month, test.day)
		}
		_ = alt // Altitude varies
	}
}