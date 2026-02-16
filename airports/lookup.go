package airports

import (
	"fmt"
	"os"
	"itinerary-prettifier/types"
)
type AirportRepository struct {
	airports map[string]types.Airport
}

func NewAirportRepository(airports map[string]types.Airport) *AirportRepository {
	return &AirportRepository{airports: airports}
}

func (r *AirportRepository) FindByCode(code string) (*types.Airport, bool) {
	airport, exists := r.airports[code]
	return &airport, exists
}

func (r *AirportRepository) GetAll() map[string]types.Airport {
	return r.airports
}

// Loader handles loading airport data from file
type Loader interface {
	Load(lookupPath string) (Repository, error)
}

type AirportLoader struct {
	parser Parser
}

func NewAirportLoader(parser Parser) *AirportLoader {
	return &AirportLoader{parser: parser}
}

func (l *AirportLoader) Load(lookupPath string) (Repository, error) {
	file, err := os.Open(lookupPath)
	if err != nil {
		return nil, fmt.Errorf("airport lookup not found: %w", err)
	}
	defer file.Close()

	airports, err := l.parser.Parse(file)
	if err != nil {
		return nil, fmt.Errorf("airport lookup malformed: %w", err)
	}

	return NewAirportRepository(airports), nil
}
