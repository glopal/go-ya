package yalib

import (
	"container/list"
	"fmt"
	"strings"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"gopkg.in/yaml.v3"
)

type Node struct {
	Node *yaml.Node
}

func (n Node) Resolver(valuePathExpression string) func(IoNode) (IoNode, error) {
	yn, err := extractNode(valuePathExpression, n.Node)
	if err != nil {
		return func(io IoNode) (IoNode, error) {
			return nil, err
		}
	}

	node := Node{yn}

	return func(io IoNode) (IoNode, error) {
		return io.Out(node.Resolve(io)), nil
	}

}
func (n Node) ValueResolver(valuePathExpression string) func(IoNode) (IoNode, error) {
	expression, err := extractStrValue(valuePathExpression, n.Node)
	if err != nil {
		panic(err)
	}

	if !strings.HasPrefix(expression, "$") && !strings.HasPrefix(expression, ".") {
		return func(io IoNode) (IoNode, error) {
			return io.Out(&yaml.Node{Value: expression}), nil
		}
	}

	expressionNode, err := yqlib.ExpressionParser.ParseExpression(expression)
	if err != nil {
		panic(err)
	}

	return func(io IoNode) (IoNode, error) {
		context, err := yqlib.NewDataTreeNavigator().GetMatchingNodes(NewYqlibContext(io), expressionNode)
		if err != nil {
			return nil, err
		}

		return io.Out(context.MatchingNodes), nil
	}
}

func extractNode(valuePathExpression string, node *yaml.Node) (*yaml.Node, error) {
	if valuePathExpression == "." {
		return node, nil
	}

	expressionNode, err := yqlib.ExpressionParser.ParseExpression(valuePathExpression)
	if err != nil {
		return nil, err
	}

	inputCandidates := list.New()
	inputCandidates.PushBack(toCandidateNode(node))

	ctx := yqlib.Context{
		MatchingNodes: inputCandidates,
		Variables:     map[string]*list.List{},
	}

	context, err := yqlib.NewDataTreeNavigator().GetMatchingNodes(ctx, expressionNode)
	if err != nil {
		return nil, err
	}

	if context.MatchingNodes.Len() == 0 {
		return nil, fmt.Errorf("failed to extract node at '%s' (zero matches)", valuePathExpression)
	}

	cnode := context.MatchingNodes.Front().Value.(*yqlib.CandidateNode)

	if cnode.Value == "null" {
		return nil, fmt.Errorf("failed to extract node at '%s' (null value)", valuePathExpression)
	}

	return cnode.MarshalYAML()
}

func extractStrValue(valuePathExpression string, node *yaml.Node) (string, error) {
	if valuePathExpression == "." {
		return node.Value, nil
	}

	expressionNode, err := yqlib.ExpressionParser.ParseExpression(valuePathExpression)
	if err != nil {
		return "", err
	}

	inputCandidates := list.New()
	inputCandidates.PushBack(toCandidateNode(node))

	ctx := yqlib.Context{
		MatchingNodes: inputCandidates,
		Variables:     map[string]*list.List{},
	}

	context, err := yqlib.NewDataTreeNavigator().GetMatchingNodes(ctx, expressionNode)
	if err != nil {
		return "", err
	}

	if context.MatchingNodes.Len() == 0 {
		return "", fmt.Errorf("failed to extract value at '%s'", valuePathExpression)
	}

	cnode := context.MatchingNodes.Front().Value.(*yqlib.CandidateNode)
	if cnode.Tag != "!!str" {
		return "", fmt.Errorf("value at '%s' is not !!str", valuePathExpression)
	}

	return cnode.Value, nil
}

func (n Node) At(valuePathExpression string) Node {
	if valuePathExpression == "." {
		return n
	}

	return n
}

func (n Node) Resolve(io IoNode) Node {
	nn := deepCopyNode(n.Node, nil)

	yqNodes := getScalarNodesByTag(nn, "!yq")

	if len(yqNodes) == 0 {
		return Node{nn}
	}

	ctx := NewYqlibContext(io)

	for _, yn := range yqNodes {
		expressionNode, err := yqlib.ExpressionParser.ParseExpression(yn.Value)
		if err != nil {
			panic(err)
		}

		context, err := yqlib.NewDataTreeNavigator().GetMatchingNodes(ctx, expressionNode)
		if err != nil {
			panic(err)
		}
		if context.MatchingNodes.Len() == 0 {
			panic("no match")
		}

		node, err := context.MatchingNodes.Front().Value.(*yqlib.CandidateNode).MarshalYAML()
		if err != nil {
			panic(err)
		}

		*yn = *node
	}

	return Node{nn}
}

func getScalarNodesByTag(node *yaml.Node, tag string) []*yaml.Node {
	if node.Kind == yaml.ScalarNode && node.Tag == tag {
		return []*yaml.Node{node}
	}

	scalarNodes := []*yaml.Node{}

	if node.Kind <= yaml.MappingNode {
		for _, n := range node.Content {
			scalarNodes = append(scalarNodes, getScalarNodesByTag(n, tag)...)
		}
	}

	return scalarNodes
}

func deepCopyNode(node *yaml.Node, cache map[*yaml.Node]*yaml.Node) *yaml.Node {
	if n, ok := cache[node]; ok {
		return n
	}
	if cache == nil {
		cache = make(map[*yaml.Node]*yaml.Node)
	}
	copy := *node
	cache[node] = &copy
	copy.Content = nil
	for _, elem := range node.Content {
		copy.Content = append(copy.Content, deepCopyNode(elem, cache))
	}
	if node.Alias != nil {
		copy.Alias = deepCopyNode(node.Alias, cache)
	}
	return &copy
}
