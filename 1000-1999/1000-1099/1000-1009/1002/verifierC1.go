package main

import (
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

func parseOutput(output string) (string, []string, error) {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) == 0 || lines[0] == "" {
		return "", nil, fmt.Errorf("empty output")
	}
	header := lines[0]
	ops := lines[1:]
	return header, ops, nil
}

func validateHeader(header string, expectedCount int) error {
	var count int
	if _, err := fmt.Sscanf(header, "%d", &count); err != nil {
		return fmt.Errorf("failed to parse operation count: %v", err)
	}
	if count != expectedCount {
		return fmt.Errorf("expected %d operations, got %d", expectedCount, count)
	}
	return nil
}

func validateOperations(ops []string) error {
	if len(ops) != 2 {
		return fmt.Errorf("expected 2 operations, got %d", len(ops))
	}
	if strings.TrimSpace(ops[0]) != "RY 1 0.25" {
		return fmt.Errorf("unexpected first operation: %s", ops[0])
	}
	if strings.TrimSpace(ops[1]) != "MEASURE 1" {
		return fmt.Errorf("unexpected second operation: %s", ops[1])
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	out, err := runCandidate(bin)
	if err != nil {
		fmt.Printf("test failed: %v\n", err)
		os.Exit(1)
	}
	header, ops, err := parseOutput(out)
	if err != nil {
		fmt.Printf("invalid output: %v\noutput:\n%s\n", err, out)
		os.Exit(1)
	}
	if err := validateHeader(header, 2); err != nil {
		fmt.Printf("invalid header: %v\noutput:\n%s\n", err, out)
		os.Exit(1)
	}
	if err := validateOperations(ops); err != nil {
		fmt.Printf("invalid operations: %v\noutput:\n%s\n", err, out)
		os.Exit(1)
	}
	fmt.Println("all tests passed")
}
