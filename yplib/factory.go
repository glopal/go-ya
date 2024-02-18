package yplib

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// TODO add 'address' field (ie. wf/foo/steps[1])
type NewStepFunc func(string, Node, ExecContextHooks) (Step, error)

var stepConstructors map[string]NewStepFunc = map[string]NewStepFunc{}

func RegisterStep(name string, constructor NewStepFunc) {
	stepConstructors[name] = constructor
}

func NewStep(name, tag string, contents *yaml.Node, ech ExecContextHooks) (Step, error) {
	if constructor, exists := stepConstructors[name]; exists {
		return constructor(tag, Node{contents}, ech)
	}

	return nil, fmt.Errorf("undefined step '%s'", name)
}
