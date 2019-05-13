// +build ignore

package main

import (
	"log"

	"examples/templates"

	"github.com/shurcooL/vfsgen"
)

func main() {
	err := vfsgen.Generate(templates.Templates, vfsgen.Options{
		PackageName:  "templates",
		BuildTags:    "!dev",
		VariableName: "Templates",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
