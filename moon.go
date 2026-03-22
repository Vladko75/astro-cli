package main

import (
	"math"
	"time"
)

func normalizeAngle(x float64) float64 {
	x = math.Mod(x, 360)
	if x < 0 {
		x += 360
	}
	return x
}

func julianDay(t time.Time) float64 {
	utc := t.UTC()
	y := float64(utc.Year())
	m := float64(utc.Month())
	d := float64(utc.Day()) + float64(utc.Hour())/24.0 + float64(utc.Minute())/1440.0 + float64(utc.Second())/86400.0 + float64(utc.Nanosecond())/8.64e13
	if m <= 2 {
		y -= 1
		m += 12
	}
	A := math.Floor(y / 100)
	B := 2 - A + math.Floor(A/4)
	JD := math.Floor(365.25*(y+4716)) + math.Floor(30.6001*(m+1)) + d + B - 1524.5
	return JD
}

func moonPhaseAndAltitude(t time.Time, latDeg, lonDeg float64) (string, float64, float64) {
	JD := julianDay(t)

	// Approximate moon phase
	T := (JD - 2451545.0) / 36525.0

	// Mean longitude of the moon
	L := normalizeAngle(218.3164477 + 481267.88123421*T - 0.0015786*T*T + T*T*T/538841.0 - T*T*T*T/65194000.0)

	// Mean longitude of the sun
	M := normalizeAngle(280.46646 + 36000.76983*T + 0.0003032*T*T)

	// Moon's argument of latitude
	F := normalizeAngle(93.2720993 + 483202.0175273*T - 0.0036539*T*T - T*T*T/3526000.0 + T*T*T*T/863310000.0)

	// Mean elongation of the moon
	D := normalizeAngle(297.8502042 + 445267.1115168*T - 0.0016300*T*T + T*T*T/545868.0 - T*T*T*T/113065000.0)

	// Phase angle
	phaseAngle := L - M

	// Illumination
	illum := (1 + math.Cos(phaseAngle*math.Pi/180.0)) / 2

	// Phase name
	phase := phaseAngle / 360.0
	phaseName := "New Moon"
	if phase < 0.125 {
		phaseName = "New Moon"
	} else if phase < 0.375 {
		phaseName = "Waxing Crescent"
	} else if phase < 0.625 {
		phaseName = "First Quarter"
	} else if phase < 0.875 {
		phaseName = "Waxing Gibbous"
	} else {
		phaseName = "Full Moon"
	}
	if phase > 0.5 {
		if phaseName == "Waxing Crescent" {
			phaseName = "Waning Crescent"
		} else if phaseName == "Waxing Gibbous" {
			phaseName = "Waning Gibbous"
		} else if phaseName == "First Quarter" {
			phaseName = "Last Quarter"
		}
	}

	// Approximate moon position
	ra := L + 6.289*math.Sin(F*math.Pi/180.0) + 1.274*math.Sin((2*D-L)*math.Pi/180.0) + 0.658*math.Sin(2*D*math.Pi/180.0) + 0.214*math.Sin(2*F*math.Pi/180.0) - 0.110*math.Sin(D*math.Pi/180.0)
	dec := 5.128*math.Sin(F*math.Pi/180.0) + 0.280*math.Sin((L-F)*math.Pi/180.0) + 0.233*math.Sin((2*D)*math.Pi/180.0) + 0.214*math.Sin(2*F*math.Pi/180.0) - 0.110*math.Sin(D*math.Pi/180.0)

	// Altitude calculation
	lst := localSiderealTime(JD, lonDeg)
	hourAngle := normalizeAngle(lst - ra)
	HA := hourAngle * math.Pi / 180.0
	latRad := latDeg * math.Pi / 180.0
	decRad := dec * math.Pi / 180.0
	alt := math.Asin(math.Sin(latRad)*math.Sin(decRad) + math.Cos(latRad)*math.Cos(decRad)*math.Cos(HA))
	altDeg := alt * 180.0 / math.Pi

	return phaseName, illum, altDeg
}
