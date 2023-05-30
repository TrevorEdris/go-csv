package config

import (
	"fmt"
	"os"

	"github.com/TrevorEdris/go-csv/app/csv"
	"gopkg.in/yaml.v3"
)

type (
	Runtime struct {
		LogLevel string
		Input    string
		Output   string
	}

	Config struct {
		Runtime    *Runtime
		Columns    csv.Headers `yaml:"columns"`
		Metadata   Metadata    `yaml:"metadata"`
		Corruption Corruption  `yaml:"corruption"`
	}

	Metadata struct {
		Delimiter     string `yaml:"delimiter"`
		DelimiterRune rune
		RowCount      int `yaml:"rowCount"`
	}

	Corruption struct {
		Method  string  `yaml:"method"`
		Chance  float64 `yaml:"chance"`
		Min     int     `yaml:"min"`
		Max     int     `yaml:"max"`
		Enabled bool    `yaml:"enabled"`
	}
)

func New(runtime *Runtime) (*Config, error) {
	c := Config{}
	data, err := os.ReadFile(runtime.Input)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", runtime.Input, err)
	}

	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml from file %s: %w", runtime.Input, err)
	}

	err = c.Validate()
	if err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	c.Runtime = runtime

	return &c, nil
}

func (c *Config) Validate() error {
	err := c.Columns.Validate()
	if err != nil {
		return err
	}

	c.Columns = c.Columns.Sort()

	err = c.Metadata.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (m *Metadata) Validate() error {
	if m.RowCount <= 0 {
		return fmt.Errorf("metadata rowcount invalid (%d): must be > 0", m.RowCount)
	}

	if m.Delimiter == "" {
		m.DelimiterRune = ','
	} else if len(m.Delimiter) == 1 {
		m.DelimiterRune = rune(m.Delimiter[0])
	} else {
		return fmt.Errorf("metadata delimiter (%s) must be a string of length 1", m.Delimiter)
	}

	return nil
}
