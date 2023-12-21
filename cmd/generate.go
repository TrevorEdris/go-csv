/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/TrevorEdris/go-csv/pkg/generator"
	"github.com/TrevorEdris/go-csv/pkg/log"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a CSV file based on an existing schema",
	Long:  `Generate a CSV file based on an existing schema`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := log.ToCtx(context.Background(), log.FromCtx(context.Background()))
		schemaFile, err := cmd.Flags().GetString("schema")
		if err != nil {
			panic(err)
		}
		outputFile, err := cmd.Flags().GetString("output")
		if err != nil {
			panic(err)
		}
		if outputFile == "" {
			// Cut the extension off the provided schemaFile
			ext := filepath.Ext(schemaFile)
			filename := filepath.Base(schemaFile)
			filename, _ = strings.CutSuffix(filename, ext)

			outputDir := filepath.Join(".", "output")
			err = os.MkdirAll(outputDir, os.ModePerm)
			if err != nil {
				panic(err)
			}
			outputFile = filepath.Join(outputDir, fmt.Sprintf("%s.csv", filename))
		}

		g, err := generator.NewGenerator(generator.GeneratorConfig{})
		if err != nil {
			panic(err)
		}
		err = g.GenerateFromFile(ctx, schemaFile, outputFile)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	generateCmd.Flags().StringP("schema", "s", "./schemas/schema.yaml", "Path to schema file")
	generateCmd.Flags().StringP("output", "o", "", "Path to output file. If unspecified, the filename will match the schema file with the '.csv' extension, such as 'schema.csv'")
}
