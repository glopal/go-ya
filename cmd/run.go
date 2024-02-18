package cmd

import (
	"fmt"
	"os"

	"github.com/glopal/go-yp/yplib"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			fmt.Println(arg)
		}
		ec := yplib.NewExecContext(yplib.ExecOptions{
			Dir:          dir,
			OutputFormat: outputFormat,
		})

		if ec.Cmd != nil {
			os.Args = append([]string{ec.Cmd.Use}, args...)
			err := ec.Cmd.Execute()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			os.Exit(0)
		}

		err := ec.Run(yplib.StdInIo(cmd.InOrStdin()))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
