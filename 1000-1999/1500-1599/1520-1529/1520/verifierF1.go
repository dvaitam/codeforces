package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const expectedOutput = "Problem F1 is interactive and cannot be automatically solved."

func run(bin string) (string, error) {
	cmd := exec.Command(bin)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		out, err := run(bin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on repetition %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if out != expectedOutput {
			fmt.Fprintf(os.Stderr, "unexpected output on repetition %d: %q\n", i+1, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
