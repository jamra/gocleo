package main

import (
	"github.com/jamra/gocleo"
	"net/http"
)

func main() {
	cleo.BuildIndexes("./w1_fixed.txt", nil)
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		panic(err)
	}
}
