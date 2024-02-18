package steps

import (
	"fmt"

	"github.com/glopal/go-yp/yplib"
)

func init() {
	yplib.RegisterStep("wf", NewWf)
}

type Wf struct {
	Id  string
	ech yplib.ExecContextHooks
}

func NewWf(tag string, node yplib.Node, ech yplib.ExecContextHooks) (yplib.Step, error) {
	return Wf{
		Id:  node.Node.Value,
		ech: ech,
	}, nil
}

func (wf Wf) Run(ion yplib.IoNode) (yplib.IoNode, error) {
	workflow := wf.ech.GetWorkflow(wf.Id)
	if workflow == nil {
		return nil, fmt.Errorf("could not find workflow '%s'", wf.Id)
	}
	return workflow.Run(ion)
}
