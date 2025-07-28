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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	out, err := run(bin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed:", err)
		os.Exit(1)
	}
	if out != "Yes" {
		fmt.Fprintf(os.Stderr, "expected Yes got %s\n", out)
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
