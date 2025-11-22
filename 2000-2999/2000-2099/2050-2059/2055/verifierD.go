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

type testCase struct {
	n      int
	k, l   int64
	coords []int64
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[len(os.Args)-1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := generateTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n%s", err, refOut)
		os.Exit(1)
	}
	refVals, err := parseOutputs(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse reference output: %v\n", err)
		os.Exit(1)
	}

	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\n%s", err, candOut)
		os.Exit(1)
	}
	candVals, err := parseOutputs(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse candidate output: %v\n", err)
		os.Exit(1)
	}

	for i := range tests {
		if refVals[i] != candVals[i] {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %d got %d\n", i+1, refVals[i], candVals[i])
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier location")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "ref-2055D-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracle2055D")
	cmd := exec.Command("go", "build", "-o", outPath, "2055D.go")
	cmd.Dir = dir
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("%v\n%s", err, buf.String())
	}
	cleanup := func() { os.RemoveAll(tmpDir) }
	return outPath, cleanup, nil
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return errBuf.String(), err
	}
	if errBuf.Len() > 0 {
		return errBuf.String(), fmt.Errorf("unexpected stderr output")
	}
	return out.String(), nil
}

func parseOutputs(out string, t int) ([]int64, error) {
	tokens := strings.Fields(out)
	if len(tokens) != t {
		return nil, fmt.Errorf("expected %d tokens, got %d", t, len(tokens))
	}
	res := make([]int64, t)
	for i, tok := range tokens {
		v, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("token %q is not integer", tok)
		}
		res[i] = v
	}
	return res, nil
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	totalN := 0
	const maxSumN = 200000

	add := func(tc testCase) {
		if totalN+tc.n > maxSumN {
			return
		}
		tests = append(tests, tc)
		totalN += tc.n
	}

	// Deterministic edge cases and sample-inspired tests
	add(testCase{n: 1, k: 3, l: 5, coords: []int64{0}})
	add(testCase{n: 3, k: 2, l: 5, coords: []int64{2, 5, 5}})
	add(testCase{n: 1, k: 10, l: 10, coords: []int64{10}})
	add(testCase{n: 10, k: 1, l: 100, coords: []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}})
	add(testCase{n: 2, k: 1, l: 20, coords: []int64{0, 20}})
	add(testCase{n: 2, k: 1, l: 20, coords: []int64{2, 20}})
	add(testCase{n: 2, k: 1, l: 30, coords: []int64{2, 30}})
	add(testCase{n: 2, k: 2, l: 4, coords: []int64{1, 4}})
	add(testCase{n: 2, k: 2, l: 4, coords: []int64{5, 5}})

	for len(tests) < 120 && totalN < maxSumN {
		n := rng.Intn(4000) + 1
		if totalN+n > maxSumN {
			n = maxSumN - totalN
		}
		k := int64(rng.Intn(100_000_000) + 1)
		l := int64(rng.Intn(100_000_000) + 1)
		if l < k {
			l = k
		}
		coords := make([]int64, n)
		prev := int64(0)
		for i := 0; i < n; i++ {
			// keep positions sorted and within [0, l]
			step := int64(rng.Intn(5) + rng.Intn(1000))
			pos := prev + step
			if pos > l {
				pos = l
			}
			coords[i] = pos
			prev = pos
		}
		if coords[0] > 0 && rng.Intn(2) == 0 {
			coords[0] = 0 // ensure sometimes a scarecrow at origin
		}
		add(testCase{n: n, k: k, l: l, coords: coords})
	}

	return tests
}

func buildInput(tests []testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%d %d %d\n", tc.n, tc.k, tc.l)
		for i, v := range tc.coords {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", v)
		}
		b.WriteByte('\n')
	}
	return b.String()
}
