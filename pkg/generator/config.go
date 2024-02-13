package generator

import (
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
)

type (
	Config struct {
		Fields    []Field `validate:"required,dive,required"`
		Delimiter *string
	}
	Field struct {
		Name              string             `validate:"required"`
		Type              string             `validate:"required"`
		NumConstraint     *NumConstraint     `yaml:"numConstraint"`
		TimeConstraint    *TimeConstraint    `yaml:"timeConstraint"`
		StringConstraint  *StringConstraint  `yaml:"stringConstraint"`
		GenericConstraint *GenericConstraint `yaml:"genericConstraint"`
	}

	NumConstraint struct {
		Min *float64
		Max *float64
	}

	TimeConstraint struct {
		Min    *time.Time
		Max    *time.Time
		Format string `validate:"required"`
	}

	StringConstraint struct {
		OneOf   []string
		Pattern *string
	}

	GenericConstraint struct {
		ChanceToOmit *float64 `validate:"required,gte=0,lte=1"`
	}
)

func NewConfig(configFilename string) (Config, error) {
	if validate == nil {
		validate = validator.New()
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

	if cfg.Delimiter == nil {
		cfg.Delimiter = &DefaultDelimiter
	}

	fmt.Printf("%+v\n", cfg)

	return cfg, nil
}

func (c Config) Validate() error {
	err := validate.Struct(c)
	if err != nil {
		return err
	}

	// TODO: Custom validation of constraints

	return nil
}

func (f Field) hasValueConstraints() bool {
	return (f.NumConstraint != nil) || (f.StringConstraint != nil) || (f.TimeConstraint != nil)
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
		} else {
			return "INVALID_STRING_CONSTRAINT"
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
