package schema

import (
	"time"

	"github.com/TrevorEdris/go-csv/pkg/errors"
	"github.com/rotisserie/eris"
)

type (
	NumericConstraint struct {
		Min float64 `yaml:"min"`
		Max float64 `yaml:"max"`
	}

	StringConstraint struct {
		Regex *string  `yaml:"regex,omitempty"`
		OneOf []string `yaml:"oneOf,omitempty"`
	}

	TimestampConstraint struct {
		Start  *time.Time `yaml:"after,omitempty"`
		End    *time.Time `yaml:"before,omitempty"`
		Format string     `yaml:"format"`
	}

	GeneralConstraint struct {
		SkipChance float64 `yaml:"skipChance"`
	}
)

func (nc *NumericConstraint) Validate() error {
	// Automatically reconfigure the min/max if they're swapped
	if nc.Min < nc.Max {
		nc.Min, nc.Max = nc.Max, nc.Min
	}

	return nil
}

func (sc *StringConstraint) Validate() error {
	if len(sc.OneOf) == 0 && *sc.Regex == "" {
		return eris.Wrap(errors.InvalidSchemaError, "stringConstraint must specify either regex or oneOf properties")
	}

	return nil
}

func (tc *TimestampConstraint) Validate() error {
	zeroTime := time.Time{}
	if (*tc.Start == zeroTime && *tc.End != zeroTime) || (*tc.Start != zeroTime && *tc.End == zeroTime) {
		return eris.Wrap(errors.InvalidSchemaError, "timestampConstraint must specify BOTH or NEITHER before and after")
	}

	if tc.Format == "" {
		return eris.Wrap(errors.InvalidSchemaError, "timestampConstraint format must be non-empty")
	}

	return nil
}

func (gc *GeneralConstraint) Validate() error {
	if gc.SkipChance < 0.0 {
		return eris.Wrapf(errors.InvalidSchemaError, "generalConstraint skipChance must be >= 0.0: %.03f", gc.SkipChance)
	}
	return nil
}
