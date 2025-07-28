package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const expectedOutput = "Problem B is interactive and cannot be automatically solved.\n"

func run(bin string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	out, err := run(bin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if out != expectedOutput {
		fmt.Printf("expected %q got %q\n", expectedOutput, out)
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
