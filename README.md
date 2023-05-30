[![Go Report Card](https://goreportcard.com/badge/github.com/TrevorEdris/go-csv)](https://goreportcard.com/report/github.com/TrevorEdris/go-csv)
![CodeQL](https://github.com/TrevorEdris/go-csv/workflows/CodeQL/badge.svg)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![GoT](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](https://go.dev)

# go-csv
CSV generator written in Go

## Dependencies

- [Go 1.20](https://github.com/moovweb/gvm)

## Quick Usage

```
❯ go run cmd/gocsv/main.go generate  --input input/default.yaml --output output.csv
2023-05-30T10:18:23.846-0400 INFO Generating CSV
2023-05-30T10:18:23.847-0400 INFO Completed in 200.4µs

❯ cat output.csv
Group,First Name,Last Name,Country of Origin,Birthday,Some Text
6,Karen,Labadie,United Kingdom of Great Britain and Northern Ireland,2018-01-30,9Ewxfa
7,Blaise,Greenfelder,Benin,1954-08-06,lllD3t
7,Wilburn,Osinski,Haiti,1920-07-29,Tvhz_l
4,Nicolas,Gerhold,"Moldova, Republic of",1922-07-29,FkwWCT
5,Oleta,Feil,Bahrain,1972-10-26,qAOdct
```

**Generate command help**

```
❯ go run cmd/gocsv/main.go generate
Error: required flag(s) "input" not set
Usage:
  gocsv generate [flags]

Flags:
  -h, --help              help for generate
      --input string      Relative filename to the configuration yaml file (default "input/default.yaml")
      --logLevel string   Log level (DEBUG|INFO|WARN) (default "INFO")
      --output string     Filename to write output to (default "output/default.csv")

required flag(s) "input" not set
exit status 1
```

## Installation

Use the provided `Makefile` to assist with installation.

```
❯ make install
go install ./cmd/gocsv
```

The `gocsv` command should now be accessible via CLI.

```
❯ gocsv -h
Interact with CSV files

Usage:
  gocsv [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  generate    Generate CSV file with randomized data
  help        Help about any command
  version     Print the version of gocsv

Flags:
  -h, --help   help for gocsv

Use "gocsv [command] --help" for more information about a command.
```

## Configuration

See [CONFIGURATION.md](./CONFIGURATION.md) for the configuration file schema.

## Development

### Requirements

- Docker
- Docker Compose

The provided `Makefile` has multiple targets to assist with development.

```
❯ make help
help     List of available commands
up       Run the application and follow the logs
down     Stop all containers
restart  Restart all containers
logs     Print logs in stdout
version  Automatically calculate the semantic version based on the number of commits since the last change to the VERSION file
```

The typical workflow would be:

1. `make up`
2. _Make code changes in your editor of choice_
3. _The changes are automatically detected and a rebuild / rerun is initiated_
4. `Ctrl+C` to quit
5. `make restart` to restart the development environment

**Example**

```
❯ make up
docker-compose -f docker-compose.dev.yaml up -d
Creating network "go-csv_default" with the default driver
Creating app ... done
make -s logs
Attaching to app
app    |
app    |   __    _   ___
app    |  / /\  | | | |_)
app    | /_/--\ |_| |_| \_ , built with Go
app    |
app    | watching .
app    | watching app
app    | watching app/config
app    | watching app/csv
app    | watching app/faker
app    | watching app/log
app    | watching cmd
app    | watching cmd/gocsv
app    | watching container
app    | watching input
app    | watching output
app    | !exclude tmp
app    | building...
app    | go: downloading github.com/spf13/cobra v1.7.0
app    | go: downloading gopkg.in/yaml.v3 v3.0.1
app    | go: downloading github.com/TrevorEdris/banner v1.1.0
app    | go: downloading go.uber.org/zap v1.24.0
app    | go: downloading github.com/brianvoe/gofakeit/v6 v6.21.0
app    | go: downloading github.com/spf13/pflag v1.0.5
app    | go: downloading go.uber.org/multierr v1.6.0
app    | go: downloading go.uber.org/atomic v1.7.0
app    | running...
app    | 2023-05-30T14:24:02.534Z INFO Generating CSV
app    | 2023-05-30T14:24:02.534Z INFO Completed in 222.5µs
app    | app/csv/row.go has changed
app    | building...
app    | running...
app    | 2023-05-30T14:24:12.204Z INFO Generating CSV
app    | 2023-05-30T14:24:12.204Z INFO Completed in 224.8µs
```

In this example, the command specified in [.air.toml](./.air.toml) was automatically ran. Code changes
were made to `app/csv/row.go` and the binary was rebuilt and reran with the same command.

This functionality is provided by [cosmtrek/air](github.com/cosmtrek/air) and the configuration can be
found in [.air.toml](./.air.toml).

The command that is run is controlled by this section of the file:

```toml
[build]
  bin = "./tmp/main generate --input input/default.yaml --output output/default.csv"
```

**Note:** When `.air.toml` is modified, the development environment must be restarted to pick up the change.

## Contributing

Refer to the [Contribution Guidelines](./CONTRIBUTING.md) for instructions on contributing to this project.
