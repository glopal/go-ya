package cmd

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"

	"github.com/glopal/go-yp/yplib"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var installCmd = &cobra.Command{
	Use:   "install",
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

		fi, err := os.OpenFile("./output.txt", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			panic(err)
		}
		enc := gob.NewEncoder(fi)

		err = enc.Encode(ec)
		if err != nil {
			log.Fatal("encode error:", err)
		}

		fi, err = os.Open("./output.txt")
		if err != nil {
			panic(err)
		}
		dec := gob.NewDecoder(fi)

		nec := yplib.ExecContext{}
		err = dec.Decode(&nec)
		if err != nil {
			log.Fatal("decode error:", err)
		}

		err = nec.Run(yplib.StdInIo(cmd.InOrStdin()))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
