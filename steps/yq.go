package steps

import (
	"github.com/glopal/go-yp/yplib"
)

func init() {
	yplib.RegisterStep("yq", NewYq)
}

type Yq struct {
	Val yplib.Dval
}

func NewYq(tag string, node yplib.Node, ech yplib.ExecContextHooks) (yplib.Step, error) {
	yq := Yq{
		Val: node.ValueResolver("."),
	}
	return yq, nil
}

func (yq Yq) Run(ion yplib.IoNode) (yplib.IoNode, error) {
	yqIo, err := yq.Val(ion)
	if err != nil {
		return nil, err
	}

	// node, err := cnode.MarshalYAML()
	// if err != nil {
	// 	return nil, err
	// }

	return ion.Out(yqIo.GetCandidateNodes()), nil
}
