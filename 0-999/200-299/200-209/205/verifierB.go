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

const refSource = "0-999/200-299/200-209/205/205B.go"

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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
		expectRaw, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, expectRaw)
			os.Exit(1)
		}
		expect, err := parseAnswer(expectRaw)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, expectRaw)
			os.Exit(1)
		}

		gotRaw, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, gotRaw)
			os.Exit(1)
		}
		got, err := parseAnswer(gotRaw)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, gotRaw)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected %d got %d\ninput:\n%s", idx+1, tc.name, expect, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-205B-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref205B.bin")
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

func parseAnswer(output string) (int64, error) {
	fields := strings.Fields(output)
	if len(fields) != 1 {
		return 0, fmt.Errorf("expected single integer, got %d tokens", len(fields))
	}
	val, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q", fields[0])
	}
	if val < 0 {
		return 0, fmt.Errorf("negative answer %d", val)
	}
	return val, nil
}

func buildTests() []testCase {
	tests := []testCase{
		newManualTest("single_element", []int64{7}),
		newManualTest("already_sorted", []int64{1, 2, 3, 3, 10}),
		newManualTest("strictly_decreasing", []int64{5, 4, 3, 2, 1}),
		newManualTest("plateau_drop", []int64{3, 3, 3, 1}),
		newManualTest("large_values", []int64{1_000_000_000, 999_999_999, 999_999_999}),
		newManualTest("late_drop", []int64{2, 5, 5, 5, 4}),
	}
	tests = append(tests,
		newCustomTest("n1e5_almost_sorted", buildAlmostSorted(100000)),
		newCustomTest("n1e5_descending", buildDescending(100000)),
	)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng, i))
	}
	return tests
}

func newManualTest(name string, arr []int64) testCase {
	return newCustomTest(name, arr)
}

func newCustomTest(name string, arr []int64) testCase {
	var sb strings.Builder
	sb.Grow(len(arr)*12 + 20)
	sb.WriteString(fmt.Sprintf("%d\n", len(arr)))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	return testCase{name: name, input: sb.String()}
}

func randomTest(rng *rand.Rand, idx int) testCase {
	n := rng.Intn(1000) + 1
	if rng.Intn(20) == 0 {
		n = rng.Intn(100000) + 1
	}
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Int63n(1_000_000_000) + 1
	}
	name := fmt.Sprintf("random_%d_n%d", idx+1, n)
	return newCustomTest(name, arr)
}

func buildAlmostSorted(n int) []int64 {
	arr := make([]int64, n)
	cur := int64(0)
	for i := 0; i < n; i++ {
		cur += int64(i % 3)
		arr[i] = cur
	}
	// introduce occasional drops
	for i := 100; i < n; i += 5000 {
		arr[i] -= 5
	}
	return arr
}

func buildDescending(n int) []int64 {
	arr := make([]int64, n)
	val := int64(1_000_000_000)
	for i := 0; i < n; i++ {
		arr[i] = val - int64(i)
	}
	return arr
}
