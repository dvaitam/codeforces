package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
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

func parseGraph(data []byte) (int, error) {
	reader := bytes.NewReader(data)
	var u, v int
	maxNode := -1
	for {
		_, err := fmt.Fscan(reader, &u, &v)
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, fmt.Errorf("failed to parse graph: %v", err)
		}
		if u > maxNode {
			maxNode = u
		}
		if v > maxNode {
			maxNode = v
		}
	}
	if maxNode < 0 {
		return 0, fmt.Errorf("graph has no edges")
	}
	return maxNode + 1, nil
}

func validatePartition(output string, n int) error {
	lines := strings.Split(output, "\n")
	seen := make([]bool, n)
	covered := 0
	communities := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		communities++
		for _, f := range fields {
			idx, err := strconv.Atoi(f)
			if err != nil {
				return fmt.Errorf("invalid node index %q", f)
			}
			if idx < 0 || idx >= n {
				return fmt.Errorf("node index %d out of range [0,%d)", idx, n)
			}
			if seen[idx] {
				return fmt.Errorf("node %d appears in multiple communities", idx)
			}
			seen[idx] = true
			covered++
		}
	}
	if communities == 0 {
		return fmt.Errorf("no communities provided")
	}
	if covered != n {
		return fmt.Errorf("expected %d nodes, but got %d", n, covered)
	}
	for i := 0; i < n; i++ {
		if !seen[i] {
			return fmt.Errorf("node %d missing from partition", i)
		}
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA1.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[len(os.Args)-1]
	if target == "--" {
		fmt.Println("usage: go run verifierA1.go /path/to/binary")
		os.Exit(1)
	}
	inputData, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read input: %v\n", err)
		os.Exit(1)
	}
	n, err := parseGraph(inputData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	_, srcFile, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(srcFile)
	refPath := filepath.Join(baseDir, "1377A1.go")
	refOut, err := runProgram(refPath, inputData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference solution failed: %v\n", err)
		os.Exit(1)
	}
	if err := validatePartition(refOut, n); err != nil {
		fmt.Fprintf(os.Stderr, "reference output invalid: %v\n", err)
		os.Exit(1)
	}

	out, err := runProgram(target, inputData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target runtime error: %v\n", err)
		os.Exit(1)
	}
	if err := validatePartition(out, n); err != nil {
		fmt.Fprintf(os.Stderr, "invalid partition: %v\noutput:\n%s", err, out)
		os.Exit(1)
	}
	fmt.Println("all tests passed")
}
