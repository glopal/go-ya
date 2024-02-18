package yalib

import (
	"strings"

	"gopkg.in/yaml.v3"
)

type DocType int

const (
	UNKNOWN DocType = iota
	WF
	MAIN
)

type YamlResource struct {
	Id   string
	Type DocType
	Node *yaml.Node
	Path string
}

func NewYamlResource(path string, node *yaml.Node) (YamlResource, error) {
	id, docType := parseIdAndType(node)
	yr := YamlResource{
		Id:   id,
		Type: docType,
		Node: node.Content[0],
		Path: path,
	}

	return yr, nil
}

func parseIdAndType(node *yaml.Node) (string, DocType) {
	docTagTokens := strings.Split(node.Content[0].ShortTag(), "/")
	if len(docTagTokens) != 2 {
		return "", UNKNOWN
	}

	switch docTagTokens[0] {
	case "!yp":
		if docTagTokens[1] == "main" {
			return "main", MAIN
		}
	case "!wf":
		return docTagTokens[1], WF
	}

	return "", UNKNOWN
}

// type KindMap struct {
// 	Kind string `yaml:"kind"`
// }

// func determineKind(node *yaml.Node) string {
// 	docTagTokens := strings.Split(node.Content[0].ShortTag(), "/")
// 	if len(docTagTokens) == 2 && docTagTokens[0] == "!kind" {
// 		return docTagTokens[1]
// 	}

// 	if node.Content[0].Kind == yaml.MappingNode {
// 		kindMap := &KindMap{}

// 		err := node.Decode(kindMap)
// 		if err != nil {
// 			return ""
// 		}

// 		return kindMap.Kind
// 	}

// 	return ""
// }
