package generator

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/TrevorEdris/go-csv/pkg/logctx"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/schollz/progressbar/v3"
	"go.uber.org/zap"
)

const (
	UnsupportedFieldType = "UNSUPPORTED_FIELD_TYPE"
)

var (
	// Base set of functions that don't require a custom random seed
	baseValueFunctions = map[string]RandomValueFunc{
		"EMPTY":     func(_ FieldParams) string { return "" },
		"UINT64":    func(_ FieldParams) string { return strconv.FormatUint(uint64(gofakeit.Uint64()), 10) },
		"UINT32":    func(_ FieldParams) string { return strconv.FormatUint(uint64(gofakeit.Uint32()), 10) },
		"UINT8":     func(_ FieldParams) string { return strconv.FormatUint(uint64(gofakeit.Uint8()), 10) },
		"FIRSTNAME": func(_ FieldParams) string { return gofakeit.FirstName() },
		"LASTNAME":  func(_ FieldParams) string { return gofakeit.LastName() },
		"STREET":    func(_ FieldParams) string { return gofakeit.Street() },
		"CITY":      func(_ FieldParams) string { return gofakeit.City() },
		"STATE":     func(_ FieldParams) string { return gofakeit.State() },
		"ZIP":       func(_ FieldParams) string { return gofakeit.Zip() },
		"PHONE":     func(_ FieldParams) string { return gofakeit.Phone() },
		"EMAIL":     func(_ FieldParams) string { return gofakeit.Email() },
		"COMPANY":   func(_ FieldParams) string { return gofakeit.Company() },
	}

	ErrKeyConflict = errors.New("value function with specified key already registered")
)

type (
	Generator struct {
		cfg    Config
		output string
		r      *rand.Rand
		fakers map[string]RandomValueFunc
		log    *zap.Logger
	}

	RandomValueFunc func(p FieldParams) string
)

// TODO: Maybe add support for NewGenerator(cfg Config)
func NewGenerator(input, output string) (*Generator, error) {
	cfg, err := NewConfig(input)
	if err != nil {
		return nil, err
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	g := &Generator{
		cfg:    cfg,
		output: output,
		r:      r,
		fakers: map[string]RandomValueFunc{},
		// TODO: Support passing in a logger via options pattern
		log: logctx.FromCtx(context.Background()),
	}

	// Register all base functions
	for key, f := range baseValueFunctions {
		err := g.RegisterNewValueFunction(key, f)
		if err != nil {
			return nil, err
		}
	}

	// Register all custom functions
	_ = g.RegisterNewValueFunction("yes_or_no", g.yesOrNo)
	_ = g.RegisterNewValueFunction("consistently_increasing_id", g.consistentlyIncreasingID)

	return g, nil
}

func (g *Generator) RegisterNewValueFunction(key string, f RandomValueFunc) error {
	key = strings.ToUpper(key)
	if _, exists := g.fakers[key]; exists {
		g.log.Error("Failed to register value function", zap.String("key", key), zap.Error(ErrKeyConflict))
		return ErrKeyConflict
	}
	g.fakers[key] = f
	g.log.Debug("Successfully registered value function", zap.String("key", key))
	return nil
}

func (g *Generator) Generate(count int) error {
	f, err := os.Create(g.output)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	defer writer.Flush()

	// First, write the header column
	for i, field := range g.cfg.Fields {
		if i != 0 {
			fmt.Fprintf(writer, "%s", *g.cfg.Delimiter)
		}
		fmt.Fprintf(writer, "%s", field.Name)
	}
	fmt.Fprintln(writer)

	// Then, for each cell, generate the random value
	pb := progressbar.Default(int64(count))
	for i := 0; i < count; i++ {
		for j, field := range g.cfg.Fields {
			// Value calculation _may_ depend on the row number or other dynamic attributes
			field.registerParams(FieldParams{i})

			if j != 0 { // Include the delimiter for all cells after the first cell
				fmt.Fprintf(writer, *g.cfg.Delimiter)
			}
			fmt.Fprintf(writer, "%s", g.randomValue(field))
		}
		fmt.Fprintln(writer)
		pb.Add(1)
	}

	g.log.Info("CSV generation complete", zap.String("filename", g.output))

	return nil
}

func (g *Generator) randomValue(field *Field) string {
	// Return empty string if field was randomly selected to be omitted
	if g.shouldOmit(field) {
		return ""
	}
	// If the field has some constraints, generate the value based on those constraints
	if field.hasValueConstraints() {
		return g.wrap(field.valueFromConstraint())
	}
	// Otherwise, generate the value from the registered faker function
	fieldType := strings.ToUpper(field.Type)
	fakerFunc, exists := g.fakers[fieldType]
	if !exists {
		return g.wrap(UnsupportedFieldType)
	}
	return g.wrap(fakerFunc(field.params))
}

// wrap will surround the value in "" if the value contains the delimiter.
// TODO: This will break if the value also contains ".
func (g *Generator) wrap(s string) string {
	if strings.Contains(s, *g.cfg.Delimiter) {
		return fmt.Sprintf("\"%s\"", s)
	}
	return s
}

func (g *Generator) shouldOmit(f *Field) bool {
	if f.GenericConstraint != nil {
		return g.r.Float64() < *f.GenericConstraint.ChanceToOmit
	}
	return false
}

// yesOrNo returns `y` or `n` 50% of the time respectively.
// This could be done via a `StringConstraint` with the `OneOf`
// property, but if a weight other than 50-50 is desired,
// a custom function like this would be necessary.
func (g *Generator) yesOrNo(_ FieldParams) string {
	if g.r.Float64() > 0.5 {
		return "y"
	}
	return "n"
}

func (g *Generator) consistentlyIncreasingID(fp FieldParams) string {
	possibleIDs := []string{
		"0000001",
		"0000002",
		"0000003",
		"0000004",
		"0000005",
		"0000006",
		"0000007",
		"0000008",
		"0000009",
		"0000010",
		"0000011",
		"0000012",
		"0000013",
		"0000014",
		"0000015",
		"0000016",
		"0000017",
		"0000018",
		"0000019",
		"0000020",
	}
	return possibleIDs[fp.RowNumber%len(possibleIDs)]
}
