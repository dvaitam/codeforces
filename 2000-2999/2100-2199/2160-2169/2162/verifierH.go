package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

const refSource = "2162H.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	oraclePath, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build oracle: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)
	var inputBuilder strings.Builder
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
		inputBuilder.WriteString(line)
		inputBuilder.WriteByte('\n')
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read input: %v\n", err)
		os.Exit(1)
	}

	if len(lines) == 0 {
		fmt.Fprintln(os.Stderr, "empty input")
		os.Exit(1)
	}

	input := inputBuilder.String()
	expOut, err := runProgram(oraclePath, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle runtime error: %v\n", err)
		os.Exit(1)
	}
	gotOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}

	expVals := parseInts(strings.TrimSpace(expOut))
	gotVals := parseInts(strings.TrimSpace(gotOut))

	if len(expVals) != len(gotVals) {
		fmt.Fprintf(os.Stderr, "length mismatch: expected %d values, got %d\n", len(expVals), len(gotVals))
		os.Exit(1)
	}

	for i := range expVals {
		if expVals[i] != gotVals[i] {
			fmt.Fprintf(os.Stderr, "mismatch at position %d: expected %d, got %d\n", i+1, expVals[i], gotVals[i])
			os.Exit(1)
		}
	}

	fmt.Println("All tests passed")
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("failed to determine path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2162H-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracle2162H")
	cmd := exec.Command("go", "build", "-o", outPath, refSource)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("build oracle failed: %v\n%s", err, string(out))
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

func parseInts(out string) []int64 {
	if strings.TrimSpace(out) == "" {
		return nil
	}
	fields := strings.Fields(out)
	result := make([]int64, len(fields))
	for i, f := range fields {
		v, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil
		}
		result[i] = v
	}
	return result
}
