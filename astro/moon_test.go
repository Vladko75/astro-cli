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

// Additional moon tests for better coverage
func TestMoonAltitudeAtDifferentLatitudes(t *testing.T) {
	testTime := time.Date(2026, 3, 22, 12, 0, 0, 0, time.UTC)
	
	// Test at equator
	_, _, alt1 := moonPhaseAndAltitude(testTime, 0, 0)
	if alt1 < -90 || alt1 > 90 {
		t.Errorf("Moon altitude out of range at equator: %f", alt1)
	}
	
	// Test at north pole
	_, _, alt2 := moonPhaseAndAltitude(testTime, 90, 0)
	if alt2 < -90 || alt2 > 90 {
		t.Errorf("Moon altitude out of range at north pole: %f", alt2)
	}
	
	// Test at south pole
	_, _, alt3 := moonPhaseAndAltitude(testTime, -90, 0)
	if alt3 < -90 || alt3 > 90 {
		t.Errorf("Moon altitude out of range at south pole: %f", alt3)
	}
	
	// Test at mid-northern hemisphere
	_, _, alt4 := moonPhaseAndAltitude(testTime, 45, 45)
	if alt4 < -90 || alt4 > 90 {
		t.Errorf("Moon altitude out of range at mid-latitude: %f", alt4)
	}
}

func TestMoonPhaseVariations(t *testing.T) {
	lat, lon := 40.0, -74.0
	
	tests := []struct {
		month   int
		dayDiff int // Days to test in the month
	}{
		{3, 7},
		{3, 14},
		{3, 21},
		{3, 28},
	}
	
	for _, test := range tests {
		testTime := time.Date(2026, time.Month(test.month), test.dayDiff, 12, 0, 0, 0, time.UTC)
		phase, illum, _ := moonPhaseAndAltitude(testTime, lat, lon)
		
		if phase == "" {
			t.Errorf("Phase empty for day %d month %d", test.dayDiff, test.month)
		}
		if illum < 0 || illum > 1 {
			t.Errorf("Invalid illumination %f on day %d", illum, test.dayDiff)
		}
	}
}

func TestNormalizeAngleEdgeCases(t *testing.T) {
	tests := []struct {
		input    float64
		expected float64
	}{
		{0.0, 0.0},
		{360.0, 0.0},
		{-360.0, 0.0},
		{1080.0, 0.0},
		{359.999, 359.999},
		{0.001, 0.001},
		{-0.001, 359.999},
	}
	
	for _, test := range tests {
		result := normalizeAngle(test.input)
		if math.Abs(result-test.expected) > 0.01 {
			t.Errorf("normalizeAngle(%f) = %f; want %f", test.input, result, test.expected)
		}
	}
}

func TestJulianDayConsistency(t *testing.T) {
	// Test that consecutive days have ~1.0 difference
	time1 := time.Date(2026, 3, 22, 12, 0, 0, 0, time.UTC)
	time2 := time.Date(2026, 3, 23, 12, 0, 0, 0, time.UTC)
	
	jd1 := julianDay(time1)
	jd2 := julianDay(time2)
	
	diff := jd2 - jd1
	if math.Abs(diff-1.0) > 0.01 {
		t.Errorf("Julian day difference between consecutive days should be ~1.0, got %f", diff)
	}
}

func TestJulianDayLeapYear(t *testing.T) {
	// Test leap year handling
	time1 := time.Date(2024, 2, 28, 12, 0, 0, 0, time.UTC)
	time2 := time.Date(2024, 2, 29, 12, 0, 0, 0, time.UTC)
	time3 := time.Date(2024, 3, 1, 12, 0, 0, 0, time.UTC)
	
	jd1 := julianDay(time1)
	jd2 := julianDay(time2)
	jd3 := julianDay(time3)
	
	// Each should differ by ~1.0
	if math.Abs(jd2-jd1-1.0) > 0.01 {
		t.Errorf("Day difference Feb 28-29 should be ~1.0, got %f", jd2-jd1)
	}
	if math.Abs(jd3-jd2-1.0) > 0.01 {
		t.Errorf("Day difference Feb 29-Mar 1 should be ~1.0, got %f", jd3-jd2)
	}
}