package steps

import (
	"bytes"
	"html/template"

	"github.com/glopal/go-yp/yplib"
)

func init() {
	yplib.RegisterStep("tmpl", NewTmpl)
}

type Tmpl struct {
	tmpl *template.Template
}

func NewTmpl(tag string, node yplib.Node, ech yplib.ExecContextHooks) (yplib.Step, error) {
	tmpl, err := template.New("test").Parse(node.Node.Value)
	if err != nil {
		return nil, err
	}

	return Tmpl{
		tmpl: tmpl,
	}, nil
}

func (t Tmpl) Run(ion yplib.IoNode) (yplib.IoNode, error) {
	var ctx interface{}
	err := ion.GetNode().Decode(&ctx)
	if err != nil {
		panic(err)
	}

	var d bytes.Buffer

	err = t.tmpl.Execute(&d, ctx)
	if err != nil {
		panic(err)
	}

	return ion.Out(d.Bytes()), nil
}
