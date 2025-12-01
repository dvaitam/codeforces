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
	refSourceC3 = "./331C3.go"
	testCount   = 200
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC3.go /path/to/binary-or-source")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := generateTests()
	for i, n := range tests {
		input := fmt.Sprintf("%d\n", n)

		refOut, err := runExecutable(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (n=%d): %v\n", i+1, n, err)
			os.Exit(1)
		}
		expected, err := parseSingleInt(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid reference output on test %d (n=%d): %v\noutput:\n%s\n", i+1, n, err, refOut)
			os.Exit(1)
		}

		candOut, err := runCandidate(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (n=%d): %v\ninput:\n%s\n", i+1, n, err, input)
			os.Exit(1)
		}
		got, err := parseSingleInt(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d (n=%d): %v\noutput:\n%s\n", i+1, n, err, candOut)
			os.Exit(1)
		}
		if expected != got {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (n=%d): expected %d got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n",
				i+1, n, expected, got, input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	tmp, err := os.CreateTemp("", "ref331C3-*")
	if err != nil {
		return "", nil, fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmp.Name()
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmpPath, refSourceC3)
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmpPath)
		return "", nil, fmt.Errorf("go build failed: %v\n%s", err, string(out))
	}
	cleanup := func() { os.Remove(tmpPath) }
	return tmpPath, cleanup, nil
}

func runExecutable(path string, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func runCandidate(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		abs, err := filepath.Abs(target)
		if err != nil {
			return "", fmt.Errorf("resolve path: %w", err)
		}
		cmd = exec.Command("go", "run", abs)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func parseSingleInt(out string) (int64, error) {
	fields := strings.Fields(out)
	if len(fields) != 1 {
		return 0, fmt.Errorf("expected 1 integer, got %d tokens", len(fields))
	}
	val, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, err
	}
	if val < 0 {
		return 0, fmt.Errorf("negative answer %d", val)
	}
	return val, nil
}

func generateTests() []int64 {
	tests := []int64{0, 1, 2, 9, 10, 11, 19, 24, 55, 99, 100, 101, 12345, 999999999, 1000000000000000000}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < testCount {
		switch rng.Intn(4) {
		case 0:
			tests = append(tests, int64(rng.Intn(1000)))
		case 1:
			tests = append(tests, rng.Int63n(1_000_000_000))
		case 2:
			tests = append(tests, rng.Int63n(1_000_000_000_000_000))
		default:
			tests = append(tests, rng.Int63n(1_000_000_000_000_000_000))
		}
	}
	return tests
}
