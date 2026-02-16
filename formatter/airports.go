package formatter

import (
	"itinerary-prettifier/airports"
	"regexp"
	"strings"
)

type AirportCodeReplacer struct{}

func NewAirportFormatter() *AirportCodeReplacer {
	return &AirportCodeReplacer{}
}

func (f *AirportCodeReplacer) ReplaceAirportCodes(text string, service airports.Service) string {
	// Less restrictive patterns - allow codes at start/end of line or surrounded by whitespace/punctuation
	patterns := []string{
		`(\*?##[A-Z]{3,4})(?:\s|$|\.|,|;|!|\?|\)|"|')`, // *##ABCD or ##ABCD
		`(\*?#[A-Z]{3})(?:\s|$|\.|,|;|!|\?|\)|"|')`,    // *#ABC or #ABC
	}
	
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			// Extract the code part (remove trailing delimiter)
			code := strings.TrimRight(match, " .,;!?)\"'")
			
			var replacement string
			if strings.HasPrefix(code, "*") {
				replacement = service.GetCityName(code)
			} else {
				replacement = service.GetAirportName(code)
			}
			
			// Only replace if we found the airport
			if replacement != code {
				// Restore the trailing delimiter
				if len(match) > len(code) {
					return replacement + match[len(code):]
				}
				return replacement
			}
			
			return match
		})
	}
	
	return text
}