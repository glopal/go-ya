package steps

import (
	"encoding/gob"

	"github.com/glopal/go-ya/yalib"
)

func init() {
	gob.Register(Yq{})
	yalib.RegisterStep("yq", NewYq)
}

type Yq struct {
	Val yalib.Dval
}

func NewYq(tag string, node yalib.Node, ech yalib.ExecContextHooks) (yalib.Step, error) {
	yq := Yq{
		Val: node.ValueResolver("."),
	}
	return yq, nil
}

func (yq Yq) Run(ion yalib.IoNode) (yalib.IoNode, error) {
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
