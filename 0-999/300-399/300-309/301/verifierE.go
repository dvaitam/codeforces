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

const refSource = "0-999/300-399/300-309/301/301E.go"
const mod = 1000000007

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
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
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, expectRaw)
			os.Exit(1)
		}
		expect, err := parseAnswer(expectRaw)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, expectRaw)
			os.Exit(1)
		}

		gotRaw, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, gotRaw)
			os.Exit(1)
		}
		got, err := parseAnswer(gotRaw)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, gotRaw)
			os.Exit(1)
		}
		if expect != got {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected %d got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
				idx+1, tc.name, expect, got, tc.input, expectRaw, gotRaw)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-301E-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref301E.bin")
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

func buildTests() []testCase {
	tests := make([]testCase, 0, 1500)
	for n := 1; n <= 10; n++ {
		for m := 1; m <= 10; m++ {
			for k := 1; k <= 10; k++ {
				name := fmt.Sprintf("det_n%d_m%d_k%d", n, m, k)
				tests = append(tests, testCase{name: name, input: fmt.Sprintf("%d %d %d\n", n, m, k)})
			}
		}
	}
	edgeCases := []struct {
		n, m, k int
		name    string
	}{
		{1, 1, 1, "all_min"},
		{100, 100, 100, "all_max"},
		{100, 1, 100, "max_n_min_m"},
		{1, 100, 100, "min_n_max_m"},
		{100, 100, 1, "k_min"},
	}
	for _, e := range edgeCases {
		tests = append(tests, testCase{
			name:  e.name,
			input: fmt.Sprintf("%d %d %d\n", e.n, e.m, e.k),
		})
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 300; i++ {
		n := rng.Intn(100) + 1
		m := rng.Intn(100) + 1
		k := rng.Intn(100) + 1
		tests = append(tests, testCase{
			name:  fmt.Sprintf("random_%d", i+1),
			input: fmt.Sprintf("%d %d %d\n", n, m, k),
		})
	}
	return tests
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
	if val < 0 || val >= mod {
		return 0, fmt.Errorf("answer %d out of [0,%d)", val, mod)
	}
	return val, nil
}
