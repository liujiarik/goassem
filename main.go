package main

import (
	"os"
	"fmt"

	"github.com/liujiarik/goassem/run"
	"bytes"
	"io"
	"strings"
)

func main() {

	allArgs := os.Args

	if allArgs[len(allArgs)-1] == "-" {
		stdin := &bytes.Buffer{}
		if _, err := io.Copy(stdin, os.Stdin); err == nil {
			stdinArgs := strings.Fields(stdin.String())
			allArgs = append(allArgs[:len(allArgs)-1], stdinArgs...)
		}
	}

	t, err := run.Run(os.Stdout, allArgs)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Happen error | %+v\n", t)
		os.Exit(2)
	}
	//fmt.Fprintf(os.Stdout, t)
	fmt.Println(t)

}
