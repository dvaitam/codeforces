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

func parseOutput(out string, t int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) < t {
		return nil, fmt.Errorf("expected %d outputs, got %d", t, len(fields))
	}
	ans := make([]int64, t)
	for i := 0; i < t; i++ {
		if _, err := fmt.Sscan(fields[i], &ans[i]); err != nil {
			return nil, fmt.Errorf("failed to parse integer at test %d: %v", i+1, err)
		}
	}
	return ans, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[len(os.Args)-1]
	if target == "--" {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}

	inputData, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read input: %v\n", err)
		os.Exit(1)
	}
	reader := strings.NewReader(string(inputData))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read t: %v\n", err)
		os.Exit(1)
	}

	_, src, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(src)
	refPath := filepath.Join(baseDir, "2029E.go")

	refOut, err := runProgram(refPath, inputData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n", err)
		os.Exit(1)
	}
	expected, err := parseOutput(refOut, t)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference output parse error: %v\n", err)
		os.Exit(1)
	}

	out, err := runProgram(target, inputData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target runtime error: %v\n", err)
		os.Exit(1)
	}
	ans, err := parseOutput(out, t)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target output parse error: %v\n", err)
		os.Exit(1)
	}
	for i := 0; i < t; i++ {
		if ans[i] != expected[i] {
			fmt.Fprintf(os.Stderr, "wrong answer on test case %d: expected %d got %d\n", i+1, expected[i], ans[i])
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
