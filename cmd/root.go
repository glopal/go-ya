package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var dir string
var outputFormat string

var rootCmd = &cobra.Command{
	Use:   "yp",
	Short: "Yaml Pipeline",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// TODO convert the dir flag to just be the first arg
	rootCmd.PersistentFlags().StringVarP(&dir, "dir", "d", ".", "The directory that contains main.yml")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output-format", "o", "yaml", "The output format hint. Defaults to 'yaml'. Valid values: yaml,json")
}
