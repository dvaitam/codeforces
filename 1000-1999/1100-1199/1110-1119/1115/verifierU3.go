package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func runProgram(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func normalizeMatrix(out string, size int) ([]string, error) {
	raw := strings.Split(strings.ReplaceAll(out, "\r\n", "\n"), "\n")
	matrix := make([]string, 0, size)
	for _, line := range raw {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		line = strings.ReplaceAll(line, " ", "")
		line = strings.ReplaceAll(line, "\t", "")
		if len(line) != size {
			return nil, fmt.Errorf("expected line with %d characters, got %d (%q)", size, len(line), line)
		}
		matrix = append(matrix, line)
	}
	if len(matrix) != size {
		return nil, fmt.Errorf("expected %d lines, got %d", size, len(matrix))
	}
	return matrix, nil
}

func expectedMatrix(n int) []string {
	size := 1 << n
	half := size >> 1
	matrix := make([]string, size)
	for i := 0; i < size; i++ {
		var row strings.Builder
		row.Grow(size)
		for j := 0; j < size; j++ {
			var ch byte
			switch {
			case i < half && j < half:
				if i+j == half-1 {
					ch = 'X'
				} else {
					ch = '.'
				}
			case i < half && j >= half:
				ch = '.'
			case i >= half && j < half:
				ch = '.'
			default:
				ch = 'X'
			}
			row.WriteByte(ch)
		}
		matrix[i] = row.String()
	}
	return matrix
}

func checkCase(bin string, n int) error {
	input := fmt.Sprintf("%d\n", n)
	out, err := runProgram(bin, input)
	if err != nil {
		return err
	}
	size := 1 << n
	got, err := normalizeMatrix(out, size)
	if err != nil {
		return err
	}
	want := expectedMatrix(n)
	for i := 0; i < size; i++ {
		if got[i] != want[i] {
			return fmt.Errorf("line %d mismatch: expected %q got %q", i+1, want[i], got[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierU3.go /path/to/binary")
		os.Exit(1)
	}
	bin, err := filepath.Abs(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to resolve binary path: %v\n", err)
		os.Exit(1)
	}
	total := 0
	for n := 2; n <= 5; n++ {
		if err := checkCase(bin, n); err != nil {
			fmt.Fprintf(os.Stderr, "case N=%d failed: %v\n", n, err)
			os.Exit(1)
		}
		total++
	}
	fmt.Printf("All %d tests passed\n", total)
}
