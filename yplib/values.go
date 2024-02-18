package yplib

import (
	"container/list"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
)

func init() {
	yqlib.InitExpressionParser()
}

type Dval func(IoNode) (IoNode, error)

func NewYqlibContext(io IoNode) yqlib.Context {
	in := io.GetCandidateNode()
	inputCandidates := list.New()
	inputCandidates.PushBack(in)

	inVarList := list.New()
	inVarList.PushBack(in)

	return yqlib.Context{
		MatchingNodes: inputCandidates,
		Variables: map[string]*list.List{
			"in": inVarList,
		},
	}
}
