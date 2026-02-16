package airports

import (
	"encoding/csv"
	"fmt"
	"io"
	"itinerary-prettifier/types"
	"os"
	"strings"
)

// Parser handles CSV parsing for airport data
type Parser interface {
	Parse(file *os.File) (map[string]types.Airport, error)
}

type CSVParser struct {
	requiredColumns []string
}

func NewCSVParser() *CSVParser {
	return &CSVParser{
		requiredColumns: []string{"name", "iso_country", "municipality", "icao_code", "iata_code", "coordinates"},
	}
}

func (p *CSVParser) Parse(file *os.File) (map[string]types.Airport, error) {
	reader := csv.NewReader(file)

	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	columnMap, err := p.validateHeaders(headers)
	if err != nil {
		return nil, err
	}

	airportMap := make(map[string]types.Airport)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading record: %w", err)
		}

		airport, err := p.parseRecord(record, columnMap)
		if err != nil {
			return nil, err
		}

		p.addAirportToMap(airportMap, airport)
	}

	return airportMap, nil
}

func (p *CSVParser) validateHeaders(headers []string) (map[string]int, error) {
	columnMap := make(map[string]int)
	for i, header := range headers {
		cleanHeader := strings.TrimSpace(strings.ToLower(header))
		columnMap[cleanHeader] = i
	}

	for _, required := range p.requiredColumns {
		if _, exists := columnMap[required]; !exists {
			return nil, fmt.Errorf("missing required column: %s", required)
		}
	}

	return columnMap, nil
}

func (p *CSVParser) parseRecord(record []string, columnMap map[string]int) (*types.Airport, error) {
	// Check if record has enough columns
	if len(record) < len(p.requiredColumns) {
		return nil, fmt.Errorf("record has insufficient columns")
	}

	// Validate no blank fields in required columns
	for _, required := range p.requiredColumns {
		idx := columnMap[required]
		if idx >= len(record) || strings.TrimSpace(record[idx]) == "" {
			return nil, fmt.Errorf("blank field in required column: %s", required)
		}
	}

	return &types.Airport{
		Name:         strings.TrimSpace(record[columnMap["name"]]),
		ISOCountry:   strings.TrimSpace(record[columnMap["iso_country"]]),
		Municipality: strings.TrimSpace(record[columnMap["municipality"]]),
		ICAO:         strings.TrimSpace(record[columnMap["icao_code"]]),
		IATA:         strings.TrimSpace(record[columnMap["iata_code"]]),
		Coordinates:  strings.TrimSpace(record[columnMap["coordinates"]]),
	}, nil
}

func (p *CSVParser) addAirportToMap(airportMap map[string]types.Airport, airport *types.Airport) {
	if airport.IATA != "" {
		airportMap["#"+airport.IATA] = *airport
		airportMap["*#"+airport.IATA] = *airport
	}
	if airport.ICAO != "" {
		airportMap["##"+airport.ICAO] = *airport
		airportMap["*##"+airport.ICAO] = *airport
	}
}
