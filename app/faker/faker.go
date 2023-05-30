package faker

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

type (
	Source string

	NumericConstraint struct {
		Min float64 `yaml:"min"`
		Max float64 `yaml:"max"`
	}

	StringConstraint struct {
		Regex string   `yaml:"regex"`
		OneOf []string `yaml:"oneOf"`
	}

	TimestampConstraint struct {
		Start  time.Time `yaml:"after"`
		End    time.Time `yaml:"before"`
		Format string    `yaml:"format"`
	}

	GeneralConstraint struct {
		SkipChance float64 `yaml:"skipChance"`
	}
)

const (
	dobFormat = "2006-01-02"
)

var (
	TypeFirstName Source = "FIRST_NAME"
	TypeLastName  Source = "LAST_NAME"
	TypeEmail     Source = "EMAIL"
	TypeUUID      Source = "UUID"
	TypeTimestamp Source = "TIMESTAMP"
	TypeInt       Source = "INTEGER"
	TypeString    Source = "STRING"
	TypeBool      Source = "BOOL"
	TypeCompany   Source = "COMPANY"
	TypePhone     Source = "PHONE"
	TypeStreet    Source = "STREET"
	TypeCity      Source = "CITY"
	TypeState     Source = "STATE"
	TypeZip       Source = "ZIP"
	TypeCountry   Source = "COUNTRY"
	TypeMonth     Source = "MONTH"
	TypeBeerName  Source = "BEER_NAME"
	TypeBeerStyle Source = "BEER_STYLE"
	TypeYesNo     Source = "YES_NO"

	// The list of "faker" values the application supports.
	// Note: This is just a subset of the fakers available in https://github.com/brianvoe/gofakeit
	//       in addition to custom ones.
	AllSources = []Source{
		TypeFirstName, TypeLastName, TypeBool, TypeEmail,
		TypeUUID, TypeTimestamp, TypeInt, TypeString,
		TypeBool, TypePhone, TypeStreet, TypeZip, TypeCompany,
		TypeMonth, TypeBeerName, TypeBeerStyle, TypeCity,
		TypeState, TypeYesNo, TypeCountry,
	}
)

func FakeValue(
	faker Source,
	sc *StringConstraint,
	nc *NumericConstraint,
	tc *TimestampConstraint,
	gc *GeneralConstraint,
) string {
	v := ""
	switch faker {
	case TypeFirstName:
		v = gofakeit.FirstName()
	case TypeLastName:
		v = gofakeit.LastName()
	case TypeEmail:
		v = gofakeit.Email()
	case TypeUUID:
		v = gofakeit.UUID()
	case TypeTimestamp:
		if tc != nil {
			zeroTime := time.Time{}
			if tc.Start != zeroTime && tc.End != zeroTime {
				v = gofakeit.DateRange(tc.Start, tc.End).Format(tc.Format)
			} else {
				v = gofakeit.Date().Format(tc.Format)
			}
		} else {
			v = gofakeit.Date().Format(dobFormat)
		}
	case TypeInt:
		if nc != nil {
			v = fmt.Sprintf("%d", gofakeit.Number(int(nc.Min), int(nc.Max)))
		} else {
			v = fmt.Sprintf("%d", gofakeit.Int64())
		}
	case TypeString:
		if sc != nil {
			if len(sc.OneOf) != 0 {
				v = gofakeit.RandomString(sc.OneOf)
			} else if sc.Regex != "" {
				v = gofakeit.Regex(sc.Regex)
			} else {
				v = "INVALID_STRING_CONSTRAINT"
			}
		} else {
			v = gofakeit.LoremIpsumWord()
		}
	case TypeBool:
		v = fmt.Sprintf("%v", gofakeit.Bool())
	case TypeCompany:
		v = gofakeit.Company()
	case TypePhone:
		v = gofakeit.Phone()
	case TypeStreet:
		v = gofakeit.Street()
	case TypeCity:
		v = gofakeit.City()
	case TypeState:
		v = gofakeit.State()
	case TypeZip:
		v = gofakeit.Zip()
	case TypeMonth:
		v = gofakeit.MonthString()
	case TypeBeerName:
		v = gofakeit.BeerName()
	case TypeBeerStyle:
		v = gofakeit.BeerStyle()
	case TypeYesNo:
		v = gofakeit.RandomString([]string{"Y", "N"})
	case TypeCountry:
		v = gofakeit.Country()
	default:
	}

	if gc != nil {
		if rand.Float64() < gc.SkipChance {
			return ""
		}
	}

	return v
}

func (nc *NumericConstraint) Validate() error {
	// Automatically reconfigure the min/max if they're swapped
	if nc.Min < nc.Max {
		nc.Min, nc.Max = nc.Max, nc.Min
	}

	return nil
}

func (sc *StringConstraint) Validate() error {
	if len(sc.OneOf) == 0 && sc.Regex == "" {
		return fmt.Errorf("stringConstraint must specify either regex or oneOf properties")
	}

	return nil
}

func (tc *TimestampConstraint) Validate() error {
	zeroTime := time.Time{}
	if (tc.Start == zeroTime && tc.End != zeroTime) || (tc.Start != zeroTime && tc.End == zeroTime) {
		return fmt.Errorf("timestampConstraint must specify BOTH or NEITHER before and after")
	}

	if tc.Format == "" {
		return fmt.Errorf("timestampConstraint format must be non-empty")
	}

	return nil
}

func (gc *GeneralConstraint) Validate() error {
	if gc.SkipChance < 0.0 {
		return fmt.Errorf("generalConstraint skipChance must be >= 0.0: %.03f", gc.SkipChance)
	}
	return nil
}
