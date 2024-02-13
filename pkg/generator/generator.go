package generator

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

const (
	UnsupportedFieldType = "UNSUPPORTED_FIELD_TYPE"
)

var (
	// Base set of functions that don't require a custom random seed
	baseFakerFuncs = map[string]randomValueFunc{
		"UINT64":    func() string { return strconv.FormatUint(uint64(gofakeit.Uint64()), 10) },
		"UINT32":    func() string { return strconv.FormatUint(uint64(gofakeit.Uint32()), 10) },
		"UINT8":     func() string { return strconv.FormatUint(uint64(gofakeit.Uint8()), 10) },
		"EMPTY":     func() string { return "" },
		"FIRSTNAME": gofakeit.FirstName,
		"LASTNAME":  gofakeit.LastName,
		"STREET":    gofakeit.Street,
	}

	ErrKeyConflict = errors.New("faker with specified key already registered")
)

type (
	Generator struct {
		cfg       Config
		output    string
		delimiter string
		r         *rand.Rand
		fakers    map[string]randomValueFunc
	}

	randomValueFunc func() string
)

func NewGenerator(input, output, delimiter string) (*Generator, error) {
	cfg, err := NewConfig(input)
	if err != nil {
		return nil, err
	}
	if *cfg.Delimiter != delimiter {
		cfg.Delimiter = &delimiter
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	g := &Generator{
		cfg:       cfg,
		output:    output,
		delimiter: delimiter,
		r:         r,
		fakers:    baseFakerFuncs,
	}
	// Register all custom functions that rely on the Generator's random source
	_ = g.RegisterNewFaker("yes_or_no", g.yesOrNo)
	return g, nil
}

func (g *Generator) RegisterNewFaker(key string, f randomValueFunc) error {
	key = strings.ToUpper(key)
	if _, exists := g.fakers[key]; exists {
		return ErrKeyConflict
	}
	g.fakers[key] = f
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

	for i, field := range g.cfg.Fields {
		if i != 0 {
			fmt.Fprintf(writer, "%s", g.delimiter)
		}
		fmt.Fprintf(writer, "%s", field.Name)
	}
	fmt.Fprintln(writer)

	for i := 0; i < count; i++ {
		for j, field := range g.cfg.Fields {
			if j != 0 {
				fmt.Fprintf(writer, g.delimiter)
			}
			fmt.Fprintf(writer, "%s", g.randomValue(field))
		}
		fmt.Fprintln(writer)
	}

	return nil
}

func (g *Generator) randomValue(field Field) string {
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
	return g.wrap(fakerFunc())
}

func (g *Generator) wrap(s string) string {
	if strings.Contains(s, g.delimiter) {
		return fmt.Sprintf("\"%s\"", s)
	}
	return s
}

func (g *Generator) shouldOmit(f Field) bool {
	if f.GenericConstraint != nil {
		return g.r.Float64() > *f.GenericConstraint.ChanceToOmit
	}
	return false
}

func (g *Generator) yesOrNo() string {
	if g.r.Float64() > 0.5 {
		return "y"
	}
	return "n"
}
