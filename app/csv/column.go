package csv

import (
	"fmt"

	"github.com/TrevorEdris/go-csv/app/faker"
)

type (
	Headers []Column

	Column struct {
		// These are pointers because not every column will define a value for them
		NumericConstraint   *faker.NumericConstraint   `yaml:"numericConstraint"`
		StringConstraint    *faker.StringConstraint    `yaml:"stringConstraint"`
		TimestampConstraint *faker.TimestampConstraint `yaml:"timestampConstraint"`
		GeneralConstraint   *faker.GeneralConstraint   `yaml:"generalConstraint"`

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

func (h Headers) Sort() Headers {
	newH := make(Headers, len(h))

	withOrder, withoutOrder := separateColumnsByOrder(h)
	for _, c := range withOrder {
		newH[*c.Order] = c
	}

	withoutOrderIndex := 0
	emptyCol := Column{}
	for newHIndex := 0; newHIndex < len(newH); newHIndex++ {
		if newH[newHIndex] == emptyCol {
			newH[newHIndex] = withoutOrder[withoutOrderIndex]
			withoutOrderIndex++
		}
	}
	return newH
}

func separateColumnsByOrder(h Headers) (withOrder, withoutOrder Headers) {
	for _, col := range h {
		if col.Order != nil {
			withOrder = append(withOrder, col)
		} else {
			withoutOrder = append(withoutOrder, col)
		}
	}
	return withOrder, withoutOrder
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
