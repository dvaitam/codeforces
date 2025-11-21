package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func runProgram(path string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return out.String(), nil
}

func parseOutput(out string, q int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) < q {
		return nil, fmt.Errorf("expected %d outputs, got %d", q, len(fields))
	}
	ans := make([]int64, q)
	for i := 0; i < q; i++ {
		if _, err := fmt.Sscan(fields[i], &ans[i]); err != nil {
			return nil, fmt.Errorf("failed to parse integer at query %d: %v", i+1, err)
		}
	}
	return ans, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierK.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[len(os.Args)-1]
	if target == "--" {
		fmt.Println("usage: go run verifierK.go /path/to/binary")
		os.Exit(1)
	}

	inputData, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read input: %v\n", err)
		os.Exit(1)
	}
	reader := strings.NewReader(string(inputData))
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read n: %v\n", err)
		os.Exit(1)
	}
	for i := 0; i < n-1; i++ {
		var u, v int
		var w int64
		fmt.Fscan(reader, &u, &v, &w)
	}
	var q int
	if _, err := fmt.Fscan(reader, &q); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read q: %v\n", err)
		os.Exit(1)
	}

	_, src, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(src)
	refPath := filepath.Join(baseDir, "1866K.go")

	refOut, err := runProgram(refPath, inputData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n", err)
		os.Exit(1)
	}
	expected, err := parseOutput(refOut, q)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference output parse error: %v\n", err)
		os.Exit(1)
	}

	out, err := runProgram(target, inputData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target runtime error: %v\n", err)
		os.Exit(1)
	}
	ans, err := parseOutput(out, q)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target output parse error: %v\n", err)
		os.Exit(1)
	}

	for i := 0; i < q; i++ {
		if ans[i] != expected[i] {
			fmt.Fprintf(os.Stderr, "wrong answer on query %d: expected %d got %d\n", i+1, expected[i], ans[i])
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
