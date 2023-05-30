package main

import (
	"fmt"
	"os"

	"github.com/TrevorEdris/go-csv/app"
	"github.com/TrevorEdris/go-csv/app/config"
	"github.com/spf13/cobra"
)

var Version string

func main() {
	rootCmd := &cobra.Command{
		Use:   "gocsv",
		Short: "Interact with CSV files",
	}

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version of gocsv",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(Version)
		},
	}
	rootCmd.AddCommand(versionCmd)

	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate CSV file with randomized data",
		Run: func(cmd *cobra.Command, args []string) {
			logLevel, err := cmd.Flags().GetString("logLevel")
			printAndExit(err)
			input, err := cmd.Flags().GetString("input")
			printAndExit(err)
			output, err := cmd.Flags().GetString("output")
			printAndExit(err)

			generate(logLevel, input, output)
		},
	}
	generateCmd.Flags().String("logLevel", "INFO", "Log level (DEBUG|INFO|WARN)")
	generateCmd.Flags().String("input", "input/default.yaml", "Relative filename to the configuration yaml file")
	generateCmd.Flags().String("output", "output/default.csv", "Filename to write output to")
	err := generateCmd.MarkFlagRequired("input")
	printAndExit(err)

	rootCmd.AddCommand(generateCmd)

	err = rootCmd.Execute()
	printAndExit(err)
}

func generate(logLevel, inputFilename, outputFilename string) {
	a, err := app.New(&config.Runtime{
		LogLevel: logLevel,
		Input:    inputFilename,
		Output:   outputFilename,
	})
	printAndExit(err)

	printAndExit(a.Generate())
}

func printAndExit(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
