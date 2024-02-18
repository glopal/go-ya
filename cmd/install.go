package cmd

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/glopal/go-ya/yalib"
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
		ec := yalib.NewExecContext(yalib.ExecOptions{
			Dir:          dir,
			OutputFormat: outputFormat,
		})

		var b bytes.Buffer
		enc := gob.NewEncoder(io.Writer(&b))

		err := os.WriteFile("./output.txt", b.Bytes(), 0644)
		if err != nil {
			panic(err)
		}

		err = enc.Encode(ec)
		if err != nil {
			log.Fatal("encode error:", err)
		}

		fi, err := os.Open("./output.txt")
		if err != nil {
			panic(err)
		}
		dec := gob.NewDecoder(fi)

		nec := yalib.ExecContext{}
		err = dec.Decode(&nec)
		if err != nil {
			log.Fatal("decode error:", err)
		}

		err = nec.Run(yalib.StdInIo(cmd.InOrStdin()))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
