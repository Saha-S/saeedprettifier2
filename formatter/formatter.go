package formatter

import "itinerary-prettifier/airports"

// Formatter orchestrates all text formatting operations
type Formatter interface {
	Prettify(text string, airportService airports.Service) string
}

// AirportFormatter replaces airport codes with full names using airports.Service
type AirportFormatter interface {
	ReplaceAirportCodes(text string, airportService airports.Service) string
}

type TextFormatter struct {
	whitespaceFormatter WhitespaceFormatter
	airportFormatter    AirportFormatter
	dateFormatter       DateFormatter
}

func NewTextFormatter() *TextFormatter {
	return &TextFormatter{
		whitespaceFormatter: NewWhitespaceFormatter(),
		airportFormatter:    NewAirportFormatter(),
		dateFormatter:       NewDateFormatter(),
	}
}

func (f *TextFormatter) Prettify(text string, airportService airports.Service) string {
	// Apply transformations in correct order
	text = f.whitespaceFormatter.ConvertControlChars(text)
	text = f.whitespaceFormatter.CollapseBlankLines(text)
	text = f.airportFormatter.ReplaceAirportCodes(text, airportService)
	text = f.dateFormatter.ReplaceTimesThenDates(text)
	text = f.whitespaceFormatter.TrimExcessiveWhitespace(text)
	return text
}
