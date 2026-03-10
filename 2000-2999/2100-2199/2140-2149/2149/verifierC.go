package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const refSource = "./2149C.go"

type testCase struct {
	n   int
	k   int
	arr []int
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	candidate := args[0]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}
	refAnswers, err := parseOutput(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\noutput:\n%s", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	candAnswers, err := parseOutput(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\noutput:\n%s", err, candOut)
		os.Exit(1)
	}

	for i, tc := range tests {
		if candAnswers[i] != refAnswers[i] {
			fmt.Fprintf(os.Stderr, "test %d mismatch: expected %d, got %d\nn=%d k=%d arr=%v\n",
				i+1, refAnswers[i], candAnswers[i], tc.n, tc.k, tc.arr)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func buildReference() (string, error) {
	outPath := "./ref_2149C.bin"
	cmd := exec.Command("go", "build", "-o", outPath, refSource)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return outPath, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func buildTests() []testCase {
	var tests []testCase

	// Manual edge cases
	add := func(n, k int, arr []int) {
		tests = append(tests, testCase{n: n, k: k, arr: append([]int(nil), arr...)})
	}

	add(1, 0, []int{0})
	add(1, 1, []int{0})
	add(1, 1, []int{1})
	add(3, 1, []int{0, 2, 3})
	add(5, 5, []int{0, 1, 2, 3, 4})
	add(6, 2, []int{0, 3, 4, 2, 6, 2})
	add(7, 4, []int{0, 1, 5, 4, 4, 7, 3})

	// Exhaustive small cases
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for n := 1; n <= 6; n++ {
		for k := 0; k <= n; k++ {
			// Try many random arrays of length n with values in [0, n]
			for trial := 0; trial < 20; trial++ {
				arr := make([]int, n)
				for i := range arr {
					arr[i] = rng.Intn(n + 1)
				}
				add(n, k, arr)
			}
		}
	}

	// Larger random cases
	for trial := 0; trial < 100; trial++ {
		n := rng.Intn(15) + 1
		k := rng.Intn(n + 1)
		arr := make([]int, n)
		for i := range arr {
			arr[i] = rng.Intn(n + 1)
		}
		add(n, k, arr)
	}

	return tests
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
		for i, v := range tc.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutput(out string, expected int) ([]int, error) {
	lines := strings.Fields(out)
	if len(lines) != expected {
		return nil, fmt.Errorf("expected %d outputs, got %d", expected, len(lines))
	}
	ans := make([]int, expected)
	for i, s := range lines {
		val, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", s)
		}
		ans[i] = val
	}
	return ans, nil
}
