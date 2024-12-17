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
	if !found {
		fmt.Println(`Not a valid tag name.
 Use <epoch>#<major>.<minor>.<patch>.<build>-<descwithoutspaces> as tag.
 For example:
    1#0.7.3.12-testcustomerx`)
		return
	}
	if len(epoch) == 1 {
		switch epoch {
		case "1", "2":
			fmt.Fprintln(f, "tag=debug")
		default:
			fmt.Fprintln(f, "tag=prod")
		}
		refName = strings.Replace(refName, "#", ":", 1)
		fmt.Fprintf(f, "fullversion=%s\n", refName)
		fmt.Fprintf(f, "versionnr=%s:%s\n", epoch, version)
		fmt.Fprintf(f, "filenameversion=%s%%3A%s\n", epoch, version)
		return
	}
	fmt.Println("Epoch numbers should be 1 digit only: 1 -> 9")
}
