package main

import "./../../gocleo"

func main() {
	cleo.InitAndRun("./w1_fixed.txt", "8080", func(s1, s2 string) (result float64) {
		result = 1.0
		return
	})
}
