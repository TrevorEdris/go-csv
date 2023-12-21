/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"

	"github.com/TrevorEdris/go-csv/pkg/log"
	"github.com/TrevorEdris/go-csv/pkg/schema"
	"github.com/spf13/cobra"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate an existing schema file",
	Long: `Validate an existing schema file. If the file does not exist or
the schema does not meet the validation criteria, an error will be shown.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		ctx = log.ToCtx(ctx, log.FromCtx(ctx))
		schemaFile, err := cmd.Flags().GetString("schema")
		if err != nil {
			panic(err)
		}
		err = schema.ValidateSchemaFile(ctx, schemaFile)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	schemaCmd.AddCommand(validateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// validateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// validateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	validateCmd.Flags().StringP("schema", "s", "./schemas/schema.yaml", "Path to schema file to validate")
}
