package main

import (
	"net/http"
	"os"
	"path"

	"github.com/jamra/gocleo"
)

func main() {
	cleo.BuildIndexes(path.Join(goSrc(), "github.com/jamra/gocleo/examples/w1_fixed.txt"), nil)
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		panic(err)
	}
}

func goSrc() string {
	return path.Join(os.Getenv("GOPATH"), "src")
}
