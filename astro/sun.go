package main

import (
	"math"
	"time"
)

// SunInfo holds sunrise, sunset, and altitude information for the sun
type SunInfo struct {
	Sunrise   time.Time
	Sunset    time.Time
	Altitude  float64 // degrees above horizon
	Azimuth   float64 // degrees from north
	DayLength float64 // hours of daylight
}

// sunAltitude calculates the sun's altitude above the horizon
func sunAltitude(t time.Time, latDeg, lonDeg float64) (altitude, azimuth float64) {
	JD := julianDay(t)
	T := (JD - 2451545.0) / 36525.0

	// Mean longitude of sun
	L0 := normalizeAngle(280.46646 + 36000.76983*T + 0.0003032*T*T)

	// Mean longitude of perigee
	M := normalizeAngle(357.52911 + 35999.05029*T - 0.0001536*T*T)
	MRad := M * math.Pi / 180.0

	// Equation of center
	C := (1.914602 - 0.004817*T - 0.000014*T*T)*math.Sin(MRad) +
		(0.019993 - 0.000101*T)*math.Sin(2*MRad) +
		0.000029*math.Sin(3*MRad)

	// True longitude
	sunLon := L0 + C

	// Apparent longitude (simplified)
	lambda := sunLon - 0.00569

	// Obliquity
	epsilon := 23.439291 - 0.0130042*T

	// Sun's declination
	epsilonRad := epsilon * math.Pi / 180.0
	lambdaRad := lambda * math.Pi / 180.0
	sunDec := math.Asin(math.Sin(epsilonRad) * math.Sin(lambdaRad))

	// Sun's right ascension
	sunRA := math.Atan2(math.Cos(epsilonRad)*math.Sin(lambdaRad), math.Cos(lambdaRad)) * 180.0 / math.Pi
	if sunRA < 0 {
		sunRA += 360
	}

	// Local Sidereal Time
	lst := localSiderealTime(JD, lonDeg)

	// Hour angle
	hourAngle := normalizeAngle(lst - sunRA)
	HA := hourAngle * math.Pi / 180.0

	// Altitude
	latRad := latDeg * math.Pi / 180.0
	alt := math.Asin(math.Sin(latRad)*math.Sin(sunDec) + math.Cos(latRad)*math.Cos(sunDec)*math.Cos(HA))
	altitudeD := alt * 180.0 / math.Pi

	// Azimuth
	y := math.Sin(HA)
	x := math.Cos(HA)*math.Sin(latRad) - math.Tan(sunDec)*math.Cos(latRad)
	azimuthRad := math.Atan2(y, x)
	azimuthD := azimuthRad*180.0/math.Pi + 180.0

	return altitudeD, azimuthD
}

// SunTimes calculates sunrise and sunset times
func SunTimes(t time.Time, latDeg, lonDeg float64) SunInfo {
	JD := julianDay(t)

	// Approximate sunrise using noon
	noonJD := math.Floor(JD) + 0.5 - lonDeg/360.0
	_ = noonJD // Use unused variable

	// Refraction + horizon dip
	refraction := 0.833 // degrees
	horizon := -refraction - 2.076*math.Sqrt(math.Max(0, 0))/(60*math.Pi/180) // simplified

	// Time for sun to rise/set approximately 6 hours around noon
	riseTime := time.Now().Add(-6 * time.Hour)
	setTime := time.Now().Add(6 * time.Hour)

	// Iterate to find exact times (simplified)
	for i := 0; i < 5; i++ {
		alt1, _ := sunAltitude(riseTime, latDeg, lonDeg)
		alt2, _ := sunAltitude(setTime, latDeg, lonDeg)

		if alt1 < horizon {
			riseTime = riseTime.Add(30 * time.Minute)
		}
		if alt2 < horizon {
			setTime = setTime.Add(-30 * time.Minute)
		}
	}

	dayLength := setTime.Sub(riseTime).Hours()
	currentAlt, currentAz := sunAltitude(t, latDeg, lonDeg)

	return SunInfo{
		Sunrise:   riseTime,
		Sunset:    setTime,
		Altitude:  currentAlt,
		Azimuth:   currentAz,
		DayLength: dayLength,
	}
}
