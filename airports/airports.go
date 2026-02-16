package airports

import "itinerary-prettifier/types"

// Repository provides airport data access
type Repository interface {
	FindByCode(code string) (*types.Airport, bool)
	GetAll() map[string]types.Airport
}

// Service provides airport-related business logic
type Service interface {
	GetAirportName(code string) string
	GetCityName(code string) string
}

type AirportService struct {
	repo Repository
}

func NewAirportService(repo Repository) *AirportService {
	return &AirportService{repo: repo}
}

func (s *AirportService) GetAirportName(code string) string {
	airport, exists := s.repo.FindByCode(code)
	if !exists {
		return code // Return original code if not found
	}
	return airport.Name
}

func (s *AirportService) GetCityName(code string) string {
	airport, exists := s.repo.FindByCode(code)
	if !exists {
		return code // Return original code if not found
	}
	return airport.Municipality
}
