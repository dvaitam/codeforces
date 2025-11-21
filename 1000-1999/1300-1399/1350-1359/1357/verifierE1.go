package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const tol = 1e-8

type testCase struct {
	P  int
	n  int
	re []float64
	im []float64
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-1357E1-")
	if err != nil {
		return "", nil, err
	}
	path := filepath.Join(tmpDir, "oracleE1")
	cmd := exec.Command("go", "build", "-o", path, "1357E1.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() { os.RemoveAll(tmpDir) }
	return path, cleanup, nil
}

func runBinary(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseOutput(out string, size int) ([]float64, []float64, error) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != size {
		return nil, nil, fmt.Errorf("expected %d lines, got %d", size, len(lines))
	}
	re := make([]float64, size)
	im := make([]float64, size)
	for i := 0; i < size; i++ {
		fields := strings.Fields(lines[i])
		if len(fields) != 2 {
			return nil, nil, fmt.Errorf("line %d: expected two floats", i+1)
		}
		valRe, err := strconv.ParseFloat(fields[0], 64)
		if err != nil {
			return nil, nil, fmt.Errorf("line %d: invalid real part: %v", i+1, err)
		}
		valIm, err := strconv.ParseFloat(fields[1], 64)
		if err != nil {
			return nil, nil, fmt.Errorf("line %d: invalid imaginary part: %v", i+1, err)
		}
		re[i] = valRe
		im[i] = valIm
	}
	return re, im, nil
}

func compareComplex(expRe, expIm, gotRe, gotIm []float64) error {
	if len(expRe) != len(gotRe) {
		return fmt.Errorf("length mismatch")
	}
	for i := range expRe {
		if math.Abs(expRe[i]-gotRe[i]) > tol || math.Abs(expIm[i]-gotIm[i]) > tol {
			return fmt.Errorf("index %d mismatch: expected (%.10f, %.10f) got (%.10f, %.10f)", i, expRe[i], expIm[i], gotRe[i], gotIm[i])
		}
	}
	return nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	size := 1 << tc.n
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.P, tc.n))
	for i := 0; i < size; i++ {
		sb.WriteString(fmt.Sprintf("%.10f %.10f\n", tc.re[i], tc.im[i]))
	}
	return sb.String()
}

func deterministicTests() []testCase {
	return []testCase{
		{
			P: 1, n: 1,
			re: []float64{1, 0},
			im: []float64{0, 0},
		},
		{
			P: 2, n: 2,
			re: []float64{1, 0, 0, 1},
			im: []float64{0, 0, 1, 0},
		},
	}
}

func randomState(rng *rand.Rand, size int) ([]float64, []float64) {
	re := make([]float64, size)
	im := make([]float64, size)
	for i := 0; i < size; i++ {
		re[i] = rng.Float64()*2 - 1
		im[i] = rng.Float64()*2 - 1
	}
	// normalize to unit norm
	var norm float64
	for i := 0; i < size; i++ {
		norm += re[i]*re[i] + im[i]*im[i]
	}
	norm = math.Sqrt(norm)
	if norm == 0 {
		norm = 1
	}
	for i := 0; i < size; i++ {
		re[i] /= norm
		im[i] /= norm
	}
	return re, im
}

func randomTest(rng *rand.Rand) testCase {
	P := rng.Intn(1000) + 2
	n := rng.Intn(4) + 1
	size := 1 << n
	re, im := randomState(rng, size)
	return testCase{P: P, n: n, re: re, im: im}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE1.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng))
	}

	for idx, tc := range tests {
		input := buildInput(tc)
		size := 1 << tc.n

		expOut, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		expRe, expIm, err := parseOutput(expOut, size)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid oracle output on test %d: %v\noutput:\n%s\n", idx+1, err, expOut)
			os.Exit(1)
		}

		gotOut, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		gotRe, gotIm, err := parseOutput(gotOut, size)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target output invalid on test %d: %v\noutput:\n%s\ninput:\n%s\n", idx+1, err, gotOut, input)
			os.Exit(1)
		}

		if err := compareComplex(expRe, expIm, gotRe, gotIm); err != nil {
			fmt.Fprintf(os.Stderr, "mismatch on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
