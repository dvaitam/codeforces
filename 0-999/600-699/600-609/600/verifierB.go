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

const refSource = "0-999/600-699/600-609/600/600B.go"

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
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}
		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		if err := compareOutputs(refOut, candOut); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: %v\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
				idx+1, tc.name, err, tc.input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-600B-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref600B.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
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

func compareOutputs(refOut, candOut string) error {
	exp := normalizeOutput(refOut)
	got := normalizeOutput(candOut)
	if len(exp) != len(got) {
		return fmt.Errorf("expected %d answers got %d", len(exp), len(got))
	}
	for i := range exp {
		if exp[i] != got[i] {
			return fmt.Errorf("answer %d mismatch: expected %d got %d", i+1, exp[i], got[i])
		}
	}
	return nil
}

func normalizeOutput(out string) []int {
	fields := strings.Fields(out)
	res := make([]int, len(fields))
	for i, f := range fields {
		val, err := strconv.Atoi(f)
		if err != nil {
			return nil
		}
		res[i] = val
	}
	return res
}

func buildTests() []testCase {
	tests := []testCase{
		{name: "single_element", input: formatCase([]int{5}, []int{3, 5, 6})},
		{name: "negative_values", input: formatCase([]int{-5, 0, 10}, []int{-10, -5, 0, 5, 10})},
		{name: "all_equal", input: formatCase([]int{1, 1, 1}, []int{0, 1, 2})},
		{name: "descending_queries", input: formatCase([]int{2, 4, 6, 8}, []int{10, 8, 6, 4, 2, 0})},
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomCase(rng, i))
	}
	tests = append(tests, stressCase())
	return tests
}

func formatCase(a, b []int) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", len(a), len(b))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for i, v := range b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func randomCase(rng *rand.Rand, idx int) testCase {
	n := rng.Intn(2000) + 1
	m := rng.Intn(2000) + 1
	a := make([]int, n)
	b := make([]int, m)
	for i := range a {
		a[i] = rng.Intn(2_000_001) - 1_000_000
	}
	for i := range b {
		b[i] = rng.Intn(2_000_001) - 1_000_000
	}
	return testCase{
		name:  fmt.Sprintf("random_%d", idx+1),
		input: formatCase(a, b),
	}
}

func stressCase() testCase {
	n := 200000
	m := 200000
	a := make([]int, n)
	b := make([]int, m)
	for i := 0; i < n; i++ {
		a[i] = i - n/2
	}
	for i := 0; i < m; i++ {
		b[i] = i - m/2
	}
	return testCase{
		name:  "stress_max",
		input: formatCase(a, b),
	}
}
