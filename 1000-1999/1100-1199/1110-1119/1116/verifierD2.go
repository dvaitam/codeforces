package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	refSource        = "1116D2.go"
	tempOraclePrefix = "oracle-1116D2-"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD2.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build oracle: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	nValues := []int{2, 3, 4, 5}
	for idx, n := range nValues {
		input := fmt.Sprintf("%d\n", n)
		exp, err := runProgram(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle runtime error on test %d (n=%d): %v\n", idx+1, n, err)
			os.Exit(1)
		}
		got, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (n=%d): %v\n", idx+1, n, err)
			os.Exit(1)
		}
		if err := comparePatterns(exp, got, n); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (n=%d) failed: %v\n", idx+1, n, err)
			fmt.Println("Expected pattern:")
			fmt.Print(exp)
			fmt.Println("Candidate output:")
			fmt.Print(got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(nValues))
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("failed to determine working directory")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", tempOraclePrefix)
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleD2")
	cmd := exec.Command("go", "build", "-o", outPath, refSource)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("build failed: %v\n%s", err, string(out))
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
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

func comparePatterns(expected, got string, n int) error {
	expLines := filterLines(expected)
	gotLines := filterLines(got)
	size := 1 << n
	if len(expLines) != size {
		return fmt.Errorf("oracle produced %d lines, expected %d", len(expLines), size)
	}
	if len(gotLines) != size {
		return fmt.Errorf("candidate produced %d lines, expected %d", len(gotLines), size)
	}
	for i := 0; i < size; i++ {
		if len(expLines[i]) != size {
			return fmt.Errorf("oracle line %d length %d, expected %d", i+1, len(expLines[i]), size)
		}
		if len(gotLines[i]) != size {
			return fmt.Errorf("candidate line %d length %d, expected %d", i+1, len(gotLines[i]), size)
		}
		for j := 0; j < size; j++ {
			expChar := normalizeChar(expLines[i][j])
			gotChar := normalizeChar(gotLines[i][j])
			if expChar != gotChar {
				return fmt.Errorf("mismatch at (%d,%d): expected %c got %c", i+1, j+1, expChar, gotChar)
			}
		}
	}
	return nil
}

func filterLines(s string) []string {
	lines := strings.Split(s, "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			out = append(out, line)
		}
	}
	return out
}

func normalizeChar(c byte) byte {
	if c == 'X' || c == 'x' {
		return 'X'
	}
	return '.'
}
