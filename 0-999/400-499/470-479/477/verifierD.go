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
	refSource = "0-999/400-499/470-479/477/477D.go"
	mod       = 1000000007
)

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
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}
		refCount, refMin, err := parseOutput(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candCount, candMin, err := parseOutput(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		if candCount != refCount || candMin != refMin {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected (%d,%d) got (%d,%d)\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
				idx+1, tc.name, refCount, refMin, candCount, candMin, tc.input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-477D-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref477D.bin")
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

func parseOutput(out string) (int64, int64, error) {
	lines := splitLines(out)
	if len(lines) != 2 {
		return 0, 0, fmt.Errorf("expected two lines, got %d", len(lines))
	}
	first, err := strconv.ParseInt(lines[0], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid first line %q", lines[0])
	}
	second, err := strconv.ParseInt(lines[1], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid second line %q", lines[1])
	}
	first = modNormalize(first)
	second = modNormalize(second)
	return first, second, nil
}

func modNormalize(x int64) int64 {
	x %= mod
	if x < 0 {
		x += mod
	}
	return x
}

func splitLines(out string) []string {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			result = append(result, line)
		}
	}
	return result
}

func buildTests() []testCase {
	tests := []testCase{
		{name: "1", input: "1\n"},
		{name: "101", input: "101\n"},
		{name: "111", input: "111\n"},
		{name: "1000", input: "1000\n"},
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	// random lengths up to ~200
	for i := 0; i < 200; i++ {
		tests = append(tests, randomCase(rng, i, 200))
	}
	// include large length near limit (~25000)
	tests = append(tests, randomCase(rng, 201, 25000))
	return tests
}

func randomCase(rng *rand.Rand, idx int, maxLen int) testCase {
	length := rng.Intn(maxLen-1) + 1
	if length == 0 {
		length = 1
	}
	bytes := make([]byte, length)
	bytes[0] = '1'
	for i := 1; i < length; i++ {
		if rng.Intn(2) == 0 {
			bytes[i] = '0'
		} else {
			bytes[i] = '1'
		}
	}
	return testCase{
		name:  fmt.Sprintf("random_%d_len_%d", idx+1, length),
		input: string(bytes) + "\n",
	}
}
