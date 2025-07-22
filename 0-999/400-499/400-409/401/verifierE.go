package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runCandidate(bin string) (string, error) {
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
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		out, err := runCandidate(bin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if out != "" {
			fmt.Fprintf(os.Stderr, "case %d failed: expected empty output got %s\n", i+1, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
