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

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new schema file",
	Long: `Create a new schema file. The resulting file will be filled out with
a basic example of what a schema file may contain.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		ctx = log.ToCtx(ctx, log.FromCtx(ctx))
		outputFile, err := cmd.Flags().GetString("output")
		if err != nil {
			panic(err)
		}
		err = schema.CreateExample(ctx, outputFile)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	schemaCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	newCmd.Flags().StringP("output", "o", "./schemas/schema.example.yaml", "Output file for new schema")
}
