package steps

import (
	"encoding/gob"

	"github.com/glopal/go-ya/yalib"
)

func init() {
	gob.Register(Context{})
	yalib.RegisterStep("context", NewContext)
}

type Context struct {
	Node yalib.Node
}

func NewContext(tag string, node yalib.Node, ech yalib.ExecContextHooks) (yalib.Step, error) {
	return Context{
		Node: node,
	}, nil
}

func (c Context) Run(ion yalib.IoNode) (yalib.IoNode, error) {
	node := c.Node.Resolve(ion)
	return ion.Out(node.Node), nil
}
