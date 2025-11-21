package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	name  string
	input string
	count int
}

type singleCase struct {
	arr []int64
	x   int64
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2164A-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleA")
	cmd := exec.Command("go", "build", "-o", outPath, "2164A.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func formatInput(cases []singleCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(cases)))
	sb.WriteByte('\n')
	for _, cs := range cases {
		sb.WriteString(strconv.Itoa(len(cs.arr)))
		sb.WriteByte('\n')
		for i, v := range cs.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		sb.WriteString(strconv.FormatInt(cs.x, 10))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func generateTests() []testCase {
	const limit int64 = 1_000_000_000
	tests := []testCase{
		{
			name: "single_equal",
			input: formatInput([]singleCase{
				{arr: []int64{0}, x: 0},
			}),
			count: 1,
		},
		{
			name: "single_not_equal",
			input: formatInput([]singleCase{
				{arr: []int64{5}, x: 3},
			}),
			count: 1,
		},
		{
			name: "multi_mixed",
			input: formatInput([]singleCase{
				{arr: []int64{-5, 10, 3}, x: 7},
				{arr: []int64{-5, 10, 3}, x: -6},
				{arr: []int64{-5, 10, 3}, x: 10},
			}),
			count: 3,
		},
		{
			name: "extremes",
			input: formatInput([]singleCase{
				{arr: []int64{-limit, limit}, x: 0},
				{arr: []int64{-limit, limit}, x: limit},
				{arr: []int64{-limit, limit}, x: limit + 1},
			}),
			count: 3,
		},
	}

	// Alternating large test with n=100 and T=1
	alt := make([]int64, 100)
	for i := range alt {
		if i%2 == 0 {
			alt[i] = limit
		} else {
			alt[i] = -limit
		}
	}
	tests = append(tests, testCase{
		name: "wide_range",
		input: formatInput([]singleCase{
			{arr: alt, x: 1},
		}),
		count: 1,
	})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 50; i++ {
		tcCount := rng.Intn(10) + 1
		cases := make([]singleCase, tcCount)
		for j := 0; j < tcCount; j++ {
			n := rng.Intn(100) + 1
			arr := make([]int64, n)
			var minVal, maxVal int64
			for k := 0; k < n; k++ {
				val := rng.Int63n(2*limit+1) - limit
				arr[k] = val
				if k == 0 || val < minVal {
					minVal = val
				}
				if k == 0 || val > maxVal {
					maxVal = val
				}
			}
			var x int64
			switch rng.Intn(3) {
			case 0:
				// inside range
				if maxVal == minVal {
					x = minVal
				} else {
					x = rng.Int63n(maxVal-minVal+1) + minVal
				}
			case 1:
				x = maxVal + rng.Int63n(limit) + 1
			default:
				x = minVal - rng.Int63n(limit) - 1
			}
			cases[j] = singleCase{arr: arr, x: x}
		}
		tests = append(tests, testCase{
			name:  fmt.Sprintf("random_small_batch_%d", i+1),
			input: formatInput(cases),
			count: len(cases),
		})
	}

	// Large batch hitting upper limits
	const maxT = 500
	largeCases := make([]singleCase, maxT)
	for i := 0; i < maxT; i++ {
		n := 100
		arr := make([]int64, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Int63n(2*limit+1) - limit
		}
		x := rng.Int63n(2*limit+1) - limit
		largeCases[i] = singleCase{arr: arr, x: x}
	}
	tests = append(tests, testCase{
		name:  "max_constraints",
		input: formatInput(largeCases),
		count: maxT,
	})

	return tests
}

func normalizeAnswers(output string) []string {
	fields := strings.Fields(output)
	res := make([]string, len(fields))
	for i, f := range fields {
		res[i] = strings.ToUpper(f)
	}
	return res
}

func compareOutputs(expected, actual string, count int) error {
	exp := normalizeAnswers(expected)
	act := normalizeAnswers(actual)
	if len(exp) != count {
		return fmt.Errorf("oracle produced %d answers, expected %d", len(exp), count)
	}
	if len(act) != count {
		return fmt.Errorf("expected %d answers, got %d", count, len(act))
	}
	for i := 0; i < count; i++ {
		if exp[i] != act[i] {
			return fmt.Errorf("mismatch at case %d: expected %s got %s", i+1, exp[i], act[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := generateTests()
	for idx, tc := range tests {
		expected, err := runBinary(oracle, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		actual, err := runBinary(target, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) runtime error: %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		if err := compareOutputs(expected, actual, tc.count); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: %v\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", idx+1, tc.name, err, tc.input, expected, actual)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
