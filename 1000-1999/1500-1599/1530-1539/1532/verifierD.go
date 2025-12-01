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

const refSource = "1532D.go"

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
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
		expect, err := runAndParse(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		got, err := runAndParse(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected %d got %d\ninput:\n%s", idx+1, tc.name, expect, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-1532D-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref1532D.bin")
	source := filepath.Join(".", refSource)
	cmd := exec.Command("go", "build", "-o", binPath, source)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		_ = os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
}

func runAndParse(bin, input string) (int64, error) {
	out, err := runProgram(bin, input)
	if err != nil {
		return 0, err
	}
	return parseAnswer(out)
}

func runProgram(bin, input string) (string, error) {
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
		return 0, fmt.Errorf("expected single integer, got %d tokens (output: %q)", len(fields), output)
	}
	val, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q", fields[0])
	}
	if val < 0 {
		return 0, fmt.Errorf("answer must be non-negative, got %d", val)
	}
	return val, nil
}

func buildTests() []testCase {
	tests := []testCase{
		newTestCase("sample1", 6, []int{5, 10, 2, 3, 14, 5}),
		newTestCase("sample2", 2, []int{1, 100}),
		newTestCase("all_equal", 8, []int{7, 7, 7, 7, 7, 7, 7, 7}),
		newTestCase("already_pairs", 4, []int{1, 1, 2, 2}),
		newTestCase("spread", 10, []int{1, 10, 2, 9, 3, 8, 4, 7, 5, 6}),
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTestCase(rng, i))
	}
	return tests
}

func newTestCase(name string, n int, skills []int) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range skills {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return testCase{name: name, input: sb.String()}
}

func randomTestCase(rng *rand.Rand, idx int) testCase {
	n := (rng.Intn(50) + 1) * 2 // even between 2 and 100
	skills := make([]int, n)
	for i := range skills {
		skills[i] = rng.Intn(100) + 1
	}
	return newTestCase(fmt.Sprintf("random_%d_n%d", idx+1, n), n, skills)
}
