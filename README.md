[![Go Report Card](https://goreportcard.com/badge/github.com/TrevorEdris/go-csv)](https://goreportcard.com/report/github.com/TrevorEdris/go-csv)
![CodeQL](https://github.com/TrevorEdris/go-csv/workflows/CodeQL/badge.svg)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![GoT](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](https://go.dev)

# go-csv
CSV generator written in Go

## Dependencies

- [Go 1.21.4](https://github.com/moovweb/gvm)

## Quick Usage

```
❯ go run main.go
 100% |█████████████████████████████████████████████████████████████████████| (5/5, 29732 it/s)
2024-02-17T10:27:52.525-0500    INFO    generator/generator.go:135      CSV generation complete {"filena
me": "out.csv"}

❯ cat out.csv
FNAME,LNAME,STREET,DOB,PICK_ONE,PATTERNITIZED,SOME_NUM
Roberta,Bartoletti,345 Creekmouth,1959-11-15,R,r0CUhz,1104
Zola,Nader,510 Capeport,1976-12-25,Y,_B_kQ4,558
Dasia,Kerluke,4154 West Landington,1969-04-01,R,2VBt23,914
Lydia,Bednar,5453 Heightsville,1967-05-18,B,ZkqQsO,551
Alta,Rau,276 New Causewayville,2003-07-29,B,hwMgWW,1145
```

**Command help**

```
❯ go run main.go -h
Usage of /var/folders/0t/w8j6n37s48g0b7kypjpbjp4w0000gn/T/go-build2512557781/b001/exe/main:
  -count int
        Number of records to generate (default 5)
  -help
        Show usage
  -input string
        Input filename (default "input/example.yaml")
  -output string
        Output filename (default "out.csv")
```

## Configuration

See [input/example.yaml](./input/example.yaml) for an example configuration file with annotations.

### Available Constraints

**Generic Constraint**

A constraint that can be applied to any field and can be combined with both `type` and `xxxConstraint`.

```yaml
genericConstraint:
  chanceToOmit: 0.1 # 10% chance of an empty string
```

`chanceToOmit` must be a decimal between 0 and 1. When a randomly generated decimal is **less than** the configured `chanceToOmit`, the value for that cell will be an empty string.

**Num Constraint**

```yaml
numConstraint:
  min: 420
  max: 1337
```

Generate an integer between the specified `min` and `max` values (both required).

**Time Constraint**

```yaml
timeConstraint:
  format: 2006-01-02
  min: 1945-01-01
  max: 2005-02-12
```

Generate a timestamp with the specified `format` between the specified `min` and `max` times.

Note: The time value used in the `format` field must match the Golang reference time `Mon, 02 Jan 2006 15:04:05 MST`, see https://www.geeksforgeeks.org/time-formatting-in-golang/.

**String Constraint**

```yaml
stringConstraint:
  pattern: '^[a-zA-Z0-9_]{6}$'
  oneOf:
    - R
    - B
    - Y
```

Must specify either `pattern`, a regex pattern, or `oneOf`, a list of possible values, but cannot specify both. When using `oneOf`, a random entry in the list will be chosen as the value. All entries have equal weight. For differently weighted entries, either register a custom value function or provide multiple entries in `oneOf` corresponding to the desired weights.

```yaml
stringConstraint:
  oneOf:
    - A
    - A
    - A
    - B
```

The above configuration will choose `A` 75% of the time and `B` 25% of the time.

### Available Types

- `EMPTY`
- `UINT64`
- `UINT32`
- `UINT8`
- `FIRSTNAME`
- `LASTNAME`
- `STREET`
- `CITY`
- `STATE`
- `ZIP`
- `PHONE`
- `EMAIL`
- `COMPANY`
- `YES_OR_NO`
- `CONSISTENTLY_INCREASING_ID`

## Registering a custom value function

A custom function can be registered to a generator, exposing that function for use via the configuration yaml file. The function must match the `generator.RandomValueFunc` signature

```go
type RandomValueFunc func(p FieldParams) string
```

```go
package main

import (
	"fmt"
	"log"

	"github.com/TrevorEdris/go-csv/pkg/generator"
)

func main() {
	g, err := generator.NewGenerator("input/example.yaml", "out.csv")
	if err != nil {
		log.Fatalf("Failed to initialize generator: %s", err.Error())
	}

	err = g.RegisterNewValueFunction("ROW_NUMBER", myOwnValueFunction)
	if err != nil {
		log.Fatalf("Failed to register value function: %s", err.Error())
	}

	err = g.Generate(*count)
	if err != nil {
		log.Fatalf("Failed to generate CSV: %s", err.Error())
	}
}

func myOwnValueFunction(p generator.FieldParams) string {
	return fmt.Sprintf("ROW_NUMBER-%d", p.RowNumber)
}
```

```yaml
fields:
- name: RowNum
  type: ROW_NUMBER # The value passed in to `g.RegisterNewValueFunction`
```

## Contributing

Refer to the [Contribution Guidelines](./CONTRIBUTING.md) for instructions on contributing to this project.
