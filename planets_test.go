package main

import (
	"math"
	"testing"
	"time"
)

func TestSolveKepler(t *testing.T) {
	tests := []struct {
		M, e float64
	}{
		{0.5, 0.1},
		{1.0, 0.0}, // Circular orbit
		{2.0, 0.5},
	}
	for _, test := range tests {
		E := solveKepler(test.M, test.e)
		f := E - test.e*math.Sin(E) - test.M
		if math.Abs(f) > 1e-10 {
			t.Errorf("Kepler equation not satisfied for M=%f, e=%f: f = %e", test.M, test.e, f)
		}
	}
}

func TestPlanetEquatorialCoords(t *testing.T) {
	D := 0.0
	elem := orbitalElements{name: "Earth", a: 1.0, e: 0.0167, I: 0, Omega: 0, omega: 102.937, L0: 100.464, n: 0.9856}
	ra, dec := planetEquatorialCoords(D, elem)
	if ra < 0 || ra >= 360 {
		t.Errorf("RA %f out of range", ra)
	}
	if dec < -90 || dec > 90 {
		t.Errorf("Dec %f out of range", dec)
	}
	// Test another planet
	elem2 := orbitalElements{name: "Mars", a: 1.524, e: 0.0934, I: 1.85, Omega: 49.578, omega: 286.537, L0: 355.453, n: 0.524}
	ra2, dec2 := planetEquatorialCoords(D, elem2)
	if ra2 < 0 || ra2 >= 360 {
		t.Errorf("RA2 %f out of range", ra2)
	}
	if dec2 < -90 || dec2 > 90 {
		t.Errorf("Dec2 %f out of range", dec2)
	}
}

func TestLocalSiderealTime(t *testing.T) {
	JD := 2451545.0
	lon := 0.0
	lst := localSiderealTime(JD, lon)
	if lst < 0 || lst >= 360 {
		t.Errorf("LST %f out of range", lst)
	}
	// For JD 2451545.0 at lon 0, LST should be around 0 or known value
	expectedApprox := 280.46061837 // Approximate GST at J2000
	if math.Abs(lst-expectedApprox) > 1 {
		t.Errorf("LST %f not close to expected ~%f", lst, expectedApprox)
	}
}

func TestPlanetsAboveHorizon(t *testing.T) {
	tests := []struct {
		lat, lon float64
		year, month, day int
		expectSome bool
	}{
		{55.7558, 37.6173, 2026, 3, 22, true}, // Moscow, day
		{36.8529, 28.2744, 2026, 3, 22, true}, // Marmaris, day
		{0, 0, 2026, 3, 22, true}, // Equator
	}
	for _, test := range tests {
		testTime := time.Date(test.year, time.Month(test.month), test.day, 12, 0, 0, 0, time.UTC)
		planets := planetsAboveHorizon(test.lat, test.lon, testTime)
		if test.expectSome && len(planets) == 0 {
			t.Errorf("Expected some planets above horizon for lat=%f, lon=%f", test.lat, test.lon)
		}
		validNames := map[string]bool{"Mercury": true, "Venus": true, "Mars": true, "Jupiter": true, "Saturn": true, "Uranus": true, "Neptune": true}
		for _, p := range planets {
			if !validNames[p.Name] {
				t.Errorf("Invalid planet name: %s", p.Name)
			}
			if p.Altitude < 0 {
				t.Errorf("Altitude %f should be >=0 for visible planets", p.Altitude)
			}
		}
	}
}