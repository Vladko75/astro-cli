package main

import (
	"math"
	"sync"
	"time"
)

type orbitalElements struct {
	name string
	a    float64 // semimajor axis (AU)
	e    float64 // eccentricity
	I    float64 // inclination (deg)
	Omega float64 // longitude of ascending node (deg)
	omega float64 // argument of perihelion (deg)
	L0   float64 // mean longitude at epoch (deg)
	n    float64 // mean motion (deg/day)
}

func solveKepler(Mrad, e float64) float64 {
	E := Mrad
	for i := 0; i < 30; i++ {
		f := E - e*math.Sin(E) - Mrad
		fp := 1 - e*math.Cos(E)
		E -= f / fp
		if math.Abs(f) < 1e-12 {
			break
		}
	}
	return E
}

var planets = []orbitalElements{
	{"Mercury", 0.387098, 0.205630, 7.005, 48.331, 29.124, 252.2508, 4.09233445},
	{"Venus", 0.723332, 0.006773, 3.3946, 76.680, 54.884, 181.9798, 1.60213034},
	{"Earth", 1.00000011, 0.016710, 0.00005, 0.0, 102.93768193, 100.46457166, 0.985608},
	{"Mars", 1.523679, 0.0934, 1.850, 49.578, 286.537, 355.45332, 0.52402068},
	{"Jupiter", 5.20260, 0.0485, 1.303, 100.492, 273.867, 34.40438, 0.08308530},
	{"Saturn", 9.5549, 0.0555, 2.485, 113.715, 339.392, 49.94432, 0.03344414},
	{"Uranus", 19.2184, 0.0464, 0.773, 74.006, 96.998857, 313.23218, 0.01172502},
	{"Neptune", 30.1104, 0.0090, 1.770, 131.784, 276.340, 304.88003, 0.00598103},
}

func planetEquatorialCoords(D float64, elem orbitalElements) (raDeg, decDeg float64) {
	M := elem.n * D + elem.L0 - elem.Omega
	Mrad := M * math.Pi / 180.0
	E := solveKepler(Mrad, elem.e)
	nu := 2 * math.Atan(math.Sqrt((1+elem.e)/(1-elem.e)) * math.Tan(E/2))
	r := elem.a * (1 - elem.e*math.Cos(E))
	// ecliptic coordinates
	x := r * (math.Cos(elem.Omega*math.Pi/180.0 + nu) * math.Cos(elem.omega*math.Pi/180.0) - math.Sin(elem.Omega*math.Pi/180.0 + nu) * math.Sin(elem.omega*math.Pi/180.0) * math.Cos(elem.I*math.Pi/180.0))
	y := r * (math.Sin(elem.Omega*math.Pi/180.0 + nu) * math.Cos(elem.omega*math.Pi/180.0) + math.Cos(elem.Omega*math.Pi/180.0 + nu) * math.Sin(elem.omega*math.Pi/180.0) * math.Cos(elem.I*math.Pi/180.0))
	z := r * math.Sin(elem.Omega*math.Pi/180.0 + nu) * math.Sin(elem.I*math.Pi/180.0)
	// equatorial
	ra := math.Atan2(y, x)
	dec := math.Asin(z / math.Sqrt(x*x + y*y + z*z))
	raDeg = ra * 180.0 / math.Pi
	decDeg = dec * 180.0 / math.Pi
	if raDeg < 0 {
		raDeg += 360
	}
	return
}

type PlanetAbove struct {
	Name     string
	Altitude float64
}

func planetsAboveHorizon(latDeg, lonDeg float64, t time.Time) []PlanetAbove {
	JD := julianDay(t)
	D := JD - 2451545.0

	lst := localSiderealTime(JD, lonDeg)
	results := make(chan PlanetAbove, len(planets))
	var wg sync.WaitGroup

	for _, p := range planets {
		if p.name == "Earth" {
			continue
		}
		wg.Add(1)
		go func(planet orbitalElements) {
			defer wg.Done()
			ra, dec := planetEquatorialCoords(D, planet)
			hourAngle := normalizeAngle(lst - ra)
			HA := hourAngle * math.Pi / 180.0
			latRad := latDeg * math.Pi / 180.0
			decRad := dec * math.Pi / 180.0
			alt := math.Asin(math.Sin(latRad)*math.Sin(decRad) + math.Cos(latRad)*math.Cos(decRad)*math.Cos(HA))
			altDeg := alt * 180.0 / math.Pi
			if altDeg > 0 {
				results <- PlanetAbove{Name: planet.name, Altitude: altDeg}
			}
		}(p)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var above []PlanetAbove
	for p := range results {
		above = append(above, p)
	}
	return above
}

func localSiderealTime(JD, lonDeg float64) float64 {
	t := (JD - 2451545.0) / 36525.0
	GST := normalizeAngle(280.46061837 + 360.98564736629*(JD-2451545.0) + 0.000387933*t*t - t*t*t/38710000.0)
	LST := normalizeAngle(GST + lonDeg)
	return LST
}
