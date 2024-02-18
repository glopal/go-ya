package yalib

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
)

type ExecContextHooks interface {
	GetWorkflow(string) *Workflow
}
type ExecContext struct {
	Opts      ExecOptions
	Cmd       *cobra.Command
	Main      *Workflow
	Workflows map[string]Workflow
}

type ExecOptions struct {
	Dir          string
	OutputFormat string
}

func NewExecContext(opts ExecOptions) ExecContext {
	ec := ExecContext{
		Opts:      opts,
		Workflows: map[string]Workflow{},
	}

	yrs, err := LoadYamlResources(ec.Opts.Dir)
	if err != nil {
		panic(err)
	}

	for _, yr := range yrs {
		if yr.Type == WF {
			err := ec.AddWorkflow(yr)
			if err != nil {
				panic(err)
			}
		} else if yr.Type == MAIN {
			err := ec.AddWorkflow(yr, true)
			if err != nil {
				panic(err)
			}

			cli := Cli{}
			err = yr.Node.Decode(&cli)
			if err != nil {
				panic(err)
			}

			ec.Cmd = &cobra.Command{
				Use:   cli.Cli.Use,
				Short: cli.Cli.Short,
				RunE: func(cmd *cobra.Command, args []string) error {
					return ec.Run(StdInIo(cmd.InOrStdin()))
				},
			}
		}
	}

	return ec
}

func (ec ExecContext) Run(io IoNode) error {
	if ec.Main == nil {
		return errors.New("missing yp/main")
	}
	output, err := ec.Main.Run(io)
	if err != nil {
		return err
	}

	if ec.Opts.OutputFormat == "json" {
		output.PrettyPrintJson(os.Stdout)
	} else {
		output.PrettyPrintYaml(os.Stdout)
	}

	return nil
}

func (ec *ExecContext) AddWorkflow(yr YamlResource, isMain ...bool) error {
	wf, err := NewWorkflow(yr.Id, yr.Node, ec)
	if err != nil {
		return err
	}

	if len(isMain) > 0 && isMain[0] {
		ec.Main = &wf
	} else {
		ec.Workflows[wf.Id] = wf
	}

	return nil
}

func (ec ExecContext) GetWorkflow(name string) *Workflow {
	if wf, exists := ec.Workflows[name]; exists {
		return &wf
	}

	return nil
}
