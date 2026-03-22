# Astro CLI

A Go command-line tool for astronomical information: current time, moon phase/altitude, and visible planets above the horizon.

## Features

- Display current local time for a location
- Moon phase, illumination percentage, and altitude above horizon
- List of planets visible above the horizon with their altitudes
- Support for top 100 cities by population or custom coordinates
- Custom time input (default: now)
- Selective output (time, moon, planets, or all)

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
- `-show string`: What to show: all, time, moon, planets (comma-separated). Default: all.

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

- Custom location (New York), specific time, moon only:
  ```bash
  ./astro --lat=40.7128 --lon=-74.0060 --time=2026-03-22T12:00:00Z --show=moon
  ```

- Time and moon for Marmaris:
  ```bash
  ./astro --show=time,moon
  ```

## Output Example

```
Current time in Marmaris, Turkey: 2026-03-22 13:09:48 +03
Moon phase in Marmaris, Turkey: Waxing Crescent (illumination 16.3%), altitude above horizon: 51.45°
Планеты над горизонтом в локации (36.8529, 28.2744):
 - Venus: 46.11°
 - Mars: 42.24°
 - Saturn: 54.27°
 - Uranus: 35.87°
 - Neptune: 53.05°
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

Current coverage: ~71% (covers astronomical calculations and output functions).</content>
<parameter name="filePath">/home/haver/vscode-t1/README.md