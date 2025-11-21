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

func parseAnswers(out string, totalOps int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) < totalOps {
		return nil, fmt.Errorf("expected %d numbers, got %d", totalOps, len(fields))
	}
	ans := make([]int, totalOps)
	for i := 0; i < totalOps; i++ {
		var v int
		if _, err := fmt.Sscan(fields[i], &v); err != nil {
			return nil, fmt.Errorf("failed to parse integer: %v", err)
		}
		ans[i] = v
	}
	return ans, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[len(os.Args)-1]
	if target == "--" {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
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
	totalOps := 0
	for i := 0; i < t; i++ {
		var n, c0 int
		fmt.Fscan(reader, &n, &c0)
		totalOps += n - 1
		for j := 0; j < n-1; j++ {
			var x int
			fmt.Fscan(reader, &x)
		}
		for j := 0; j < n-1; j++ {
			var u, v int
			fmt.Fscan(reader, &u, &v)
		}
	}

	_, src, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(src)
	refPath := filepath.Join(baseDir, "1790F.go")

	refOut, err := runProgram(refPath, inputData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n", err)
		os.Exit(1)
	}
	refAns, err := parseAnswers(refOut, totalOps)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference output parse error: %v\n", err)
		os.Exit(1)
	}

	targetOut, err := runProgram(target, inputData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target runtime error: %v\n", err)
		os.Exit(1)
	}
	ans, err := parseAnswers(targetOut, totalOps)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target output parse error: %v\n", err)
		os.Exit(1)
	}

	for i := 0; i < totalOps; i++ {
		if ans[i] != refAns[i] {
			fmt.Fprintf(os.Stderr, "wrong answer at position %d: expected %d got %d\n", i+1, refAns[i], ans[i])
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
