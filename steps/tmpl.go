package steps

import (
	"bytes"
	"html/template"

	"github.com/glopal/go-ya/yalib"
)

func init() {
	yalib.RegisterStep("tmpl", NewTmpl)
}

type Tmpl struct {
	tmpl *template.Template
}

func NewTmpl(tag string, node yalib.Node, ech yalib.ExecContextHooks) (yalib.Step, error) {
	tmpl, err := template.New("test").Parse(node.Node.Value)
	if err != nil {
		return nil, err
	}

	return Tmpl{
		tmpl: tmpl,
	}, nil
}

func (t Tmpl) Run(ion yalib.IoNode) (yalib.IoNode, error) {
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
