package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]

	cmd := exec.Command(binary)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		fmt.Printf("runtime error: %v\nstderr: %s\n", err, errBuf.String())
		os.Exit(1)
	}
	out := strings.TrimSpace(outBuf.String())
	expected := "Problem F is interactive and cannot be automatically solved."
	if out != expected {
		fmt.Printf("expected %q got %q\n", expected, out)
		os.Exit(1)
	}
	fmt.Println("All 1 tests passed")
}
