package faker

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/TrevorEdris/go-csv/pkg/schema"
	"github.com/brianvoe/gofakeit/v6"
)

type (
	Source string
)

const (
	dobFormat = "2006-01-02"
)

var (
	SourceFirstName Source = "FIRST_NAME"
	SourceLastName  Source = "LAST_NAME"
	SourceEmail     Source = "EMAIL"
	SourceUUID      Source = "UUID"
	SourceTimestamp Source = "TIMESTAMP"
	SourceInt       Source = "INTEGER"
	SourceString    Source = "STRING"
	SourceBool      Source = "BOOL"
	SourceCompany   Source = "COMPANY"
	SourcePhone     Source = "PHONE"
	SourceStreet    Source = "STREET"
	SourceCity      Source = "CITY"
	SourceState     Source = "STATE"
	SourceZip       Source = "ZIP"
	SourceCountry   Source = "COUNTRY"
	SourceMonth     Source = "MONTH"
	SourceBeerName  Source = "BEER_NAME"
	SourceBeerStyle Source = "BEER_STYLE"
	SourceYesNo     Source = "YES_NO"

	// The list of "faker" values the application supports.
	// Note: This is just a subset of the fakers available in https://github.com/brianvoe/gofakeit
	//       in addition to custom ones.
	AllSources = []Source{
		SourceFirstName, SourceLastName, SourceBool, SourceEmail,
		SourceUUID, SourceTimestamp, SourceInt, SourceString,
		SourceBool, SourcePhone, SourceStreet, SourceZip, SourceCompany,
		SourceMonth, SourceBeerName, SourceBeerStyle, SourceCity,
		SourceState, SourceYesNo, SourceCountry,
	}
)

func FakeValue(
	faker Source,
	sc *schema.StringConstraint,
	nc *schema.NumericConstraint,
	tc *schema.TimestampConstraint,
	gc *schema.GeneralConstraint,
) string {
	v := ""
	switch faker {
	case SourceFirstName:
		v = gofakeit.FirstName()
	case SourceLastName:
		v = gofakeit.LastName()
	case SourceEmail:
		v = gofakeit.Email()
	case SourceUUID:
		v = gofakeit.UUID()
	case SourceTimestamp:
		if tc != nil {
			zeroTime := time.Time{}
			if *tc.Start != zeroTime && *tc.End != zeroTime {
				v = gofakeit.DateRange(*tc.Start, *tc.End).Format(tc.Format)
			} else {
				v = gofakeit.Date().Format(tc.Format)
			}
		} else {
			v = gofakeit.Date().Format(dobFormat)
		}
	case SourceInt:
		if nc != nil {
			v = fmt.Sprintf("%d", gofakeit.Number(int(nc.Min), int(nc.Max)))
		} else {
			v = fmt.Sprintf("%d", gofakeit.Int64())
		}
	case SourceString:
		if sc != nil {
			if len(sc.OneOf) != 0 {
				v = gofakeit.RandomString(sc.OneOf)
			} else if *sc.Regex != "" {
				v = gofakeit.Regex(*sc.Regex)
			} else {
				v = "INVALID_STRING_CONSTRAINT"
			}
		} else {
			v = gofakeit.LoremIpsumWord()
		}
	case SourceBool:
		v = fmt.Sprintf("%v", gofakeit.Bool())
	case SourceCompany:
		v = gofakeit.Company()
	case SourcePhone:
		v = gofakeit.Phone()
	case SourceStreet:
		v = gofakeit.Street()
	case SourceCity:
		v = gofakeit.City()
	case SourceState:
		v = gofakeit.State()
	case SourceZip:
		v = gofakeit.Zip()
	case SourceMonth:
		v = gofakeit.MonthString()
	case SourceBeerName:
		v = gofakeit.BeerName()
	case SourceBeerStyle:
		v = gofakeit.BeerStyle()
	case SourceYesNo:
		v = gofakeit.RandomString([]string{"Y", "N"})
	case SourceCountry:
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
