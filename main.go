package main

import (
	"fmt"
	"os"
	"strings"
)

/*
- name: Save state
run: echo "{name}={value}" >> $GITHUB_STATE

- name: Set output
run: echo "{name}={value}" >> $GITHUB_OUTPUT
*/

func main() {
	refName := os.Getenv("INPUT_REFNAME")
	outputFile := os.Getenv("GITHUB_OUTPUT")

	f, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	version := "nil"
	parts := strings.Split(refName, "_v")
	if len(parts) > 1 {
		version = parts[1]
	}

	fmt.Fprintf(f, `versionnr=%s`, version)
}
