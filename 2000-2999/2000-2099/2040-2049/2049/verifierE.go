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
	n int64
	p int64
	k int64
}

func callerFile() (string, bool) {
	_, file, _, ok := runtime.Caller(0)
	return file, ok
}

func buildOracle() (string, func(), error) {
	_, file, ok := callerFile()
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2049E-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleE")
	cmd := exec.Command("go", "build", "-o", outPath, "2049E.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runBinary(bin, input string) (string, error) {
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

func buildInput(tc []testCase) string {
	var sb strings.Builder
	sb.Grow(len(tc) * 24)
	sb.WriteString(strconv.Itoa(len(tc)))
	sb.WriteByte('\n')
	for _, t := range tc {
		sb.WriteString(strconv.FormatInt(t.n, 10))
		sb.WriteByte(' ')
		sb.WriteString(strconv.FormatInt(t.p, 10))
		sb.WriteByte(' ')
		sb.WriteString(strconv.FormatInt(t.k, 10))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseAnswers(out string, t int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d answers, got %d", t, len(fields))
	}
	res := make([]int64, t)
	for i, f := range fields {
		val, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", f)
		}
		res[i] = val
	}
	return res, nil
}

func compareAnswers(exp, got []int64) error {
	if len(exp) != len(got) {
		return fmt.Errorf("answer count mismatch")
	}
	for i := range exp {
		if exp[i] != got[i] {
			return fmt.Errorf("test %d expected %d, got %d", i+1, exp[i], got[i])
		}
	}
	return nil
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 8, p: 6, k: 6},
		{n: 4, p: 3, k: 2},
		{n: 16, p: 1, k: 15},
	}
}

func randomPowerOfTwo(rng *rand.Rand, minPow, maxPow int) int64 {
	p := rng.Intn(maxPow-minPow+1) + minPow
	return 1 << p
}

func randomTests() [][]testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	batches := make([][]testCase, 0, 60)
	for len(batches) < cap(batches) {
		t := rng.Intn(8) + 1
		if rng.Intn(5) == 0 {
			t = rng.Intn(60) + 20
		}
		tc := make([]testCase, t)
		for i := 0; i < t; i++ {
			n := randomPowerOfTwo(rng, 2, 12) // n in [4, 4096]
			if rng.Intn(6) == 0 {
				n = randomPowerOfTwo(rng, 2, 20) // larger occasionally
			}
			p := rng.Int63n(n) + 1
			k := rng.Int63n(n-2) + 2
			tc[i] = testCase{n: n, p: p, k: k}
		}
		batches = append(batches, tc)
	}
	return batches
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	batches := [][]testCase{deterministicTests()}
	batches = append(batches, randomTests()...)

	for batchIdx, tcList := range batches {
		input := buildInput(tcList)
		expOut, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on batch %d: %v\ninput:\n%s", batchIdx+1, err, input)
			os.Exit(1)
		}
		gotOut, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on batch %d: %v\ninput:\n%s", batchIdx+1, err, input)
			os.Exit(1)
		}
		expAns, err := parseAnswers(expOut, len(tcList))
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle output invalid on batch %d: %v\n%s", batchIdx+1, err, expOut)
			os.Exit(1)
		}
		gotAns, err := parseAnswers(gotOut, len(tcList))
		if err != nil {
			fmt.Fprintf(os.Stderr, "target output invalid on batch %d: %v\n%s", batchIdx+1, err, gotOut)
			os.Exit(1)
		}
		if err := compareAnswers(expAns, gotAns); err != nil {
			fmt.Fprintf(os.Stderr, "batch %d mismatch: %v\ninput:\n%s", batchIdx+1, err, input)
			os.Exit(1)
		}
	}

	fmt.Println("All tests passed.")
}
