package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const refSource = "./2149A.go"

type testCase struct {
	n   int
	arr []int
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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
		refAns := refAnswers[i]
		candAns := candAnswers[i]
		if refAns < 0 {
			fmt.Fprintf(os.Stderr, "reference answer negative on test %d\n", i+1)
			os.Exit(1)
		}
		if candAns < 0 {
			fmt.Fprintf(os.Stderr, "test %d: candidate output negative number %d\n", i+1, candAns)
			os.Exit(1)
		}

		// Any minimal answer is acceptable, but no answer smaller than reference (minimal) possible
		if candAns != refAns {
			fmt.Fprintf(os.Stderr, "test %d mismatch: expected %d operations, got %d\ninput: n=%d array=%v\n",
				i+1, refAns, candAns, tc.n, tc.arr)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func buildReference() (string, error) {
	outPath := "./ref_2149A.bin"
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
	add := func(arr []int) {
		tests = append(tests, testCase{n: len(arr), arr: append([]int(nil), arr...)})
	}

	baseCases := [][]int{
		{1},
		{-1},
		{0},
		{-1, -1},
		{-1, 0},
		{-1, 1},
		{0, 0},
		{0, 1},
		{1, 1},
	}
	for _, arr := range baseCases {
		add(arr)
	}

	allVals := []int{-1, 0, 1}
	for n := 1; n <= 8; n++ {
		total := int(math.Pow(3, float64(n)))
		if total > 200 {
			break
		}
		for mask := 0; mask < total; mask++ {
			arr := make([]int, n)
			tmp := mask
			for i := 0; i < n; i++ {
				arr[i] = allVals[tmp%3]
				tmp /= 3
			}
			add(arr)
		}
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 200 {
		n := rng.Intn(8) + 1
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = allVals[rng.Intn(3)]
		}
		add(arr)
	}
	return tests
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i, v := range tc.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
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
