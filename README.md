# Itinerary Prettifier

Itinerary Prettifier is a Go command-line utility that cleans and enriches raw trip itineraries. It normalizes whitespace, expands airport codes into readable names, and renders timestamp tokens into human-friendly strings by consulting an airport lookup CSV.

## Features

- Converts control characters and collapses excessive blank lines.
- Replaces IATA (`#ABC`) and ICAO (`##ABCD`) tokens with airport names from a lookup file.
- Turns city-prefixed tokens (`*#ABC`, `*##ABCD`) into municipality names for quick reference.
- Formats time tokens (`T12(...)`, `T24(...)`) with the correct offset and date tokens (`D(...)`) into `DD Mon YYYY`.
- Leaves unknown airport codes untouched, making it safe to run on partially curated inputs.

## Requirements

- Go 1.24.5 or later.
- UTF-8 text input file containing the itinerary content.
- CSV lookup file with airport metadata (see below for required columns).

## Quick Start

```bash
# Fetch dependencies (none beyond the Go standard library)
go mod tidy

# Format your input/lookup files and run the prettifier
go run . ./input.txt ./output.txt ./airport-lookup.csv
```

`input.txt` is the raw itinerary to prettify, `output.txt` is where the cleaned text is written, and `airport-lookup.csv` is the lookup table. The program prints short error messages on stderr when it cannot proceed (e.g., invalid argument count, missing files, malformed CSV).

## Airport Lookup CSV

The loader expects a header row containing at least the following columns:

- `name`
- `iso_country`
- `municipality`
- `icao_code`
- `iata_code`
- `coordinates`

Each row must populate those columns. During parsing the tool adds multiple lookup keys so that both `#IATA` and `##ICAO` tokens can resolve to the same airport.

## Token Reference

| Token | Meaning | Example Input | Output Example |
| ----- | ------- | ------------- | -------------- |
| `#ABC` | IATA airport code → airport name | `Depart #SEA` | `Depart Seattle-Tacoma International Airport` |
| `##ABCD` | ICAO airport code → airport name | `Gate ##KSEA` | `Gate Seattle-Tacoma International Airport` |
| `*#ABC` | City hint for IATA code | `Arrive *#LHR` | `Arrive London` |
| `*##ABCD` | City hint for ICAO code | `Layover *##EGLL` | `Layover London` |
| `T24(ISO timestamp)` | 24-hour clock with offset | `T24(2025-03-05T08:15:00-08:00)` | `08:15 (-08:00)` |
| `T12(ISO timestamp)` | 12-hour clock with offset | `T12(2025-03-05T16:45:00+00:00)` | `04:45PM (+00:00)` |
| `D(ISO date or datetime)` | Calendar date | `D(2025-03-05)` | `05 Mar 2025` |

Tokens remain unchanged when their lookup fails or the timestamp/date cannot be parsed, so your source data stays intact.

## Example

```
Input:
Depart #SEA on D(2025-03-05) at T24(2025-03-05T08:15:00-08:00).
Arrive *#LHR at T12(2025-03-05T16:45:00+00:00).

Output:
Depart Seattle-Tacoma International Airport on 05 Mar 2025 at 08:15 (-08:00).
Arrive London at 04:45PM (+00:00).
```

Run the sample command above to process the provided `input.txt` and inspect `output.txt` for the transformed result.

## Project Structure

```
.
├── airports/     # CSV parsing and airport lookup services
├── cli/          # Command-line argument parsing
├── config/       # Configuration validation
├── fileio/       # File reader/writer helpers
├── formatter/    # Text prettification pipeline
├── types/        # Shared data structures
├── main.go       # Composition root wiring everything together
└── go.mod        # Module definition
```
