# Astro CLI

A Go command-line tool for astronomical information: current time, moon phase/altitude, and visible planets above the horizon.

## Features

- Display current local time for a location with timezone support
- Moon phase, illumination percentage, and altitude above horizon
- Sun altitude, azimuth, and rise/set times
- Bright stars visible above the horizon (30+ star catalog)
- List of planets visible above the horizon with their altitudes and zodiac symbols
- Support for top 100 cities by population, polar regions, or custom coordinates
- Custom time input (default: current time)
- Selective output (time, moon, planets, sun, stars, or any combination)
- Color-coded output with ANSI terminal support (disable with `--no-color`)
- Timezone handling for 120+ IANA timezones

## Installation

Ensure Go 1.21+ is installed.

Clone or download the project, then:

```bash
go mod tidy
go build -o astro .
```

## Usage

```bash
./astro [options]
```

### Options

- `-city string`: City name (e.g., tokyo, london, etc.). If not specified, defaults to Marmaris.
- `-lat float`: Latitude in degrees (overrides city).
- `-lon float`: Longitude in degrees (overrides city).
- `-time string`: Time in RFC3339 format (e.g., 2026-03-22T12:00:00Z). Default: current time.
- `-show string`: What to show: all, time, moon, planets, sun, stars (comma-separated). Default: all.
- `--no-color`: Disable colored output.
- `--help`: Show help information.
- `--version`: Show version information.

### Examples

- Default (Marmaris, all info, current time):
  ```bash
  ./astro
  ```

- Tokyo, only time:
  ```bash
  ./astro --city=tokyo --show=time
  ```

- London, planets at specific time:
  ```bash
  ./astro --city=london --show=planets --time=2026-03-22T20:00:00Z
  ```

- North Pole, moon info:
  ```bash
  ./astro --city="north pole" --show=moon
  ```

- Custom location (New York), sun and stars without color:
  ```bash
  ./astro --lat=40.7128 --lon=-74.0060 --show=sun,stars --no-color
  ```

- Time and moon for Marmaris:
  ```bash
  ./astro --show=time,moon
  ```

- Show help:
  ```bash
  ./astro --help
  ```

## Output Example

```
Current time in Marmaris, Turkey: 2026-03-22 13:09:48 +03

Moon phase in Marmaris, Turkey: Waxing Crescent (illumination 16.3%), altitude above horizon: 51.45°

Sun in Marmaris, Turkey: altitude 45.23°, azimuth 112.5°
Sunrise at 06:12, Sunset at 18:45

Visible planets above the horizon at (36.8529, 28.2744):
 ♀ Venus: 46.11°
 ♂ Mars: 42.24°
 ♄ Saturn: 54.27°
 ♅ Uranus: 35.87°
 ♆ Neptune: 53.05°

Bright stars above the horizon at (36.8529, 28.2744):
 Sirius: 52.34°, magnitude 1.46
 Canopus: 38.12°, magnitude 0.74
 Alpha Centauri: 12.33°, magnitude 0.27
```

## Testing

Run unit tests:

```bash
go test .
```

Run with coverage:

```bash
go test -cover .
```

**Current coverage: 66.0%** (comprehensive test suite covering astronomical calculations, data extraction, output functions, and edge cases across various locations and times).</content>
<parameter name="filePath">/home/haver/vscode-t1/README.md