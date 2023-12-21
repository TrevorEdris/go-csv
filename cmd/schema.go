/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// schemaCmd represents the schema command
var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Schema command; please use a sub-command",
	Long:  `Schema command; please use a sub-command`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("schema called; please use a sub-command")
	},
}

func init() {
	rootCmd.AddCommand(schemaCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// schemaCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// schemaCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
