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

const refSource = "2057C.go"

type testCase struct {
	name  string
	input string
	lr    [][2]int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/candidate_binary_or_go_file")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanupRef, err := buildBinary(refSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference solution: %v\n", err)
		os.Exit(1)
	}
	defer cleanupRef()

	candBin := candidate
	candCleanup := func() {}
	if strings.HasSuffix(candidate, ".go") {
		candBin, candCleanup, err = buildBinary(candidate)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to build candidate: %v\n", err)
			os.Exit(1)
		}
	}
	defer candCleanup()

	tests := buildTests()
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		candOut, err := runProgram(candBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}

		refTriples, err := parseTriples(refOut, len(tc.lr))
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}
		candTriples, err := parseTriples(candOut, len(tc.lr))
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}

		for caseIdx, bounds := range tc.lr {
			l, r := bounds[0], bounds[1]
			refVal, err := evalTriple(refTriples[caseIdx], l, r)
			if err != nil {
				fmt.Fprintf(os.Stderr, "reference output invalid on test %d case %d: %v\ninput:\n%soutput:\n%s\n", idx+1, caseIdx+1, err, tc.input, refOut)
				os.Exit(1)
			}
			cVal, err := evalTriple(candTriples[caseIdx], l, r)
			if err != nil {
				fmt.Fprintf(os.Stderr, "candidate output invalid on test %d case %d: %v\ninput:\n%soutput:\n%s\n", idx+1, caseIdx+1, err, tc.input, candOut)
				os.Exit(1)
			}
			if cVal != refVal {
				fmt.Fprintf(os.Stderr, "test %d (%s) case %d: expected value %d, got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n",
					idx+1, tc.name, caseIdx+1, refVal, cVal, tc.input, refOut, candOut)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildBinary(src string) (string, func(), error) {
	abs, err := filepath.Abs(src)
	if err != nil {
		return "", nil, err
	}
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot locate verifier path")
	}
	dir := filepath.Dir(file)

	tmpDir, err := os.MkdirTemp("", "2057C-bin-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "bin")

	cmd := exec.Command("go", "build", "-o", binPath, abs)
	cmd.Dir = dir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("build failed: %v\n%s", err, stderr.String())
	}

	cleanup := func() {
		_ = os.RemoveAll(tmpDir)
	}
	return binPath, cleanup, nil
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

func parseTriples(output string, cases int) ([][3]int, error) {
	tokens := strings.Fields(output)
	if len(tokens) != cases*3 {
		return nil, fmt.Errorf("expected %d integers, got %d", cases*3, len(tokens))
	}
	res := make([][3]int, cases)
	for i := 0; i < cases; i++ {
		for j := 0; j < 3; j++ {
			val, err := strconv.Atoi(tokens[i*3+j])
			if err != nil {
				return nil, fmt.Errorf("invalid integer %q", tokens[i*3+j])
			}
			res[i][j] = val
		}
	}
	return res, nil
}

func evalTriple(tr [3]int, l, r int) (int64, error) {
	a, b, c := tr[0], tr[1], tr[2]
	if a == b || a == c || b == c {
		return 0, fmt.Errorf("values must be pairwise distinct")
	}
	for _, v := range []int{a, b, c} {
		if v < l || v > r {
			return 0, fmt.Errorf("value %d out of range [%d,%d]", v, l, r)
		}
	}
	val := int64(a^b) + int64(a^c) + int64(b^c)
	return val, nil
}

func buildTests() []testCase {
	tests := []testCase{
		newManualTest("example_like", [][2]int{{0, 2}}),
		newManualTest("small_span", [][2]int{{5, 8}}),
		newManualTest("mixed", [][2]int{{8, 16}, {2, 5}}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 220; i++ {
		tests = append(tests, randomTest(rng, i))
	}
	return tests
}

func newManualTest(name string, lr [][2]int) testCase {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(lr)))
	sb.WriteByte('\n')
	for _, p := range lr {
		l, r := p[0], p[1]
		sb.WriteString(fmt.Sprintf("%d %d\n", l, r))
	}
	return testCase{name: name, input: sb.String(), lr: lr}
}

func randomTest(rng *rand.Rand, idx int) testCase {
	t := rng.Intn(4) + 1
	lr := make([][2]int, t)
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(t))
	sb.WriteByte('\n')
	for i := 0; i < t; i++ {
		l := rng.Intn(1 << 20)
		diff := rng.Intn(1<<10) + 2
		r := l + diff
		if r >= 1<<30 {
			l -= (r - (1 << 30)) + 1
			r = l + diff
			if l < 0 {
				l = 0
				r = diff
			}
		}
		lr[i] = [2]int{l, r}
		sb.WriteString(fmt.Sprintf("%d %d\n", l, r))
	}
	return testCase{
		name:  fmt.Sprintf("random_%d", idx+1),
		input: sb.String(),
		lr:    lr,
	}
}
