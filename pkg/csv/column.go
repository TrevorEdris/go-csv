package csv

import (
	"fmt"

	"github.com/TrevorEdris/go-csv/pkg/faker"
	"github.com/TrevorEdris/go-csv/pkg/schema"
)

type (
	Headers []Column

	Column struct {
		// These are pointers because not every column will define a value for them
		NumericConstraint   *schema.NumericConstraint   `yaml:"numericConstraint"`
		StringConstraint    *schema.StringConstraint    `yaml:"stringConstraint"`
		TimestampConstraint *schema.TimestampConstraint `yaml:"timestampConstraint"`
		GeneralConstraint   *schema.GeneralConstraint   `yaml:"generalConstraint"`

		Order *int `yaml:"order"`

		// These are required by every column
		Label  string `yaml:"label"`
		Source string `yaml:"source"`
	}
)

func (h Headers) Validate() error {
	// Validate that all columns defined for the header row specify a valid ordering
	maxOrder := len(h) - 1
	specified := make([]bool, len(h))
	for i, c := range h {
		if c.Order != nil {
			o := *c.Order
			if o < 0 || o > maxOrder {
				return fmt.Errorf("column %s order out of range (%d); must be between 0 and %d", c.Label, *c.Order, maxOrder)
			}
			if specified[*c.Order] {
				return fmt.Errorf("column %s order duplicated (%d); each order value must be unique", c.Label, *c.Order)
			}
			specified[*c.Order] = true
		}

		err := c.Validate()
		if err != nil {
			return fmt.Errorf("column %d - %s: %w", i, c.Label, err)
		}
	}
	return nil
}

func (c *Column) Validate() error {
	if c.Label == "" {
		return fmt.Errorf("label must be non-empty")
	}

	valid := false
	for _, f := range faker.AllSources {
		if string(f) == c.Source {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("column %s unsupported source: %s; require one of %v", c.Label, c.Source, faker.AllSources)
	}

	if c.NumericConstraint != nil {
		if err := c.NumericConstraint.Validate(); err != nil {
			return err
		}
	}

	if c.StringConstraint != nil {
		if err := c.StringConstraint.Validate(); err != nil {
			return err
		}
	}

	if c.TimestampConstraint != nil {
		if err := c.TimestampConstraint.Validate(); err != nil {
			return err
		}
	}

	if c.GeneralConstraint != nil {
		if err := c.GeneralConstraint.Validate(); err != nil {
			return err
		}
	}

	return nil
}
