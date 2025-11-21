package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const refSource = "0-999/200-299/250-259/250/250A.go"

type testCase struct {
	name  string
	input string
	arr   []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}
		refK, err := parseOutput(tc, refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candK, err := parseOutput(tc, candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		if candK != refK {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected %d folders got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
				idx+1, tc.name, refK, candK, tc.input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-250A-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref250A.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() {
		_ = os.RemoveAll(dir)
	}
	return binPath, cleanup, nil
}

func runProgram(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutput(tc testCase, output string) (int, error) {
	lines := strings.FieldsFunc(strings.TrimSpace(output), func(r rune) bool { return r == '\n' || r == '\r' })
	if len(lines) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	kLine := strings.TrimSpace(lines[0])
	k, err := strconv.Atoi(kLine)
	if err != nil {
		return 0, fmt.Errorf("invalid k value %q", kLine)
	}
	if k <= 0 {
		return 0, fmt.Errorf("k must be positive, got %d", k)
	}
	var counts []int
	if len(lines) > 1 {
		for _, token := range strings.Fields(lines[1]) {
			val, err := strconv.Atoi(token)
			if err != nil {
				return 0, fmt.Errorf("invalid segment length %q", token)
			}
			counts = append(counts, val)
		}
	} else {
		return 0, fmt.Errorf("missing second line with folder sizes")
	}

	if len(counts) != k {
		return 0, fmt.Errorf("folder count mismatch: declared %d got %d sizes", k, len(counts))
	}
	total := 0
	for _, c := range counts {
		if c <= 0 {
			return 0, fmt.Errorf("non-positive folder size %d", c)
		}
		total += c
	}
	if total != len(tc.arr) {
		return 0, fmt.Errorf("sizes sum %d does not match n=%d", total, len(tc.arr))
	}

	pos := 0
	for idx, size := range counts {
		neg := 0
		for i := 0; i < size; i++ {
			if tc.arr[pos+i] < 0 {
				neg++
				if neg >= 3 {
					return 0, fmt.Errorf("folder %d contains %d negative days", idx+1, neg)
				}
			}
		}
		pos += size
	}

	return k, nil
}

func buildTests() []testCase {
	tests := []testCase{
		newManualTest("n1_all_positive", []int{5}),
		newManualTest("no_negatives", []int{1, 2, 3, 4, 5}),
		newManualTest("exact_two_negatives", []int{-1, -2}),
		newManualTest("third_negative", []int{-1, -2, -3}),
		newManualTest("mixed", []int{1, -1, 2, -2, -3, 4, -4}),
		newManualTest("long_neg_runs", []int{-1, -2, -3, -4, -5, -6}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng, i))
	}
	return tests
}

func newManualTest(name string, arr []int) testCase {
	return testCase{
		name:  name,
		arr:   arr,
		input: formatInput(arr),
	}
}

func randomTest(rng *rand.Rand, idx int) testCase {
	n := rng.Intn(100) + 1
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(201) - 100
	}
	name := fmt.Sprintf("random_%d_n%d", idx+1, n)
	return testCase{
		name:  name,
		arr:   arr,
		input: formatInput(arr),
	}
}

func formatInput(arr []int) string {
	var sb strings.Builder
	sb.Grow(len(arr)*5 + 20)
	sb.WriteString(fmt.Sprintf("%d\n", len(arr)))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}
