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

const refSource = "2003D2.go"

type testCase struct {
	name    string
	input   string
	outputs int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD2.go /path/to/candidate_binary_or_go_file")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference solution: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	candBin := candidate
	candCleanup := func() {}
	if strings.HasSuffix(candidate, ".go") {
		candBin, candCleanup, err = buildGoFile(candidate)
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
		refVals, err := parseOutputs(refOut, tc.outputs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candVals, err := parseOutputs(candOut, tc.outputs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}

		for i := 0; i < tc.outputs; i++ {
			if candVals[i] != refVals[i] {
				fmt.Fprintf(os.Stderr, "test %d (%s) failed at answer %d: expected %d got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n",
					idx+1, tc.name, i+1, refVals[i], candVals[i], tc.input, refOut, candOut)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier location")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2003D2-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "oracleD2")

	cmd := exec.Command("go", "build", "-o", binPath, refSource)
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

func runProgram(target, input string) (string, error) {
	cmd := exec.Command(target)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutputs(output string, expected int) ([]int64, error) {
	tokens := strings.Fields(output)
	if len(tokens) != expected {
		return nil, fmt.Errorf("expected %d integers, got %d", expected, len(tokens))
	}
	res := make([]int64, expected)
	for i, tok := range tokens {
		val, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q: %v", tok, err)
		}
		res[i] = val
	}
	return res, nil
}

func buildTests() []testCase {
	tests := []testCase{
		newManualTest("minimal_zero_m", `1
1 0
1 0
`, 1),
		newManualTest("small_chainable", `1
3 6
2 0 1
3 0 1 2
3 1 3 4
`, 1),
		newManualTest("duplicate_mex_entries", `1
4 7
2 0 1
3 0 1 5
4 0 1 2 5
1 2
`, 1),
		newManualTest("statement_sample", `6
3 4
2 0 2
3 2 3 3
4 7 0 1 5
3 4
5 0 2 0 4 11
1 15
1 3
0 3 3
2 50
2 1 2
2 1 2
1 17
1 2
4 1 4 9 5
4 1145142
2 2 2
5 7 3 6 0 3
3 0 1 1
5 0 9 2 1 5
5 19198101
22 324003 0 3 1416324 2 14607284 1312631 2 0 14151955 1223554 192248 2 1492515 725556
`, 6),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 250; i++ {
		tests = append(tests, randomTest(rng, i))
	}
	return tests
}

func newManualTest(name, input string, outputs int) testCase {
	return testCase{name: name, input: input, outputs: outputs}
}

func randomTest(rng *rand.Rand, idx int) testCase {
	t := rng.Intn(4) + 1
	var sb strings.Builder
	sb.Grow(2048)
	sb.WriteString(strconv.Itoa(t))
	sb.WriteByte('\n')

	for caseIdx := 0; caseIdx < t; caseIdx++ {
		n := rng.Intn(8) + 1
		m := rng.Int63n(1_000_000_000)
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))

		for i := 0; i < n; i++ {
			length := rng.Intn(10) + 1
			sb.WriteString(strconv.Itoa(length))
			for j := 0; j < length; j++ {
				val := drawValue(rng)
				sb.WriteByte(' ')
				sb.WriteString(strconv.FormatInt(val, 10))
			}
			sb.WriteByte('\n')
		}
	}

	return testCase{
		name:    fmt.Sprintf("random_%d", idx+1),
		input:   sb.String(),
		outputs: t,
	}
}

func drawValue(rng *rand.Rand) int64 {
	if rng.Intn(5) == 0 {
		return rng.Int63n(1_000_000_000)
	}
	return int64(rng.Intn(2000))
}

func buildGoFile(src string) (string, func(), error) {
	abs, err := filepath.Abs(src)
	if err != nil {
		return "", nil, err
	}
	tmpDir, err := os.MkdirTemp("", "candidate-2003D2-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "candidateD2")

	cmd := exec.Command("go", "build", "-o", binPath, abs)
	cmd.Dir = filepath.Dir(abs)
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
