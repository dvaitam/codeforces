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

const referenceSource = "2000-2999/2000-2099/2090-2099/2091/2091A.go"

type testCase struct {
	name  string
	input string
	t     int
}

func main() {
	candidatePath, err := candidatePathFromArgs()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	candidate, candCleanup, err := prepareCandidateBinary(candidatePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer candCleanup()

	refPath := referencePath()
	refBin, refCleanup, err := buildReferenceBinary(refPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer refCleanup()

	tests := buildTests()
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\n", idx+1, tc.name, err)
			fmt.Fprintln(os.Stderr, previewInput(tc.input))
			os.Exit(1)
		}
		expect, err := parseOutputs(refOut, tc.t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d (%s): %v\nraw output:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\n", idx+1, tc.name, err)
			fmt.Fprintln(os.Stderr, previewInput(tc.input))
			os.Exit(1)
		}
		have, err := parseOutputs(candOut, tc.t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\nraw output:\n%s\n", idx+1, tc.name, err, candOut)
			fmt.Fprintln(os.Stderr, previewInput(tc.input))
			os.Exit(1)
		}

		for i := 0; i < tc.t; i++ {
			if expect[i] != have[i] {
				fmt.Fprintf(os.Stderr, "test %d (%s) failed on case %d: expected %d got %d\n", idx+1, tc.name, i+1, expect[i], have[i])
				fmt.Fprintln(os.Stderr, previewInput(tc.input))
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func candidatePathFromArgs() (string, error) {
	if len(os.Args) != 2 {
		return "", fmt.Errorf("usage: go run verifierA.go /path/to/binary-or-source")
	}
	return os.Args[1], nil
}

func referencePath() string {
	if _, file, _, ok := runtime.Caller(0); ok {
		return filepath.Join(filepath.Dir(file), "2091A.go")
	}
	return referenceSource
}

func buildReferenceBinary(src string) (string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "verifier2091A-ref")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	bin := filepath.Join(tmpDir, "ref2091a")
	cmd := exec.Command("go", "build", "-o", bin, src)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	cleanup := func() { _ = os.RemoveAll(tmpDir) }
	return bin, cleanup, nil
}

func prepareCandidateBinary(path string) (string, func(), error) {
	if !strings.HasSuffix(path, ".go") {
		return path, func() {}, nil
	}
	tmpDir, err := os.MkdirTemp("", "verifier2091A-cand")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	bin := filepath.Join(tmpDir, "candidate")
	cmd := exec.Command("go", "build", "-o", bin, path)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build candidate: %v\n%s", err, out.String())
	}
	cleanup := func() { _ = os.RemoveAll(tmpDir) }
	return bin, cleanup, nil
}

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return errBuf.String(), fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return out.String(), nil
}

func parseOutputs(out string, t int) ([]int, error) {
	tokens := strings.Fields(out)
	if len(tokens) != t {
		return nil, fmt.Errorf("expected %d integers, got %d", t, len(tokens))
	}
	res := make([]int, t)
	for i, tok := range tokens {
		val, err := strconv.Atoi(tok)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", tok)
		}
		if val < 0 {
			return nil, fmt.Errorf("negative answer %d", val)
		}
		res[i] = val
	}
	return res, nil
}

func buildTests() []testCase {
	var tests []testCase

	tests = append(tests, testCase{
		name:  "samples",
		t:     4,
		input: "4\n10\n2 0 1 2 3 2 5 0 0 1\n8\n2 0 1 2 3 2 5 0\n8\n2 0 1 0 3 2 5 0\n16\n2 3 1 2 3 0 1 9 2 1 0 3 5 4 0 3\n",
	})

	tests = append(tests, testCase{
		name: "minimal-lengths",
		t:    3,
		input: strings.Join([]string{
			"3",
			"1",
			"0",
			"7",
			"0 1 0 3 2 0 2",
			"8",
			"5 2 0 2 3 0 1 0",
		}, "\n") + "\n",
	})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 150; i++ {
		tests = append(tests, randomTest(rng, i, 10))
	}
	for i := 0; i < 30; i++ {
		tests = append(tests, randomTest(rng, i, 20))
	}
	return tests
}

func randomTest(rng *rand.Rand, idx, maxT int) testCase {
	t := rng.Intn(maxT-1) + 2
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n := rng.Intn(20) + 1
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(rng.Intn(10)))
		}
		sb.WriteByte('\n')
	}
	return testCase{
		name:  fmt.Sprintf("random-%d", idx+1),
		t:     t,
		input: sb.String(),
	}
}

func previewInput(in string) string {
	const limit = 400
	if len(in) <= limit {
		return in
	}
	return in[:limit] + "..."
}
