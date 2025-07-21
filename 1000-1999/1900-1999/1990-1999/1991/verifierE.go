package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func runCase(exe string) error {
	cmd := exec.Command(exe)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if out.Len() != 0 {
		return fmt.Errorf("expected no output, got: %s", out.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	for i := 0; i < 100; i++ {
		if err := runCase(exe); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
