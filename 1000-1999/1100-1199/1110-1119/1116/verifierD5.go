package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	dim = 8
)

var expected = []string{
	"X......X",
	"X......X",
	"...X.X..",
	"...X.X..",
	".X.X.X.X",
	".X.X.X.X",
	"X.X.X.X.",
	"X.X.X.X.",
}

func runProgram(bin string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseMatrix(out string) ([]string, error) {
	lines := strings.Split(strings.ReplaceAll(out, "\r\n", "\n"), "\n")
	matrix := make([]string, 0, dim)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.Contains(line, " ") || strings.Contains(line, "\t") {
			parts := strings.Fields(line)
			line = strings.Join(parts, "")
		}
		if len(line) != dim {
			return nil, fmt.Errorf("expected length %d line, got %d (%q)", dim, len(line), line)
		}
		matrix = append(matrix, line)
	}
	if len(matrix) != dim {
		return nil, fmt.Errorf("expected %d lines, got %d", dim, len(matrix))
	}
	return matrix, nil
}

func checkPattern(matrix []string) error {
	for i := 0; i < dim; i++ {
		if matrix[i] != expected[i] {
			return fmt.Errorf("line %d mismatch: expected %q got %q", i+1, expected[i], matrix[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD5.go /path/to/binary")
		os.Exit(1)
	}
	bin, err := filepath.Abs(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to resolve binary path: %v\n", err)
		os.Exit(1)
	}
	out, err := runProgram(bin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "program execution failed: %v\n", err)
		os.Exit(1)
	}
	matrix, err := parseMatrix(out)
	if err != nil {
		fmt.Fprintf(os.Stderr, "output parse error: %v\noutput:\n%s", err, out)
		os.Exit(1)
	}
	if err := checkPattern(matrix); err != nil {
		fmt.Fprintf(os.Stderr, "pattern mismatch: %v\noutput:\n%s", err, out)
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
