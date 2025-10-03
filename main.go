package main

// Importing necessary packages from the Go standard library
import (
	"encoding/csv" // for reading CSV files
	"errors"       // for creating custom error messages
	"fmt"          // for printing messages to the console
	"io/ioutil"    // for reading and writing files
	"os"           // for accessing command-line arguments and file operations
	"regexp"       // for pattern matching using regular expressions
	"strings"      // for manipulating strings
	"time"         // for parsing and formatting dates and times
)

// Define a struct to hold airport information
type Airport struct {
	Name         string // Full name of the airport
	Municipality string // City where the airport is located
	IATA         string // 3-letter IATA code
	ICAO         string // 4-letter ICAO code
}

// The main function is the entry point of every Go program
func main() {
	// If the user runs the program with "-h", show usage instructions
	if len(os.Args) == 2 && os.Args[1] == "-h" {
		fmt.Println("itinerary usage:\ngo run . ./input.txt ./output.txt ./airport-lookup.csv")
		return
	}

	// If the user doesn't provide exactly 3 arguments, show usage instructions
	if len(os.Args) != 4 {
		fmt.Println("itinerary usage:\ngo run . ./input.txt ./output.txt ./airport-lookup.csv")
		return
	}

	// Extract file paths from command-line arguments
	inputPath := os.Args[1]
	outputPath := os.Args[2]
	lookupPath := os.Args[3]

	// Try to read the input itinerary file
	inputData, err := ioutil.ReadFile(inputPath)
	if err != nil {
		fmt.Println("Input not found")
		return
	}

	// Try to open the airport lookup CSV file
	lookupFile, err := os.Open(lookupPath)
	if err != nil {
		fmt.Println("Airport lookup not found")
		return
	}
	defer lookupFile.Close() // Make sure the file gets closed when we're done

	// Parse the CSV file and build a map of airport codes to airport info
	airports, err := parseAirportCSV(lookupFile)
	if err != nil {
		fmt.Println("Airport lookup malformed")
		return
	}

	// Process the input text to make it customer-friendly
	output := prettify(string(inputData), airports)

	// Write the processed output to the output file
	err = ioutil.WriteFile(outputPath, []byte(output), 0644)
	if err != nil {
		fmt.Println("Failed to write output")
	}
}

// This function reads and validates the airport CSV file
func parseAirportCSV(file *os.File) (map[string]Airport, error) {
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil || len(records) < 2 {
		return nil, errors.New("CSV malformed")
	}

	// Read the header row to determine column positions
	headers := records[0]
	colIndex := map[string]int{}
	for i, h := range headers {
		colIndex[h] = i
	}

	// Check that all required columns are present
	required := []string{"name", "municipality", "iata_code", "icao_code"}
	for _, r := range required {
		if _, ok := colIndex[r]; !ok {
			return nil, errors.New("missing column")
		}
	}

	// Create a map to store airport data by code
	airportMap := map[string]Airport{}
	for _, row := range records[1:] {
		// Check for blank fields
		for _, idx := range colIndex {
			if idx >= len(row) || strings.TrimSpace(row[idx]) == "" {
				return nil, errors.New("blank field")
			}
		}

		// Create an Airport struct from the row
		airport := Airport{
			Name:         row[colIndex["name"]],
			Municipality: row[colIndex["municipality"]],
			IATA:         row[colIndex["iata_code"]],
			ICAO:         row[colIndex["icao_code"]],
		}

		// Store both airport name and city name lookups
		if airport.IATA != "" {
			airportMap["#"+airport.IATA] = airport
			airportMap["*#"+airport.IATA] = airport
		}
		if airport.ICAO != "" {
			airportMap["##"+airport.ICAO] = airport
			airportMap["*##"+airport.ICAO] = airport
		}
	}
	return airportMap, nil
}

// This function applies all the transformations to the input text
func prettify(text string, airports map[string]Airport) string {
	// Normalize vertical whitespace characters
	text = strings.ReplaceAll(text, "\v", "\n")
	text = strings.ReplaceAll(text, "\f", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")

	// Remove excessive blank lines
	text = collapseBlankLines(text)

	// Replace airport codes with names or cities
	text = replaceAirportCodes(text, airports)

	// Replace date and time formats
	text = replaceDates(text)

	return text
}

// This function ensures no more than one blank line in a row
func collapseBlankLines(text string) string {
	lines := strings.Split(text, "\n")
	var result []string
	blank := false
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			if !blank {
				result = append(result, "")
				blank = true
			}
		} else {
			result = append(result, line)
			blank = false
		}
	}
	return strings.Join(result, "\n")
}

// This function replaces airport codes with readable names or cities
func replaceAirportCodes(text string, airports map[string]Airport) string {
	// Match patterns like #LAX, ##EGLL, *#LHR, *##KJFK
	re := regexp.MustCompile(`\*?#\w{3}|\*?##\w{4}`)
	return re.ReplaceAllStringFunc(text, func(code string) string {
		if airport, ok := airports[code]; ok {
			if strings.HasPrefix(code, "*") {
				return airport.Municipality // Replace with city name
			}
			return airport.Name // Replace with airport name
		}
		return code // Leave unchanged if not found
	})
}

// This function replaces date and time formats with customer-friendly versions
func replaceDates(text string) string {
	// Replace D(...) with formatted date
	dateRe := regexp.MustCompile(`D\(([^)]+)\)`)
	text = dateRe.ReplaceAllStringFunc(text, func(match string) string {
		return formatDate(match, "D")
	})

	// Replace T12(...) with 12-hour time
	t12Re := regexp.MustCompile(`T12\(([^)]+)\)`)
	text = t12Re.ReplaceAllStringFunc(text, func(match string) string {
		return formatDate(match, "T12")
	})

	// Replace T24(...) with 24-hour time
	t24Re := regexp.MustCompile(`T24\(([^)]+)\)`)
	text = t24Re.ReplaceAllStringFunc(text, func(match string) string {
		return formatDate(match, "T24")
	})

	return text
}

// This function parses and formats ISO 8601 date/time strings
func formatDate(match, prefix string) string {
	// Extract the date/time string from inside the parentheses
	re := regexp.MustCompile(prefix + `\(([^)]+)\)`)
	parts := re.FindStringSubmatch(match)
	if len(parts) < 2 {
		return match // If malformed, return original
	}
	raw := parts[1]

	// Normalize Unicode minus sign to ASCII dash
	raw = strings.ReplaceAll(raw, "−", "-")

	// Try to parse the date/time string
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return match // If parsing fails, return original
	}

	// Format the timezone offset
	offset := t.Format("-07:00")
	if strings.HasSuffix(raw, "Z") {
		offset = "(+00:00)" // Zulu time
	} else {
		offset = "(" + offset + ")"
	}

	// Format based on the prefix
	switch prefix {
	case "D":
		return t.Format("02 Jan 2006") // e.g. 05 Apr 2007
	case "T12":
		return t.Format("03:04PM ") + offset // e.g. 12:30PM (-02:00)
	case "T24":
		return t.Format("15:04 ") + offset // e.g. 12:30 (-02:00)
	default:
		return match
	}
}
