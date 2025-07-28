package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const expectedOutput = "Problem C is interactive and cannot be automatically solved."

func runCase(bin string) error {
	cmd := exec.Command(bin)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expectedOutput {
		return fmt.Errorf("expected %q got %q", expectedOutput, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		if err := runCase(bin); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
