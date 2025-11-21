package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
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

func parseOutput(out string, n int) ([]int, []int, error) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	filtered := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			filtered = append(filtered, line)
		}
	}
	if len(filtered) == 1 && strings.TrimSpace(filtered[0]) == "-1" {
		return []int{-1}, []int{-1}, nil
	}
	if len(filtered) != 2 {
		return nil, nil, fmt.Errorf("expected 2 lines of output, got %d", len(filtered))
	}
	parseLine := func(line string) ([]int, error) {
		fields := strings.Fields(line)
		if len(fields) != n {
			return nil, fmt.Errorf("expected %d numbers, got %d", n, len(fields))
		}
		res := make([]int, n)
		for i := 0; i < n; i++ {
			if _, err := fmt.Sscan(fields[i], &res[i]); err != nil {
				return nil, fmt.Errorf("failed to parse integer: %v", err)
			}
		}
		return res, nil
	}
	A, err := parseLine(filtered[0])
	if err != nil {
		return nil, nil, err
	}
	B, err := parseLine(filtered[1])
	if err != nil {
		return nil, nil, err
	}
	return A, B, nil
}

func mergeCheck(A, B []int, C []int) bool {
	n := len(A)
	merge := make([]int, 0, 2*n)
	i, j := 0, 0
	for len(merge) < 2*n {
		if j == n {
			merge = append(merge, A[i])
			i++
		} else if i == n {
			merge = append(merge, B[j])
			j++
		} else if A[i] <= B[j] {
			merge = append(merge, A[i])
			i++
		} else {
			merge = append(merge, B[j])
			j++
		}
	}
	for idx := range merge {
		if merge[idx] != C[idx] {
			return false
		}
	}
	return true
}

func isPermutation(arr []int, limit int) bool {
	if len(arr) != limit {
		return false
	}
	tmp := append([]int(nil), arr...)
	sort.Ints(tmp)
	for i, v := range tmp {
		if v != i+1 {
			return false
		}
	}
	return true
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
	reader := bufio.NewReader(bytes.NewReader(inputData))
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read n: %v\n", err)
		os.Exit(1)
	}
	C := make([]int, 2*n)
	for i := 0; i < 2*n; i++ {
		if _, err := fmt.Fscan(reader, &C[i]); err != nil {
			fmt.Fprintf(os.Stderr, "failed to read C[%d]: %v\n", i, err)
			os.Exit(1)
		}
	}
	if !isPermutation(C, 2*n) {
		fmt.Fprintf(os.Stderr, "input array C is not a permutation\n")
		os.Exit(1)
	}

	_, src, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(src)
	refPath := filepath.Join(baseDir, "1906E.go")

	refOut, err := runProgram(refPath, inputData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n", err)
		os.Exit(1)
	}
	refA, refB, err := parseOutput(refOut, n)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference output parse error: %v\n", err)
		os.Exit(1)
	}
	if refA[0] == -1 && refB[0] == -1 {
		// reference claims impossible, so target must also output -1 (case-insensitive)
		targetOut, err := runProgram(target, inputData)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error: %v\n", err)
			os.Exit(1)
		}
		if strings.TrimSpace(targetOut) != "-1" {
			fmt.Fprintf(os.Stderr, "expected -1, got:\n%s", targetOut)
			os.Exit(1)
		}
		fmt.Println("all tests passed")
		return
	}
	// reference gives valid solution, so verify target as well
	targetOut, err := runProgram(target, inputData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target runtime error: %v\n", err)
		os.Exit(1)
	}
	if strings.TrimSpace(targetOut) == "-1" {
		fmt.Fprintf(os.Stderr, "target output -1 but solution exists\n")
		os.Exit(1)
	}
	A, B, err := parseOutput(targetOut, n)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target output parse error: %v\n", err)
		os.Exit(1)
	}
	if !mergeCheck(A, B, C) {
		fmt.Fprintf(os.Stderr, "target arrays do not merge to C\n")
		os.Exit(1)
	}
	fmt.Println("all tests passed")
}
