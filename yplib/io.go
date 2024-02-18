package yplib

import (
	"bytes"
	"container/list"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/mattn/go-colorable"
	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"github.com/nwidger/jsoncolor"
	"gopkg.in/yaml.v3"
)

type IoNode interface {
	GetNode() *yaml.Node
	GetCandidateNode() *yqlib.CandidateNode
	GetCandidateNodes() []*yqlib.CandidateNode
	Out(any) IoNode
	Yq(*yqlib.ExpressionNode) IoNode
	PrettyPrintYaml(w *os.File)
	PrettyPrintJson(w *os.File)
	ToJson() ([]byte, error)
	Debug() string
}

type ioNode struct {
	raw    []byte
	node   *yaml.Node
	cnodes []*yqlib.CandidateNode
}

func StdInIo(r io.Reader) IoNode {
	fi, _ := os.Stdin.Stat()
	if (fi.Mode() & os.ModeCharDevice) != 0 {
		return &ioNode{
			cnodes: []*yqlib.CandidateNode{
				{
					Tag: "!!null",
				},
			},
		}
	}

	input := ioNode{}
	raw, err := io.ReadAll(r)
	if err == nil {
		input.raw = raw
	}

	return &input
}

func (ion *ioNode) GetNode() *yaml.Node {
	if ion.node == nil {
		if len(ion.cnodes) > 0 {
			node, err := ion.GetCandidateNode().MarshalYAML()
			if err == nil {
				ion.node = node
			}
		} else if ion.raw != nil {
			outNode := yaml.Node{}
			d := yaml.NewDecoder(bytes.NewReader([]byte(ion.raw)))
			err := d.Decode(&outNode)
			if err == nil {
				ion.node = &outNode
			}
		}
	}

	return ion.node
}

func (ion ioNode) GetCandidateNode() *yqlib.CandidateNode {
	if ion.cnodes != nil && len(ion.cnodes) > 0 {
		if len(ion.cnodes) == 1 {
			return ion.cnodes[0]
		} else {
			ynodes := []*yaml.Node{}

			for _, cnode := range ion.cnodes {
				yn, err := cnode.MarshalYAML()
				if err != nil {
					panic(err)
				}

				ynodes = append(ynodes, yn)
			}

			outNode := &yaml.Node{}
			err := outNode.Encode(&ynodes)
			if err == nil {
				ion.node = outNode
			} else {
				panic(err)
			}
		}
	}

	return toCandidateNode(ion.GetNode())
}

func (ion ioNode) GetCandidateNodes() []*yqlib.CandidateNode {
	if ion.cnodes != nil && len(ion.cnodes) > 0 {
		return ion.cnodes
	}

	return []*yqlib.CandidateNode{ion.GetCandidateNode()}
}

func (ion ioNode) Out(v any) IoNode {
	out := ioNode{
		cnodes: []*yqlib.CandidateNode{},
	}
	switch v := v.(type) {
	case *yaml.Node:
		out.node = v
	case []*yaml.Node:
		outNode := &yaml.Node{}

		err := outNode.Encode(&v)
		if err == nil {
			out.node = outNode
		}
	case []*yqlib.CandidateNode:
		out.cnodes = v
	case *yqlib.CandidateNode:
		out.cnodes = []*yqlib.CandidateNode{v}
	case *list.List:
		cur := v.Front()

		for {
			if cur == nil {
				break
			}
			out.cnodes = append(out.cnodes, cur.Value.(*yqlib.CandidateNode))
			cur = cur.Next()
		}
	case []byte:
		out.raw = v
	case Node:
		out.node = v.Node
	}

	return &out
}

func (ion ioNode) Yq(expressionNode *yqlib.ExpressionNode) IoNode {
	inputCandidates := list.New()
	inputCandidates.PushBack(ion.GetCandidateNode())

	ctx := yqlib.Context{
		MatchingNodes: inputCandidates,
		Variables:     map[string]*list.List{},
	}

	context, err := yqlib.NewDataTreeNavigator().GetMatchingNodes(ctx, expressionNode)
	if err != nil {
		return &ion
	}

	return ion.Out(context.MatchingNodes)

}

func (ion ioNode) PrettyPrintYaml(w *os.File) {
	// if io.raw != nil {
	// 	w.Write(io.raw)
	// 	return
	// }
	prefs := yqlib.NewDefaultYamlPreferences()
	prefs.UnwrapScalar = false
	printer := yqlib.NewPrinter(yqlib.NewYamlEncoder(4, shouldColorize(), prefs), yqlib.NewSinglePrinterWriter(w))

	// cnodes := io.GetCandidateNodes()
	// fmt.Println(len(cnodes))
	// for _, cn := range cnodes {
	// 	fmt.Println(cn)
	// }
	list, err := yqlib.NewAllAtOnceEvaluator().EvaluateNodes(".", ion.GetCandidateNodes()...)
	if err != nil {
		panic(err)
	}
	printer.PrintResults(list)
}

type Encoder interface {
	SetIndent(string, string)
	Encode(interface{}) error
}

func (ion ioNode) PrettyPrintJson(w *os.File) {
	if ion.raw != nil {
		w.Write(ion.raw)
		return
	}

	var enc Encoder = json.NewEncoder(w)
	if shouldColorize() {
		out := colorable.NewColorable(w) // needed for Windows
		enc = jsoncolor.NewEncoder(out)
	}

	enc.SetIndent("", "  ")

	err := enc.Encode(ion.GetCandidateNode())
	if err != nil {
		panic(err)
	}

}

func toCandidateNode(node *yaml.Node) *yqlib.CandidateNode {
	n := node
	if node == nil {
		return nil
	}
	if node.Kind == yaml.DocumentNode {
		n = node.Content[0]
	}

	anchorMap := map[string]*yqlib.CandidateNode{}
	o := &yqlib.CandidateNode{}

	o.UnmarshalYAML(n, anchorMap)

	return o
}

func shouldColorize() bool {
	colorsEnabled := false
	fileInfo, _ := os.Stdout.Stat()

	if (fileInfo.Mode() & os.ModeCharDevice) != 0 {
		colorsEnabled = true
	}

	return colorsEnabled
}
func (ion ioNode) ToJson() ([]byte, error) {
	var b bytes.Buffer
	enc := json.NewEncoder(io.Writer(&b))

	err := enc.Encode(ion.GetCandidateNode())
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
func (ion ioNode) Debug() string {
	baseNode := ion.GetNode()

	fmt.Println(baseNode.ShortTag())
	fmt.Println(nodeToString(baseNode, 2))

	// for _, n := range baseNode.C

	return ""
}

func nodeToString(node *yaml.Node, indent int) string {
	str := ""

	switch node.Kind {
	case yaml.DocumentNode:
		str += nodeToString(node.Content[0], indent)
	case yaml.MappingNode:
		str += mapToString(node.Content, indent)
	case yaml.ScalarNode:
		return fmt.Sprintf("%s%s %s", getIndent(indent), node.ShortTag(), node.Value)
	}

	return str
}

func getIndent(indent int) string {
	str := ""

	for range indent {
		str += " "
	}

	return str
}

func mapToString(node []*yaml.Node, indent int) string {
	str := ""

	for i, n := range node {
		if i == 0 || i%2 == 0 {
			str += nodeToString(n, indent)
		} else {
			str += fmt.Sprintf(": %s\n", nodeToString(n, 0))
		}
	}

	return str
}
