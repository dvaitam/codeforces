package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func run(bin string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for i := 1; i <= 100; i++ {
		got, err := run(bin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i, err)
			os.Exit(1)
		}
		if got != "Hello, World!" {
			fmt.Fprintf(os.Stderr, "case %d failed: expected Hello, World! got %s\n", i, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
