package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const ref1970C3 = "1000-1999/1900-1999/1970-1979/1970/1970C3.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC3.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to read input:", err)
		os.Exit(1)
	}
	expectedLines, err := countRounds(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	refBin, cleanup, err := buildReference(ref1970C3)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer cleanup()

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n", err)
		os.Exit(1)
	}
	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}

	refLines := nonEmptyLines(refOut)
	candLines := nonEmptyLines(candOut)
	if len(refLines) != expectedLines {
		fmt.Fprintf(os.Stderr, "reference output has %d lines, expected %d\n", len(refLines), expectedLines)
		os.Exit(1)
	}
	if len(candLines) != expectedLines {
		fmt.Fprintf(os.Stderr, "candidate output has %d lines, expected %d\n", len(candLines), expectedLines)
		os.Exit(1)
	}
	for i := 0; i < expectedLines; i++ {
		if !strings.EqualFold(refLines[i], candLines[i]) {
			fmt.Fprintf(os.Stderr, "line %d mismatch: expected %s, got %s\n", i+1, refLines[i], candLines[i])
			fmt.Fprintln(os.Stderr, "reference output:")
			fmt.Fprintln(os.Stderr, refOut)
			fmt.Fprintln(os.Stderr, "candidate output:")
			fmt.Fprintln(os.Stderr, candOut)
			os.Exit(1)
		}
	}

	fmt.Println("Accepted")
}

func countRounds(data []byte) (int, error) {
	reader := bufio.NewReader(bytes.NewReader(data))
	var n, t int
	if _, err := fmt.Fscan(reader, &n, &t); err != nil {
		return 0, fmt.Errorf("failed to read n and t: %v", err)
	}
	for i := 0; i < n-1; i++ {
		var u, v int
		if _, err := fmt.Fscan(reader, &u, &v); err != nil {
			return 0, fmt.Errorf("failed to read edge %d: %v", i+1, err)
		}
	}
	for i := 0; i < t; i++ {
		var x int
		if _, err := fmt.Fscan(reader, &x); err != nil {
			return 0, fmt.Errorf("failed to read starting node %d: %v", i+1, err)
		}
	}
	return t, nil
}

func nonEmptyLines(out string) []string {
	scanner := bufio.NewScanner(strings.NewReader(out))
	lines := make([]string, 0)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines
}

func buildReference(src string) (string, func(), error) {
	dir, err := os.MkdirTemp("", "verifier-1970C3-")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(dir, "ref.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("go build failed: %v\n%s", err, stderr.String())
	}
	return bin, func() { os.RemoveAll(dir) }, nil
}

func runProgram(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	return out.String(), cmd.Run()
}
