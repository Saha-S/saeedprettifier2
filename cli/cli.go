package cli

import (
	"errors"
	"flag"
	"fmt"
	"itinerary-prettifier/types"
	"os"
)

// Parser handles command line argument parsing
type Parser interface {
	Parse() (*types.Config, error)
}

type CLIParser struct{}

func NewCLIParser() *CLIParser {
	return &CLIParser{}
}

func (p *CLIParser) Parse() (*types.Config, error) {
	helpFlag := flag.Bool("h", false, "show usage information")
	flag.Parse()

	if *helpFlag {
		fmt.Println("itinerary usage:")
		fmt.Println("go run . ./input.txt ./output.txt ./airport-lookup.csv")
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) != 3 {
		return nil, ErrInvalidArguments
	}

	return &types.Config{
		InputPath:  args[0],
		OutputPath: args[1],
		LookupPath: args[2],
	}, nil
}

// CLI errors
var (
	ErrInvalidArguments = errors.New("invalid number of arguments")
)
