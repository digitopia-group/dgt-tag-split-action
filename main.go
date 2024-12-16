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

		fmt.Fprintf(f, `tag=%s\n`, parts[0])
		fmt.Fprintf(f, `versionnr=%s\n`, version)
		fmt.Fprintf(f, `fullversion=%s\n`, refName)
		return
	}
	epoch, version, found := strings.Cut(refName, "#")
	if found && len(epoch) == 1 {
		switch epoch {
		case "1", "2":
			fmt.Fprintln(f, "tag=debug")
		default:
			fmt.Fprintln(f, "tag=prod")
		}
		refName = strings.Replace(refName, "#", ":", 1)
		fmt.Fprintf(f, "fullversion=%s\n", refName)
		fmt.Fprintf(f, "versionnr=%s:%s\n", epoch, version)
	}
}
