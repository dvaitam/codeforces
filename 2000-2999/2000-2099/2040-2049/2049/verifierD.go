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

const refSource = "2049D.go"

type testCase struct {
	name    string
	input   string
	outputs int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/candidate_binary_or_go_file")
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

	tmpDir, err := os.MkdirTemp("", "2049D-bin-")
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
		newManualTest("single_cell", [][][]int64{
			{{5}},
		}, 1),
		newManualTest("two_by_two", [][][]int64{
			{
				{1, 2},
				{3, 4},
			},
		}, 1),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng, i))
	}
	return tests
}

func newManualTest(name string, grids [][][]int64, kVal int64) testCase {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(grids)))
	sb.WriteByte('\n')
	for _, g := range grids {
		n := len(g)
		m := len(g[0])
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, kVal))
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if j > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(strconv.FormatInt(g[i][j], 10))
			}
			sb.WriteByte('\n')
		}
	}
	return testCase{name: name, input: sb.String(), outputs: len(grids)}
}

func randomTest(rng *rand.Rand, idx int) testCase {
	t := rng.Intn(4) + 1
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(t))
	sb.WriteByte('\n')
	for i := 0; i < t; i++ {
		n := rng.Intn(6) + 1
		m := rng.Intn(6) + 1
		k := rng.Int63n(1_000_000_001)
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
		for r := 0; r < n; r++ {
			for c := 0; c < m; c++ {
				val := drawValue(rng)
				if c > 0 {
					sb.WriteByte(' ')
				}
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
	if rng.Intn(6) == 0 {
		return rng.Int63n(1_000_000_000)
	}
	return int64(rng.Intn(2000))
}
