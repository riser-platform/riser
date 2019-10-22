// +build ignore

package main

import (
	"net/http"

	"github.com/shurcooL/vfsgen"
)

const assetDir = "./assets"

func main() {
	err := vfsgen.Generate(http.Dir(assetDir), vfsgen.Options{
		Filename:     "assets.go",
		PackageName:  "main",
		VariableName: "Assets",
	})
	if err != nil {
		panic(err)
	}
}
