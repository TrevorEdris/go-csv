package csv

import "github.com/TrevorEdris/go-csv/app/faker"

type (
	Row struct {
		// Values is a map of column labels to column values. This must be a
		// map to ensure the CSV output matches values to the correct
		// column.
		Values map[string]string
	}
)

func NewRow(headers []Column) Row {
	record := Row{
		Values: make(map[string]string),
	}
	for _, header := range headers {
		record.generateValues(header)
	}
	return record
}

func (r *Row) generateValues(col Column) string {
	v := faker.FakeValue(faker.Source(col.Source), col.StringConstraint, col.NumericConstraint, col.TimestampConstraint, col.GeneralConstraint)

	r.Values[col.Label] = v

	return v
}

func (r *Row) ToStringSlice(headers []Column) []string {
	output := make([]string, len(headers))
	for i, header := range headers {
		output[i] = r.Values[header.Label]
	}
	return output
}
