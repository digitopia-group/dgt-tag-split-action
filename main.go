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

	if strings.Contains(refName, "_v") {
		version := "nil"
		parts := strings.Split(refName, "_v")
		if len(parts) > 1 {
			version = parts[1]
		}

		fmt.Fprintln(f, fmt.Sprintf(`tag=%s`, parts[0]))
		fmt.Fprintln(f, fmt.Sprintf(`versionnr=%s`, version))
		return
	}
	epoch, _, found := strings.Cut(refName, ":")
	if found {
		switch epoch {
		case "1", "2":
			fmt.Fprintln(f, "tag=debug")
		default:
			fmt.Fprintln(f, "tag=prod")
		}
		fmt.Fprintln(f, fmt.Sprintf("versionnr=%s", refName))
	}
}
