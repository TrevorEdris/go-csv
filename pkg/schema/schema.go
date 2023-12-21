package schema

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/TrevorEdris/banner"
	"github.com/TrevorEdris/go-csv/pkg/errors"
	"github.com/TrevorEdris/go-csv/pkg/log"
	"github.com/go-playground/validator/v10"
	"github.com/rotisserie/eris"
	"gopkg.in/yaml.v3"
)

type (
	Schema struct {
		Columns  Headers  `yaml:"columns"`
		Metadata Metadata `yaml:"metadata"`
	}

	Metadata struct {
		Delimiter     string        `yaml:"delimiter"`
		DelimiterRune rune          `yaml:""` // Only used internally; should not be configurable
		RowCount      int           `yaml:"rowCount" validate:"gte=1"`
		MultipleFiles MultipleFiles `yaml:"multipleFiles"`
	}

	MultipleFiles struct {
		Enabled   bool `yaml:"enabled"`
		FileCount int  `yaml:"fileCount" validate:"gte=1"`
	}
)

var (
	validate *validator.Validate

	minDoB = time.Date(1930, 1, 1, 0, 0, 0, 0, time.UTC)
	maxDoB = time.Now()

	zero    = 0
	one     = 1
	two     = 2
	example = Schema{
		Columns: []Column{
			{
				Label:  "First Name",
				Source: "FIRST_NAME",
				Order:  &one,
			},
			{
				Label:  "Last Name",
				Source: "LAST_NAME",
				Order:  &two,
			},
			{
				Label:  "Group",
				Source: "INTEGER",
				Order:  &zero,
				NumericConstraint: &NumericConstraint{
					Min: 1,
					Max: 9,
				},
			},
			{
				Label:  "Relationship",
				Source: "STRING",
				StringConstraint: &StringConstraint{
					OneOf: []string{"EMPLOYEE", "SPOUSE", "DEPENDENT"},
				},
			},
			{
				Label:  "Date of Birth",
				Source: "TIMESTAMP",
				TimestampConstraint: &TimestampConstraint{
					Start:  &minDoB,
					End:    &maxDoB,
					Format: "2006-01-02",
				},
			},
			{
				Label:  "Potentially Empty Field",
				Source: "UUID",
				GeneralConstraint: &GeneralConstraint{
					SkipChance: 0.5,
				},
			},
		},
		Metadata: Metadata{
			Delimiter:     ",",
			DelimiterRune: ',',
			RowCount:      10,
			MultipleFiles: MultipleFiles{
				Enabled:   false,
				FileCount: 1,
			},
		},
	}
)

func CreateExample(ctx context.Context, outputFile string) error {
	// Check if outputFile is also specifying a directory
	outputDir := ""
	fname := filepath.Base(outputFile)
	if fname != outputFile {
		outputDir, _ = strings.CutSuffix(outputFile, fname)
		// Strip the `/` or `\` from the path
		outputDir = outputDir[:len(outputDir)-1]
	}
	// If the output dir does not already exist, create it
	_, err := os.Stat(outputDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(outputDir, os.ModePerm)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	f, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	yamlData, err := yaml.Marshal(&example)
	if err != nil {
		return err
	}
	_, err = f.Write(yamlData)
	if err != nil {
		return err
	}
	log.FromCtx(ctx).Sugar().Infof("Created %s\n", outputFile)
	return nil
}

func ValidateSchemaFile(ctx context.Context, schemaFile string) error {
	validate = validator.New()

	schema, err := NewSchemaFromFile(ctx, schemaFile)
	if err != nil {
		return err
	}

	err = schema.Validate()
	if err != nil {
		if !errors.Is(errors.InvalidSchemaError, err) {
			return err
		}
		log.FromCtx(ctx).Sugar().Warnf("Invalid schema file %s; %s", schemaFile, err.Error())
		return nil
	}

	log.FromCtx(ctx).Info("Schema is valid")
	return nil
}

func NewSchemaFromFile(ctx context.Context, schemaFile string) (*Schema, error) {
	bytes, err := os.ReadFile(schemaFile)
	if err != nil {
		return nil, err
	}

	schema := &Schema{}
	err = yaml.Unmarshal(bytes, schema)
	if err != nil {
		return nil, err
	}

	log.FromCtx(ctx).Info(banner.New("Determining column order", banner.WithChar('='), banner.Green()))

	// Ensure proper ordering by filling-in any missing orders
	orderedCols := make(Headers, len(schema.Columns))
	colsWithNoOrder := make(Headers, 0)

	for i, col := range schema.Columns {
		if col.Order == nil {
			colsWithNoOrder = append(colsWithNoOrder, col)
			log.FromCtx(ctx).Sugar().Infof("Input [%d] %s has no order", i, col.Label)
			continue
		}
		orderedCols[*col.Order] = col
		log.FromCtx(ctx).Sugar().Infof("Input [%d] %s specified order %d", i, col.Label, *col.Order)
	}

	nextOrder := 0
	emptyCol := Column{}
	for i, col := range orderedCols {
		if col != emptyCol {
			log.FromCtx(ctx).Sugar().Infof("--- Output [%d] already set to %s", i, col.Label)
			continue
		}
		log.FromCtx(ctx).Sugar().Infof("Output [%d] unset; Setting to %s", i, colsWithNoOrder[nextOrder].Label)
		newI := i
		colsWithNoOrder[nextOrder].Order = &newI
		orderedCols[i] = colsWithNoOrder[nextOrder]
		nextOrder += 1
	}

	schema.Columns = orderedCols

	return schema, nil
}

func (s *Schema) Validate() error {
	// Struct-level validation
	err := validate.Struct(s)
	if err != nil {
		return err
	}

	// Custom validation
	err = s.Columns.Validate()
	if err != nil {
		return err
	}

	err = s.Metadata.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (m *Metadata) Validate() error {
	if m.RowCount <= 0 {
		return eris.Wrapf(errors.InvalidSchemaError, "metadata rowcount invalid (%d): must be > 0", m.RowCount)
	}

	if m.Delimiter == "" {
		m.DelimiterRune = ','
	} else if len(m.Delimiter) == 1 {
		m.DelimiterRune = rune(m.Delimiter[0])
	} else {
		return eris.Wrapf(errors.InvalidSchemaError, "metadata delimiter (%s) must be a string of length 1", m.Delimiter)
	}

	return nil
}
