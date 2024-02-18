package yalib

import (
	"gopkg.in/yaml.v3"
)

type Step interface {
	Run(IoNode) (IoNode, error)
}

type Workflow struct {
	Id    string
	Node  *yaml.Node
	Steps []Step
}

func NewWorkflow(id string, node *yaml.Node, ec ExecContextHooks) (Workflow, error) {
	wf := Workflow{
		Id:    id,
		Steps: []Step{},
	}

	for i, c := range node.Content {
		if c.Value == "steps" {
			for _, s := range node.Content[i+1].Content {
				step, err := NewStep(s.Content[0].Value, s.Content[1].Tag, s.Content[1], ec)
				if err != nil {
					return Workflow{}, err
				}

				wf.Steps = append(wf.Steps, step)
			}
		}
	}

	return wf, nil
}

func (wf Workflow) Run(in IoNode) (IoNode, error) {
	for _, step := range wf.Steps {
		out, err := step.Run(in)
		if err != nil {
			return nil, err
		}

		in = out
	}

	return in, nil
}

type workflowSchema struct {
	Steps []*yaml.Node `yaml:"steps"`
}
