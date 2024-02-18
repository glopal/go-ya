package yalib

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func LoadYamlResources(path string) ([]YamlResource, error) {
	yamlResources := []YamlResource{}

	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if isYamlFile(path) {
		yrs, err := loadYamlFile(path)
		if err != nil {
			return yamlResources, err
		}

		yamlResources = append(yamlResources, yrs...)
	} else if info.IsDir() {
		err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				return nil
			}

			if isYamlFile(path) {
				yrs, err := loadYamlFile(path)
				if err != nil {
					return err
				}

				yamlResources = append(yamlResources, yrs...)
			}

			return nil
		})
		if err != nil {
			return yamlResources, err
		}
	}

	return yamlResources, nil

}

func loadYamlFile(path string) ([]YamlResource, error) {
	yrs := []YamlResource{}
	nodes, err := loadYamlNodes(path)
	if err != nil {
		return yrs, err
	}

	for _, node := range nodes {
		yr, err := NewYamlResource(path, node)
		if err != nil {
			return yrs, err
		}

		yrs = append(yrs, yr)
	}

	return yrs, nil
}

func loadYamlNodes(file string) ([]*yaml.Node, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	nodes := []*yaml.Node{}
	d := yaml.NewDecoder(f)
	for {
		// create new spec here
		node := yaml.Node{}
		// pass a reference to spec reference
		err := d.Decode(&node)
		// check it was parsed
		if err == nil {
			nodes = append(nodes, &node)
			continue
		}
		// break the loop in case of EOF
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}

	}

	return nodes, nil
}

func isYamlFile(file string) bool {
	ext := filepath.Ext(file)
	return ext == ".yml" || ext == ".yaml"
}
