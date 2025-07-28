package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cmd := exec.Command(bin)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error executing binary: %v\n", err)
		os.Exit(1)
	}
	expected := "Problem E is interactive and cannot be solved automatically."
	output := strings.TrimSpace(string(out))
	if output != expected {
		fmt.Fprintf(os.Stderr, "expected %q got %q\n", expected, output)
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
