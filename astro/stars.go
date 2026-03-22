package main

import (
	"math"
	"sort"
	"time"
)

// Star represents a bright star with name, coordinates, and magnitude
type Star struct {
	Name      string
	RA        float64 // right ascension in degrees
	Dec       float64 // declination in degrees
	Magnitude float64 // apparent magnitude (lower = brighter)
	Altitude  float64 // altitude above horizon at observation time
}

// brightStars is a catalog of prominent bright stars (Hipparcos catalog subset)
var brightStars = []struct {
	name      string
	ra        float64
	dec       float64
	magnitude float64
}{
	{"Sirius", 101.2872, -16.7161, -1.46},
	{"Canopus", 95.9876, -52.6957, -0.74},
	{"Rigil Kentaurus", 219.9012, -60.8370, -0.27},
	{"Arcturus", 213.9154, 19.1825, -0.05},
	{"Vega", 279.2343, 38.7837, 0.03},
	{"Capella", 79.1721, 45.9980, 0.08},
	{"Rigel", 78.6344, -8.2017, 0.18},
	{"Procyon", 114.8250, 5.2250, 0.40},
	{"Achernar", 24.4283, -57.2696, 0.46},
	{"Betelgeuse", 88.7929, 7.4070, 0.50},
	{"Altair", 297.6958, 8.8683, 0.76},
	{"Acrux", 186.6496, -63.0996, 0.77},
	{"Aldebaran", 68.9803, 16.5072, 0.87},
	{"Spica", 201.2983, -11.1613, 0.98},
	{"Antares", 247.3518, -26.4322, 1.06},
	{"Pollux", 131.2871, 28.0263, 1.16},
	{"Fomalhaut", 344.4125, -29.6224, 1.17},
	{"Denebola", 168.9862, 14.5720, 2.14},
	{"Regulus", 152.0933, 11.9673, 1.35},
	{"Castor", 113.6497, 31.8886, 1.58},
	{"Shaula", 263.4020, -37.1043, 1.62},
	{"Bellatrix", 81.2828, 6.3497, 1.64},
	{"Mizar", 200.9819, 54.9247, 2.27},
	{"Alnilam", 84.0530, -1.2019, 1.69},
	{"Alnitak", 85.2656, -2.0748, 1.93},
	{"Alioth", 166.4557, 55.9598, 1.76},
	{"Dubhe", 164.0757, 61.7507, 1.79},
	{"Deneb", 309.2281, 45.2803, 1.25},
	{"Gienah", 222.6497, -17.5439, 2.58},
	{"Eltanin", 261.5168, 51.4891, 2.24},
}

// StarsAboveHorizon returns bright stars visible above horizon at given location and time
func StarsAboveHorizon(latDeg, lonDeg float64, t time.Time) []Star {
	JD := julianDay(t)
	lst := localSiderealTime(JD, lonDeg)

	var visibleStars []Star

	for _, starData := range brightStars {
		hourAngle := normalizeAngle(lst - starData.ra)
		HA := hourAngle * math.Pi / 180.0
		latRad := latDeg * math.Pi / 180.0
		decRad := starData.dec * math.Pi / 180.0

		alt := math.Asin(math.Sin(latRad)*math.Sin(decRad) + math.Cos(latRad)*math.Cos(decRad)*math.Cos(HA))
		altDeg := alt * 180.0 / math.Pi

		// Only include stars above 0 degrees altitude
		if altDeg > 0 {
			visibleStars = append(visibleStars, Star{
				Name:      starData.name,
				RA:        starData.ra,
				Dec:       starData.dec,
				Magnitude: starData.magnitude,
				Altitude:  altDeg,
			})
		}
	}

	// Sort by brightness (lower magnitude = brighter)
	sort.Slice(visibleStars, func(i, j int) bool {
		return visibleStars[i].Magnitude < visibleStars[j].Magnitude
	})

	return visibleStars
}
