package csv

import (
	"encoding/csv"
	"fmt"
	"os"
)

type (
	File struct {
		Headers     []Column
		Rows        []Row
		RecordCount int
		Delimiter   rune
	}
)

func New(headers []Column, rowCount int, delimiter rune) *File {
	return &File{
		Headers:   headers,
		Rows:      make([]Row, rowCount),
		Delimiter: delimiter,
	}
}

func (f *File) AddRow() {
	if f.RecordCount >= len(f.Rows) {
		blank := make([]Row, len(f.Rows))
		newRows := append(f.Rows, blank...)
		f.Rows = newRows
	}
	f.Rows[f.RecordCount] = NewRow(f.Headers)
	f.RecordCount += 1
}

func (ef *File) AddRowAtIndex(i int) error {
	return fmt.Errorf("not implemented")
}

func (ef *File) Corrupt(chance float64) error {
	return fmt.Errorf("not implemented")
}

func (f *File) headersAsSlice() []string {
	output := make([]string, len(f.Headers))
	for i, header := range f.Headers {
		output[i] = header.Label
	}
	return output
}

func (f *File) Write(filename string) error {
	osfile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer osfile.Close()

	writer := csv.NewWriter(osfile)
	writer.Comma = f.Delimiter
	err = writer.Write(f.headersAsSlice())
	if err != nil {
		return err
	}

	for _, row := range f.Rows {
		output := row.ToStringSlice(f.Headers)
		err = writer.Write(output)
		if err != nil {
			return err
		}
	}
	writer.Flush()
	return nil
}
