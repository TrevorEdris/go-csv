package generator

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"
)

var (
	DefaultDelimiter = ","
	DefaultCount     = 5

	validate *validator.Validate

	ErrInvalidConstraint = errors.New("invalid constraint")
)

type (
	Config struct {
		Fields    []*Field `validate:"required,dive,required"`
		Delimiter *string  `json:"delimiter,omitempty"`
	}

	Field struct {
		Name              string `validate:"required"`
		Type              string `json:"type,omitempty"`
		params            FieldParams
		NumConstraint     *NumConstraint     `yaml:"numConstraint" json:"numConstraint,omitempty"`
		TimeConstraint    *TimeConstraint    `yaml:"timeConstraint" json:"timeConstraint,omitempty"`
		StringConstraint  *StringConstraint  `yaml:"stringConstraint" json:"stringConstraint,omitempty"`
		GenericConstraint *GenericConstraint `yaml:"genericConstraint" json:"genericConstraint,omitempty"`
	}

	FieldParams struct {
		RowNumber int
	}

	NumConstraint struct {
		Min *float64 `json:"min,omitempty"`
		Max *float64 `json:"max,omitempty"`
	}

	TimeConstraint struct {
		Min    *time.Time `json:"min,omitempty"`
		Max    *time.Time `json:"max,omitempty"`
		Format string     `validate:"required"`
	}

	StringConstraint struct {
		OneOf   []string `yaml:"oneOf" validate:"dive,required" json:"oneOf,omitempty"`
		Pattern *string  `json:"pattern,omitempty"`
	}

	GenericConstraint struct {
		ChanceToOmit *float64 `validate:"required,gte=0,lte=1" json:"chanceToOmit,omitempty"`
	}
)

func NewConfig(configFilename string) (Config, error) {
	if validate == nil {
		validate = validator.New()
		validate.RegisterStructValidation(configStructLevelValidation, Config{})
		validate.RegisterStructValidation(fieldStructLevelValidation, Field{})
		validate.RegisterStructValidation(numConstraintStructLevelValidation, NumConstraint{})
		validate.RegisterStructValidation(timeConstraintStructLevelValidation, TimeConstraint{})
		validate.RegisterStructValidation(stringConstraintStructLevelValidation, StringConstraint{})
		validate.RegisterStructValidation(genericConstraintStructLevelValidation, GenericConstraint{})
	}

	data, err := os.ReadFile(configFilename)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return Config{}, err
	}

	err = validate.Struct(cfg)
	if err != nil {
		return Config{}, err
	}

	if cfg.Delimiter == nil {
		cfg.Delimiter = &DefaultDelimiter
	}

	return cfg, nil
}

func (c Config) AsJson() (string, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (c Config) AsFormattedJson() (string, error) {
	b, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (f *Field) registerParams(p FieldParams) {
	f.params = p
}

func (f *Field) hasValueConstraints() bool {
	return (f.NumConstraint != nil) || (f.StringConstraint != nil) || (f.TimeConstraint != nil)
}

func (f *Field) hasMultipleValueConstraints() bool {
	c := 0
	if f.NumConstraint != nil {
		c++
	}
	if f.StringConstraint != nil {
		c++
	}
	if f.TimeConstraint != nil {
		c++
	}
	return c > 1
}

func (f Field) valueFromConstraint() string {
	if f.NumConstraint != nil {
		return fmt.Sprintf("%d", gofakeit.Number(int(*f.NumConstraint.Min), int(*f.NumConstraint.Max)))
	}

	if f.StringConstraint != nil {
		if len(f.StringConstraint.OneOf) != 0 {
			return gofakeit.RandomString(f.StringConstraint.OneOf)
		} else if f.StringConstraint.Pattern != nil {
			return gofakeit.Regex(*f.StringConstraint.Pattern)
		}
	}

	if f.TimeConstraint != nil {
		zeroTime := time.Time{}
		if *f.TimeConstraint.Min != zeroTime && *f.TimeConstraint.Max != zeroTime {
			t := gofakeit.DateRange(*f.TimeConstraint.Min, *f.TimeConstraint.Max)
			return t.Format(f.TimeConstraint.Format)
		} else {
			return gofakeit.Date().Format(f.TimeConstraint.Format)
		}
	}

	return ""
}

func configStructLevelValidation(sl validator.StructLevel) {
	c := sl.Current().Interface().(Config)

	if c.Delimiter != nil && *c.Delimiter == "" {
		sl.ReportError(c.Delimiter, "delimiter", "Delimiter", "nonemptyDelimiter", "")
	} else if c.Delimiter != nil && len(*c.Delimiter) > 1 {
		sl.ReportError(c.Delimiter, "delimiter", "Delimiter", "singleCharDelimiter", "")
	}
}

func fieldStructLevelValidation(sl validator.StructLevel) {
	f := sl.Current().Interface().(Field)

	if f.Type == "" && !f.hasValueConstraints() {
		sl.ReportError(f.Type, "type", "Type", "typeRequiredWhenNoConstraints", "")
	} else if f.Type != "" && f.hasValueConstraints() {
		sl.ReportError(f.Type, "type", "Type", "typeOrConstraints", "")
	} else if f.hasMultipleValueConstraints() {
		sl.ReportError(f.Type, "type", "Type", "onlyUseOneConstraint", "")
	}
}

func numConstraintStructLevelValidation(sl validator.StructLevel) {
	c := sl.Current().Interface().(NumConstraint)

	if c.Min != nil && c.Max != nil && *c.Min >= *c.Max {
		sl.ReportError(c.Min, "min", "Min", "minLTmax", "")
		sl.ReportError(c.Max, "max", "Max", "minLTmax", "")
	}
}

func timeConstraintStructLevelValidation(sl validator.StructLevel) {
	c := sl.Current().Interface().(TimeConstraint)

	if c.Min != nil && c.Max != nil && c.Min.After(*c.Max) {
		sl.ReportError(c.Min, "min", "Min", "minLTmax", "")
		sl.ReportError(c.Max, "max", "Max", "minLTmax", "")
	}
}

func stringConstraintStructLevelValidation(sl validator.StructLevel) {
	c := sl.Current().Interface().(StringConstraint)

	if (len(c.OneOf) > 0 && c.Pattern != nil && *c.Pattern != "") || (len(c.OneOf) == 0 && (c.Pattern == nil || *c.Pattern == "")) {
		sl.ReportError(c.OneOf, "oneOf", "OneOf", "oneOfOrPattern", "")
		sl.ReportError(c.Pattern, "pattern", "Pattern", "oneOfOrPattern", "")
	} else if c.Pattern != nil && *c.Pattern == "" {
		sl.ReportError(c.Pattern, "pattern", "Pattern", "nonemptyPattern", "")
	}
}

func genericConstraintStructLevelValidation(sl validator.StructLevel) {
	// Placeholder for any custom validation
}
