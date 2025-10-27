package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"
)

// Airport holds CSV airport data
type Airport struct {
	Name         string
	Municipality string
	IATA         string
	ICAO         string
}

func main() {
	// help
	if len(os.Args) == 2 && os.Args[1] == "-h" {
		fmt.Println("itinerary usage:\ngo run . ./input.txt ./output.txt ./airport-lookup.csv")
		return
	}

	// expect exactly 3 args
	if len(os.Args) != 4 {
		fmt.Println("itinerary usage:\ngo run . ./input.txt ./output.txt ./airport-lookup.csv")
		return
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]
	lookupPath := os.Args[3]

	// read input
	inBytes, err := ioutil.ReadFile(inputPath)
	if err != nil {
		fmt.Println("Input not found")
		return
	}
	input := string(inBytes)

	// open lookup CSV
	lf, err := os.Open(lookupPath)
	if err != nil {
		fmt.Println("Airport lookup not found")
		return
	}
	defer lf.Close()

	airports, err := parseAirportCSV(lf)
	if err != nil {
		fmt.Println("Airport lookup malformed")
		return
	}

	// process
	out := prettify(input, airports)

	// write output
	if err := ioutil.WriteFile(outputPath, []byte(out), 0644); err != nil {
		fmt.Println("Failed to write output")
		return
	}

	fmt.Println("✨ Itinerary prettified successfully! ✨")
}

// parseAirportCSV reads CSV and supports shuffled columns; validates blanks
func parseAirportCSV(f *os.File) (map[string]Airport, error) {
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil || len(records) < 2 {
		return nil, errors.New("CSV malformed")
	}

	// map headers -> index (lowercased trimmed)
	headers := records[0]
	col := map[string]int{}
	for i, h := range headers {
		col[strings.TrimSpace(strings.ToLower(h))] = i
	}

	required := []string{"name", "municipality", "iata_code", "icao_code"}
	for _, req := range required {
		if _, ok := col[req]; !ok {
			return nil, errors.New("missing column")
		}
	}

	airportMap := map[string]Airport{}
	for _, row := range records[1:] {
		// validate required fields present and non-blank
		for _, req := range required {
			idx := col[req]
			if idx >= len(row) || strings.TrimSpace(row[idx]) == "" {
				return nil, errors.New("blank field")
			}
		}
		a := Airport{
			Name:         strings.TrimSpace(row[col["name"]]),
			Municipality: strings.TrimSpace(row[col["municipality"]]),
			IATA:         strings.TrimSpace(row[col["iata_code"]]),
			ICAO:         strings.TrimSpace(row[col["icao_code"]]),
		}
		if a.IATA != "" {
			airportMap["#"+a.IATA] = a
			airportMap["*#"+a.IATA] = a
		}
		if a.ICAO != "" {
			airportMap["##"+a.ICAO] = a
			airportMap["*##"+a.ICAO] = a
		}
	}
	return airportMap, nil
}

// prettify applies transformations in a safe order
func prettify(text string, airports map[string]Airport) string {
	// 1) Convert literal backslash sequences (e.g. backslash + 'v') to real newlines
	text = strings.ReplaceAll(text, `\v`, "\n")
	text = strings.ReplaceAll(text, `\f`, "\n")
	text = strings.ReplaceAll(text, `\r`, "\n")

	// 2) Convert any actual control characters if present
	text = strings.ReplaceAll(text, "\v", "\n")
	text = strings.ReplaceAll(text, "\f", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")

	// 3) Collapse excessive blank lines (no more than one empty line)
	text = collapseBlankLines(text)

	// 4) Replace airport codes (#LAX, ##EGLL, and * variants)
	text = replaceAirportCodes(text, airports)

	// 5) Replace date/time tokens: process T12 and T24 first, then D
	text = replaceTimesThenDates(text)

	return text
}

func collapseBlankLines(text string) string {
	lines := strings.Split(text, "\n")
	out := make([]string, 0, len(lines))
	blank := false
	for _, ln := range lines {
		if strings.TrimSpace(ln) == "" {
			if !blank {
				out = append(out, "")
				blank = true
			}
		} else {
			out = append(out, ln)
			blank = false
		}
	}
	return strings.Join(out, "\n")
}

func replaceAirportCodes(text string, airports map[string]Airport) string {
	re := regexp.MustCompile(`\*?#\w{3}|\*?##\w{4}`)
	return re.ReplaceAllStringFunc(text, func(tok string) string {
		if a, ok := airports[tok]; ok {
			if strings.HasPrefix(tok, "*") {
				return a.Municipality
			}
			return a.Name
		}
		return tok
	})
}

// replaceTimesThenDates processes T12/T24 first (so they don't get mistaken for D),
// then processes D(...).
func replaceTimesThenDates(text string) string {
	t12Re := regexp.MustCompile(`T12\(\s*([^)]+?)\s*\)`)
	text = t12Re.ReplaceAllStringFunc(text, func(m string) string {
		return formatDateToken(m, "T12")
	})

	t24Re := regexp.MustCompile(`T24\(\s*([^)]+?)\s*\)`)
	text = t24Re.ReplaceAllStringFunc(text, func(m string) string {
		return formatDateToken(m, "T24")
	})

	// Now D(...)
	dRe := regexp.MustCompile(`D\(\s*([^)]+?)\s*\)`)
	text = dRe.ReplaceAllStringFunc(text, func(m string) string {
		return formatDateToken(m, "D")
	})

	return text
}

// formatDateToken parses the token and returns formatted string (or original if malformed)
func formatDateToken(token, typ string) string {
	re := regexp.MustCompile(typ + `\(\s*([^)]+?)\s*\)`)
	parts := re.FindStringSubmatch(token)
	if len(parts) < 2 {
		return token
	}
	raw := strings.TrimSpace(parts[1])
	raw = strings.ReplaceAll(raw, "−", "-") // normalize unicode minus

	// Try multiple layouts
	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04-07:00",
		"2006-01-02T15:04:05-07:00",
		"2006-01-02", // date-only fallback
	}

	var parsed time.Time
	var err error
	for _, layout := range layouts {
		parsed, err = time.Parse(layout, raw)
		if err == nil {
			break
		}
		// If raw has a T and layout is date-only, try the part before 'T'
		if layout == "2006-01-02" && strings.Contains(raw, "T") {
			dateOnly := strings.SplitN(raw, "T", 2)[0]
			parsed, err = time.Parse(layout, dateOnly)
			if err == nil {
				break
			}
		}
	}
	if err != nil {
		return token // malformed -> leave unchanged
	}

	// Format offset: Z -> (+00:00)
	offset := "(" + parsed.Format("-07:00") + ")"
	if strings.HasSuffix(strings.ToUpper(raw), "Z") {
		offset = "(+00:00)"
	}

	switch typ {
	case "D":
		return parsed.Format("02 Jan 2006")
	case "T12":
		return parsed.Format("03:04PM ") + offset
	case "T24":
		return parsed.Format("15:04 ") + offset
	default:
		return token
	}
}
