package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runCandidate(bin string) (string, error) {
	cmd := exec.Command(bin)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	if stderr.Len() > 0 {
		return "", fmt.Errorf("unexpected stderr output: %s", stderr.String())
	}
	return stdout.String(), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA3.go /path/to/binary")
		os.Exit(1)
	}
	out, err := runCandidate(os.Args[1])
	if err != nil {
		fmt.Printf("test failed: %v\n", err)
		os.Exit(1)
	}
	if strings.TrimSpace(out) != "" {
		fmt.Printf("test failed: expected no output, got %q\n", out)
		os.Exit(1)
	}
	fmt.Println("all tests passed")
}
