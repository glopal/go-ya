package steps

import (
	"fmt"

	"github.com/glopal/go-ya/yalib"
)

func init() {
	yalib.RegisterStep("wf", NewWf)
}

type Wf struct {
	Id  string
	ech yalib.ExecContextHooks
}

func NewWf(tag string, node yalib.Node, ech yalib.ExecContextHooks) (yalib.Step, error) {
	return Wf{
		Id:  node.Node.Value,
		ech: ech,
	}, nil
}

func (wf Wf) Run(ion yalib.IoNode) (yalib.IoNode, error) {
	workflow := wf.ech.GetWorkflow(wf.Id)
	if workflow == nil {
		return nil, fmt.Errorf("could not find workflow '%s'", wf.Id)
	}
	return workflow.Run(ion)
}
