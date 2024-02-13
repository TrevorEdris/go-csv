package main

import (
	"flag"
	"log"
	"os"

	"github.com/TrevorEdris/go-csv/pkg/generator"
)

func main() {

	help := flag.Bool("help", false, "Show usage")
	input := flag.String("input", "input/example.yaml", "Input filename")
	output := flag.String("output", "out.csv", "Output filename")
	count := flag.Int("count", generator.DefaultCount, "Number of records to generate")
	delimiter := flag.String("delimiter", generator.DefaultDelimiter, "Character to use as the CSV delimiter")

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	g, err := generator.NewGenerator(*input, *output, *delimiter)
	if err != nil {
		log.Fatalf("Failed to initialize generator: %s", err.Error())
	}

	err = g.Generate(*count)
	if err != nil {
		log.Fatalf("Failed to generate CSV: %s", err.Error())
	}
}
