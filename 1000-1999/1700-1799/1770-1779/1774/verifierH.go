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
	"strings"
)

type testCase struct {
	n int
	k int
}

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

func parseInput(data []byte) ([]testCase, error) {
	reader := bufio.NewReader(bytes.NewReader(data))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, fmt.Errorf("failed to read t: %v", err)
	}
	tests := make([]testCase, 0, t)
	for ; t > 0; t-- {
		var n, k int
		if _, err := fmt.Fscan(reader, &n, &k); err != nil {
			return nil, fmt.Errorf("failed to read n,k: %v", err)
		}
		tests = append(tests, testCase{n: n, k: k})
	}
	return tests, nil
}

func parseReference(out string, tests []testCase) ([]int64, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	best := make([]int64, len(tests))
	for i, tc := range tests {
		if _, err := fmt.Fscan(reader, &best[i]); err != nil {
			return nil, fmt.Errorf("reference: failed to read value for case %d: %v", i+1, err)
		}
		for j := 0; j < tc.n; j++ {
			var tmp int
			if _, err := fmt.Fscan(reader, &tmp); err != nil {
				return nil, fmt.Errorf("reference: failed to read permutation value for case %d: %v", i+1, err)
			}
		}
	}
	return best, nil
}

func validatePermutation(arr []int, tc testCase) error {
	n := tc.n
	if len(arr) != n {
		return fmt.Errorf("expected permutation length %d, got %d", n, len(arr))
	}
	seen := make([]bool, n+1)
	for idx, v := range arr {
		if v < 1 || v > n {
			return fmt.Errorf("value %d at position %d out of range [1,%d]", v, idx+1, n)
		}
		if seen[v] {
			return fmt.Errorf("value %d appears multiple times", v)
		}
		seen[v] = true
	}
	return nil
}

func minWindowSum(arr []int, k int) int64 {
	if k > len(arr) {
		return 0
	}
	var sum int64
	for i := 0; i < k; i++ {
		sum += int64(arr[i])
	}
	best := sum
	for i := k; i < len(arr); i++ {
		sum += int64(arr[i]) - int64(arr[i-k])
		if sum < best {
			best = sum
		}
	}
	return best
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[len(os.Args)-1]
	if target == "--" {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}

	inputData, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read input: %v\n", err)
		os.Exit(1)
	}
	tests, err := parseInput(inputData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	_, src, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(src)
	refPath := filepath.Join(baseDir, "1774H.go")

	refOut, err := runProgram(refPath, inputData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n", err)
		os.Exit(1)
	}
	refVals, err := parseReference(refOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	targetOut, err := runProgram(target, inputData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target runtime error: %v\n", err)
		os.Exit(1)
	}
	reader := bufio.NewReader(strings.NewReader(targetOut))
	for i, tc := range tests {
		var reported int64
		if _, err := fmt.Fscan(reader, &reported); err != nil {
			fmt.Fprintf(os.Stderr, "case %d: failed to read reported value: %v\n", i+1, err)
			os.Exit(1)
		}
		arr := make([]int, tc.n)
		for j := 0; j < tc.n; j++ {
			if _, err := fmt.Fscan(reader, &arr[j]); err != nil {
				fmt.Fprintf(os.Stderr, "case %d: failed to read permutation element %d: %v\n", i+1, j+1, err)
				os.Exit(1)
			}
		}
		if err := validatePermutation(arr, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d invalid permutation: %v\n", i+1, err)
			os.Exit(1)
		}
		minSum := minWindowSum(arr, tc.k)
		if reported != minSum {
			fmt.Fprintf(os.Stderr, "case %d: reported value %d does not match permutation minimal window sum %d\n", i+1, reported, minSum)
			os.Exit(1)
		}
		if reported != refVals[i] {
			fmt.Fprintf(os.Stderr, "case %d: optimal value %d but reported %d\n", i+1, refVals[i], reported)
			os.Exit(1)
		}
	}

	fmt.Println("all tests passed")
}
