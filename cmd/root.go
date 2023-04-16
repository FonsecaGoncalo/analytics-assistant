package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "analytics-assistant",
	Short: "A CLI tool for interacting with a data analyst",
	Long:  `This CLI tool enables you to interact with a data analyst and connect to a MySQL database to obtain data insights using natural language.`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
