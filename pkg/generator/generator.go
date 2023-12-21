package generator

import (
	"context"

	"github.com/TrevorEdris/banner"
	"github.com/TrevorEdris/go-csv/pkg/csv"
	"github.com/TrevorEdris/go-csv/pkg/log"
	"github.com/TrevorEdris/go-csv/pkg/schema"
	gocsvschema "github.com/TrevorEdris/go-csv/pkg/schema"
	"github.com/rotisserie/eris"
)

type (
	Generator interface {
		Generate(ctx context.Context, schema *gocsvschema.Schema, outputFile string) error
		GenerateFromFile(ctx context.Context, schemaFile, outputFile string) error
	}

	generator struct {
		cfg GeneratorConfig
	}

	GeneratorConfig struct{}
)

var _ Generator = &generator{}

func NewGenerator(cfg GeneratorConfig) (Generator, error) {
	return &generator{cfg}, nil
}

func (g *generator) Generate(ctx context.Context, schema *gocsvschema.Schema, outputFile string) error {
	meta := schema.Metadata
	out := csv.NewCsv(schema)
	log.FromCtx(ctx).Info(banner.New("Generating CSV", banner.WithChar('='), banner.Green()))
	err := out.AddRows(ctx, meta.RowCount)
	if err != nil {
		return eris.Wrap(err, "failed to add rows to CSV")
	}
	log.FromCtx(ctx).Info(banner.New("Writing CSV to file", banner.WithChar('='), banner.Green()))
	err = out.Write(ctx, outputFile)
	if err != nil {
		return eris.Wrap(err, "failed to write generated CSV to file")
	}
	return nil
}

func (g *generator) GenerateFromFile(ctx context.Context, schemaFile, outputFile string) error {
	log.FromCtx(ctx).Sugar().Infof("Generating CSV using schema %s, outputting to %s", schemaFile, outputFile)
	s, err := schema.NewSchemaFromFile(ctx, schemaFile)
	if err != nil {
		return eris.Wrap(err, "failed to load schema from file")
	}
	return g.Generate(ctx, s, outputFile)
}
