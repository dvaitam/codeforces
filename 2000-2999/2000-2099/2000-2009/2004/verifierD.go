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

const refSource = "2004D.go"

type testCase struct {
	colors  []string
	queries [][2]int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := append(deterministicTests(), randomTests()...)
	input := buildInput(tests)
	expectedOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}
	actualOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\ninput:\n%s", err, input)
		os.Exit(1)
	}

	expected, err := parseOutput(expectedOut, lenAnswers(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference output invalid: %v\n%s", err, expectedOut)
		os.Exit(1)
	}
	actual, err := parseOutput(actualOut, lenAnswers(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate output invalid: %v\n%s", err, actualOut)
		os.Exit(1)
	}

	if err := compare(expected, actual); err != nil {
		fmt.Fprintf(os.Stderr, "%v\ninput:\n%s", err, input)
		os.Exit(1)
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)

	tmp, err := os.CreateTemp("", "2004D-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	cmd.Dir = dir
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, buf.String())
	}
	return tmp.Name(), nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func deterministicTests() []testCase {
	return []testCase{
		// Single city, only zero-cost queries.
		{
			colors:  []string{"BG"},
			queries: [][2]int{{1, 1}, {1, 1}},
		},
		// No possible path between distinct colors.
		{
			colors:  []string{"BG", "RY"},
			queries: [][2]int{{1, 2}, {2, 1}},
		},
		// Direct shared color beats detours.
		{
			colors:  []string{"BG", "BR", "GY", "GR"},
			queries: [][2]int{{1, 3}, {2, 4}, {1, 4}},
		},
		// Best path through an intermediate city.
		{
			colors:  []string{"BR", "GY", "RY"},
			queries: [][2]int{{1, 3}, {2, 3}},
		},
		// Mix of reachable and unreachable pairs.
		{
			colors:  []string{"BG", "GR", "RY", "BY", "BR"},
			queries: [][2]int{{1, 5}, {2, 4}, {3, 4}, {1, 3}},
		},
	}
}

func randomTests() []testCase {
	pairs := []string{"BG", "BR", "BY", "GR", "GY", "RY"}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	build := func() testCase {
		n := rng.Intn(40) + 1
		colors := make([]string, n)
		for i := 0; i < n; i++ {
			colors[i] = pairs[rng.Intn(len(pairs))]
		}
		q := rng.Intn(60) + 1
		queries := make([][2]int, q)
		for i := 0; i < q; i++ {
			x := rng.Intn(n) + 1
			y := rng.Intn(n) + 1
			queries[i] = [2]int{x, y}
		}
		return testCase{colors: colors, queries: queries}
	}

	tests := make([]testCase, 0, 200)
	for len(tests) < cap(tests) {
		tests = append(tests, build())
	}
	return tests
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.Grow(len(tests) * 128)
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(len(tc.colors)))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(len(tc.queries)))
		sb.WriteByte('\n')
		for i, c := range tc.colors {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(c)
		}
		sb.WriteByte('\n')
		for _, q := range tc.queries {
			sb.WriteString(strconv.Itoa(q[0]))
			sb.WriteByte(' ')
			sb.WriteString(strconv.Itoa(q[1]))
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func parseOutput(out string, expected int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d answers, got %d", expected, len(fields))
	}
	res := make([]int, expected)
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("answer %d is not an integer: %s", i+1, f)
		}
		res[i] = v
	}
	return res, nil
}

func compare(expected, actual []int) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("answer count mismatch: expected %d got %d", len(expected), len(actual))
	}
	for i := range expected {
		if expected[i] != actual[i] {
			return fmt.Errorf("mismatch at answer %d: expected %d got %d", i+1, expected[i], actual[i])
		}
	}
	return nil
}

func lenAnswers(tests []testCase) int {
	total := 0
	for _, tc := range tests {
		total += len(tc.queries)
	}
	return total
}
