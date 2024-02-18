package main

import (
	"os"

	"github.com/glopal/go-yp/cmd"
	_ "github.com/glopal/go-yp/steps"
	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"gopkg.in/op/go-logging.v1"
)

func init() {
	// disable yqlib debug logging
	leveled := logging.AddModuleLevel(logging.NewLogBackend(os.Stderr, "", 0))
	leveled.SetLevel(logging.ERROR, "")
	yqlib.GetLogger().SetBackend(leveled)
}

func main() {
	// 	ys := `
	// Origin: "104.246.138.141"
	// Url: "https://httpbin.org/get"
	// Host: "httpbin.org"
	// `
	// 	// f, err := os.Open("./test.yml")
	// 	// if err != nil {
	// 	// 	panic(err)
	// 	// }
	// 	d := yaml.NewDecoder(bytes.NewReader([]byte(ys)))

	// 	node := yaml.Node{}
	// 	err := d.Decode(&node)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	fmt.Println(node.Content[0].Content[0].Value)

	cmd.Execute()
}
