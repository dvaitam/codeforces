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

const (
	refSource = "1000-1999/1100-1199/1180-1189/1184/1184D2.go"
	mod       = 1000000007
)

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD2.go /path/to/binary")
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
		refVal, err := parseOutput(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candVal, err := parseOutput(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}

		if refVal != candVal {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected %d got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
				idx+1, tc.name, refVal, candVal, tc.input, refOut, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-1184D2-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref1184D2.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		_ = os.RemoveAll(dir)
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

func parseOutput(output string) (int64, error) {
	trimmed := strings.TrimSpace(output)
	if trimmed == "" {
		return 0, fmt.Errorf("empty output")
	}
	tokens := strings.Fields(trimmed)
	if len(tokens) != 1 {
		return 0, fmt.Errorf("expected single integer, got %q", trimmed)
	}
	val, err := strconv.ParseInt(tokens[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q", tokens[0])
	}
	return val % mod, nil
}

func buildTests() []testCase {
	tests := []testCase{
		makeManualTest("minimal", 1, 1, 1),
		makeManualTest("already_end_left", 5, 1, 5),
		makeManualTest("already_end_right", 6, 6, 6),
		makeManualTest("middle_no_growth", 2, 1, 2),
		makeManualTest("middle_can_grow", 3, 2, 5),
		makeManualTest("large_m", 7, 3, 250),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng, i))
	}
	return tests
}

func makeManualTest(name string, n, k, m int) testCase {
	return testCase{
		name:  name,
		input: fmt.Sprintf("%d %d %d\n", n, k, m),
	}
}

func randomTest(rng *rand.Rand, idx int) testCase {
	n := rng.Intn(250) + 1
	m := n + rng.Intn(251-n)
	k := rng.Intn(n) + 1
	name := fmt.Sprintf("random_%d_n%d_k%d_m%d", idx+1, n, k, m)
	return testCase{
		name:  name,
		input: fmt.Sprintf("%d %d %d\n", n, k, m),
	}
}
