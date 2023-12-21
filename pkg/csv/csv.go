package csv

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/TrevorEdris/banner"
	"github.com/TrevorEdris/go-csv/pkg/log"
	gocsvschema "github.com/TrevorEdris/go-csv/pkg/schema"
	"github.com/rotisserie/eris"
)

type (
	Csv interface {
		AddRows(ctx context.Context, rowCount int) error
		Write(ctx context.Context, outputFile string) error
	}

	csvFile struct {
		schema   *gocsvschema.Schema
		rows     []row
		rowCount int
	}
)

var _ Csv = &csvFile{}

func NewCsv(schema *gocsvschema.Schema) Csv {
	return &csvFile{
		schema:   schema,
		rows:     make([]row, 0),
		rowCount: 0,
	}
}

func (c *csvFile) addRow(ctx context.Context) error {
	if c.rowCount >= c.schema.Metadata.RowCount {
		log.FromCtx(ctx).Sugar().Infof("CSV populated with %d rows; skipping", c.schema.Metadata.RowCount)
		return nil
	}

	rowToAdd := newRow(c.schema.Columns)
	c.rows = append(c.rows, rowToAdd)
	c.rowCount += 1

	return nil
}

func (c *csvFile) AddRows(ctx context.Context, rowCount int) error {
	tenPct := rowCount / 10
	for i := 0; i < rowCount; i++ {
		err := c.addRow(ctx)
		if err != nil {
			return err
		}
		if i%tenPct == 0 {
			log.FromCtx(ctx).Info(banner.New(fmt.Sprintf("%d%%", i*100/rowCount), banner.WithChar('-'), banner.WithLength(40), banner.Blue()))
		}
	}
	log.FromCtx(ctx).Info(banner.New("100%", banner.WithChar('-'), banner.WithLength(40), banner.Blue()))
	log.FromCtx(ctx).Info(banner.New(fmt.Sprintf("Created %d rows", c.rowCount), banner.WithLength(40), banner.WithChar('-'), banner.Green()))
	return nil
}

func (c *csvFile) headersAsSlice(ctx context.Context) []string {
	output := make([]string, len(c.schema.Columns))
	for i, col := range c.schema.Columns {
		if col.Order != nil {
			log.FromCtx(ctx).Sugar().Infof("Header row [%d] set to %s", *col.Order, col.Label)
			output[*col.Order] = col.Label
		} else {
			log.FromCtx(ctx).Sugar().Warnf("col [%d] %s unordered", i, col.Label)
			output[i] = col.Label
		}
	}
	return output
}

func (c *csvFile) Write(ctx context.Context, outputFile string) error {
	// Ensure the outputFile's directory exists
	outputDir := filepath.Base(outputFile)
	if outputDir != outputFile {
		os.MkdirAll(outputDir, os.ModePerm)
	}
	f, err := os.Create(outputFile)
	if err != nil {
		return eris.Wrap(err, "failed to create output file")
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	writer.Comma = c.schema.Metadata.DelimiterRune

	// Write header row
	h := c.headersAsSlice(ctx)
	log.FromCtx(ctx).Sugar().Infof("FINAL HEADER ROW: %s", strings.Join(h, ","))
	err = writer.Write(h)
	if err != nil {
		return eris.Wrap(err, "failed to write header row")
	}

	tenPct := c.schema.Metadata.RowCount / 10
	for i, row := range c.rows {
		output := row.ToStringSlice(c.schema.Columns)
		err = writer.Write(output)
		if err != nil {
			return eris.Wrapf(err, "failed to write row %d", i)
		}
		if i%tenPct == 0 {
			log.FromCtx(ctx).Info(banner.New(fmt.Sprintf("%d%%", i*100/c.schema.Metadata.RowCount), banner.WithChar('-'), banner.WithLength(40), banner.Blue()))
		}
	}
	writer.Flush()
	log.FromCtx(ctx).Sugar().Info(banner.New(fmt.Sprintf("Wrote CSV to %s", outputFile), banner.WithChar('='), banner.Green()))
	return nil
}
