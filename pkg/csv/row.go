package csv

import (
	"github.com/TrevorEdris/go-csv/pkg/faker"
	"github.com/TrevorEdris/go-csv/pkg/schema"
)

type (
	row struct {
		// Values is a map of column labels to column values. This must be a
		// map to ensure the CSV output matches values to the correct
		// column.
		Values map[string]string
	}
)

func newRow(headers schema.Headers) row {
	record := row{
		Values: make(map[string]string),
	}
	for _, header := range headers {
		record.generateValues(header)
	}
	return record
}

func (r *row) generateValues(col schema.Column) string {
	v := faker.FakeValue(faker.Source(col.Source), col.StringConstraint, col.NumericConstraint, col.TimestampConstraint, col.GeneralConstraint)

	r.Values[col.Label] = v

	return v
}

func (r *row) ToStringSlice(headers schema.Headers) []string {
	output := make([]string, len(headers))
	for i, header := range headers {
		if header.Order != nil {
			output[*header.Order] = r.Values[header.Label]
		} else {
			output[i] = r.Values[header.Label]
		}
	}
	return output
}
