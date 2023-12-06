package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	refName := os.Getenv("INPUT_REFNAME")

	parts := strings.Split(refName, "_v")
	if len(parts) < 2 {
		fmt.Println(`::set-output name=myOutput::nil`)
	}

	fmt.Println(fmt.Sprintf(`::set-output name=myOutput::%s`, parts[1]))
}
