package config

import (
	"errors"
	"itinerary-prettifier/types"
)

// Validator validates configuration
type Validator interface {
	Validate(config *types.Config) error
}

type ConfigValidator struct{}

func NewConfigValidator() *ConfigValidator {
	return &ConfigValidator{}
}

func (v *ConfigValidator) Validate(config *types.Config) error {
	if config.InputPath == "" {
		return ErrInputPathRequired
	}
	if config.OutputPath == "" {
		return ErrOutputPathRequired
	}
	if config.LookupPath == "" {
		return ErrLookupPathRequired
	}
	return nil
}

// Configuration errors
var (
	ErrInputPathRequired  = errors.New("input path is required")
	ErrOutputPathRequired = errors.New("output path is required")
	ErrLookupPathRequired = errors.New("lookup path is required")
)
