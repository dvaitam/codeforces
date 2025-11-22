package main

import (
	"bufio"
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
	refSource2057H = "2000-2999/2000-2099/2050-2059/2057/2057H.go"
	maxTotalN      = 250000
)

type testCase struct {
	n int
	a []uint64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference(refSource2057H)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	for idx, tc := range tests {
		input := serialize(tc)

		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		expected, err := parseOutput(refOut, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d: %v\ninput:\n%soutput:\n%s", idx+1, err, input, refOut)
			os.Exit(1)
		}

		candOut, runErr := runProgram(candidate, input)
		if runErr != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\ninput:\n%s", idx+1, runErr, input)
			os.Exit(1)
		}
		got, parseErr := parseOutput(candOut, tc.n)
		if parseErr != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d: %v\ninput:\n%soutput:\n%s", idx+1, parseErr, input, candOut)
			os.Exit(1)
		}

		if err := compareSlices(expected, got); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%soutput:\n%s", idx+1, err, input, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference(source string) (string, error) {
	tmp, err := os.CreateTemp("", "2057H-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	srcPath, err := resolveSourcePath(source)
	if err != nil {
		os.Remove(tmp.Name())
		return "", err
	}

	cmd := exec.Command("go", "build", "-o", tmp.Name(), srcPath)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, buf.String())
	}
	return tmp.Name(), nil
}

func resolveSourcePath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(cwd, path), nil
}

func runProgram(target, input string) (string, error) {
	cmd := commandFor(target)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return stdout.String(), nil
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func parseOutput(out string, n int) ([]uint64, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	res := make([]uint64, 0, n)
	for i := 0; i < n; i++ {
		var v uint64
		if _, err := fmt.Fscan(reader, &v); err != nil {
			return nil, fmt.Errorf("expected %d numbers, got %d (%v)", n, i, err)
		}
		res = append(res, v)
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("extra output detected starting with %q", extra)
	}
	return res, nil
}

func compareSlices(expect, got []uint64) error {
	if len(expect) != len(got) {
		return fmt.Errorf("length mismatch: expected %d values, got %d", len(expect), len(got))
	}
	for i := range expect {
		if expect[i] != got[i] {
			return fmt.Errorf("mismatch at position %d: expected %d, got %d", i+1, expect[i], got[i])
		}
	}
	return nil
}

func serialize(tc testCase) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte('\n')
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatUint(v, 10))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func buildTests() []testCase {
	tests := []testCase{
		{n: 2, a: []uint64{8, 0}},
		{n: 5, a: []uint64{2, 2, 2, 2, 2}},
		{n: 5, a: []uint64{0, 0, 9, 0, 0}},
		{n: 1, a: []uint64{0}},
		{n: 1, a: []uint64{1000000000}},
		{n: 3, a: []uint64{1, 0, 1}},
		{n: 4, a: []uint64{1000000000, 0, 1000000000, 0}},
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	total := 0
	for _, tc := range tests {
		total += tc.n
	}

	add := func(tc testCase) {
		tests = append(tests, tc)
		total += tc.n
	}

	add(randomCase(rng, 12, 5, 9))
	add(randomCase(rng, 25, 3, 20))
	add(randomCase(rng, 50, 1, 100))
	add(randomCase(rng, 200, 0, 1000000000))
	add(waveCase(30, 1000))

	for total < maxTotalN {
		n := rng.Intn(4000) + 1
		maxVal := uint64(1_000_000_000)
		if rng.Intn(4) == 0 {
			maxVal = 3
		}
		add(randomCase(rng, n, 0, maxVal))
	}

	return tests
}

func randomCase(rng *rand.Rand, n int, minVal, maxVal uint64) testCase {
	arr := make([]uint64, n)
	for i := 0; i < n; i++ {
		arr[i] = minVal + uint64(rng.Int63n(int64(maxVal-minVal+1)))
	}
	return testCase{n: n, a: arr}
}

func waveCase(n int, peak uint64) testCase {
	arr := make([]uint64, n)
	for i := 0; i < n; i++ {
		if i%2 == 0 {
			arr[i] = peak
		} else {
			arr[i] = 0
		}
	}
	return testCase{n: n, a: arr}
}
