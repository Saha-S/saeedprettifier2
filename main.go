package main

import (
	"fmt"
	"itinerary-prettifier/airports"
	"itinerary-prettifier/cli"
	"itinerary-prettifier/config"
	"itinerary-prettifier/fileio"
	"itinerary-prettifier/formatter"
	"strings"
)

func main() {
	// Initialize dependencies
	cliParser := cli.NewCLIParser()
	configValidator := config.NewConfigValidator()
	fileReader := fileio.NewFileReader()
	fileWriter := fileio.NewFileWriter()
	fileChecker := fileio.NewFileChecker()
	csvParser := airports.NewCSVParser()
	airportLoader := airports.NewAirportLoader(csvParser)
	textFormatter := formatter.NewTextFormatter()

	// Parse command line arguments
	config, err := cliParser.Parse()
	if err != nil {
		fmt.Println("itinerary usage:")
		fmt.Println("go run . ./input.txt ./output.txt ./airport-lookup.csv")
		return
	}

	// Validate configuration
	if err := configValidator.Validate(config); err != nil {
		fmt.Println("itinerary usage:")
		fmt.Println("go run . ./input.txt ./output.txt ./airport-lookup.csv")
		return
	}

	// Check if output file exists (to prevent overwrite on error)
	if fileChecker.Exists(config.OutputPath) {
		// We'll proceed but ensure we don't overwrite on error
	}

	// Read input file
	input, err := fileReader.ReadFile(config.InputPath)
	if err != nil {
		fmt.Println("Input not found")
		return
	}

	// Load airport data
	airportRepo, err := airportLoader.Load(config.LookupPath)
	if err != nil {
		// Check what type of error it is
		if strings.Contains(err.Error(), "airport lookup not found") {
			fmt.Println("Airport lookup not found")
		} else if strings.Contains(err.Error(), "airport lookup malformed") {
			fmt.Println("Airport lookup malformed")
		} else {
			fmt.Println("Airport lookup not found") // default case
		}
		return
	}

	// Create airport service
	airportService := airports.NewAirportService(airportRepo)

	// Process and format text
	output := textFormatter.Prettify(input, airportService)

	// Write output file
	if err := fileWriter.WriteFile(config.OutputPath, output); err != nil {
		fmt.Println("Failed to write output")
		return
	}
}
