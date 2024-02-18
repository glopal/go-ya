package steps

import (
	"github.com/glopal/go-yp/yplib"
	"gopkg.in/yaml.v3"
)

func init() {
	yplib.RegisterStep("for", NewFor)
}

type For struct {
	Seq *yaml.Node
	Wf  yplib.Workflow
}

func NewFor(tag string, node yplib.Node, ech yplib.ExecContextHooks) (yplib.Step, error) {
	wf, err := yplib.NewWorkflow("", node.Node, ech)
	if err != nil {
		return nil, err
	}

	r := For{
		Wf: wf,
	}
	return r, nil
}

func (r For) Run(ion yplib.IoNode) (yplib.IoNode, error) {
	if ion.GetNode().Kind != yaml.SequenceNode {
		panic("expected !!seq (TODO: wrap input in seq)")
	}

	results := []*yaml.Node{}
	for _, n := range ion.GetNode().Content {
		io, err := r.Wf.Run(ion.Out(n))
		if err != nil {
			return nil, err
		}

		results = append(results, io.GetNode())
	}

	return ion.Out(results), nil
}
