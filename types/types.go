package types

// Airport holds CSV airport data
type Airport struct {
	Name         string
	ISOCountry   string
	Municipality string
	ICAO         string
	IATA         string
	Coordinates  string
}

// Config holds application configuration
type Config struct {
	InputPath  string
	OutputPath string
	LookupPath string
}

// ProcessingResult holds the result of processing
type ProcessingResult struct {
	Output string
	Error  error
}
