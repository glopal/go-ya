package steps

import (
	"github.com/glopal/go-yp/yplib"
)

func init() {
	yplib.RegisterStep("context", NewContext)
}

type Context struct {
	node yplib.Node
}

func NewContext(tag string, node yplib.Node, ech yplib.ExecContextHooks) (yplib.Step, error) {
	return Context{
		node: node,
	}, nil
}

func (c Context) Run(ion yplib.IoNode) (yplib.IoNode, error) {
	node := c.node.Resolve(ion)
	return ion.Out(node.Node), nil
}
