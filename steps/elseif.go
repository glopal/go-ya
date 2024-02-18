package steps

import (
	"strconv"

	"github.com/glopal/go-ya/yalib"
	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"gopkg.in/yaml.v3"
)

func init() {
	yalib.RegisterStep("elseif", NewElseIf)
}

type ElseIf struct {
	Conditions map[*yqlib.ExpressionNode]yalib.Workflow
}

func NewElseIf(tag string, node yalib.Node, ech yalib.ExecContextHooks) (yalib.Step, error) {
	elseIf := ElseIf{
		Conditions: map[*yqlib.ExpressionNode]yalib.Workflow{},
	}

	rawConditions := map[string]yaml.Node{}
	err := node.Node.Decode(&rawConditions)
	if err != nil {
		return nil, err
	}

	for expr, n := range rawConditions {
		en, err := yqlib.ExpressionParser.ParseExpression(expr)
		if err != nil {
			return nil, err
		}

		wf, err := yalib.NewWorkflow("", &n, ech)
		if err != nil {
			return nil, err
		}

		elseIf.Conditions[en] = wf
	}

	return elseIf, nil
}

func (ei ElseIf) Run(ion yalib.IoNode) (yalib.IoNode, error) {
	for en, wf := range ei.Conditions {
		cn := ion.Yq(en).GetCandidateNodes()
		if len(cn) > 0 {
			boolValue, err := strconv.ParseBool(cn[0].Value)
			if err != nil {
				continue
			}

			if boolValue {
				return wf.Run(ion)
			}
		}
	}

	return ion, nil
}
