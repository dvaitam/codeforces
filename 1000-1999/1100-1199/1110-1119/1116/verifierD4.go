package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type testCase struct {
	n int
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-1116D4-")
	if err != nil {
		return "", nil, err
	}
	path := filepath.Join(tmpDir, "oracleD4")
	cmd := exec.Command("go", "build", "-o", path, "1116D4.go")
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

func parseMatrix(out string, size int) ([]string, error) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != size {
		return nil, fmt.Errorf("expected %d rows, got %d", size, len(lines))
	}
	for i, line := range lines {
		if len(line) != size {
			return nil, fmt.Errorf("row %d has length %d, expected %d", i+1, len(line), size)
		}
		for _, ch := range line {
			if ch != 'X' && ch != '.' {
				return nil, fmt.Errorf("invalid character %q in row %d", ch, i+1)
			}
		}
	}
	return lines, nil
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 2},
		{n: 3},
		{n: 4},
		{n: 5},
	}
}

func randomTest(rng *rand.Rand) testCase {
	n := rng.Intn(4) + 2 // values between 2 and 5
	return testCase{n: n}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD4.go /path/to/binary")
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
	for i := 0; i < 20; i++ {
		tests = append(tests, randomTest(rng))
	}

	for idx, tc := range tests {
		input := fmt.Sprintf("%d\n", tc.n)
		expOut, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d (N=%d): %v\n", idx+1, tc.n, err)
			os.Exit(1)
		}
		gotOut, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d (N=%d): %v\n", idx+1, tc.n, err)
			os.Exit(1)
		}
		size := 1 << tc.n
		expMat, err := parseMatrix(expOut, size)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle output invalid on test %d: %v\noutput:\n%s\n", idx+1, err, expOut)
			os.Exit(1)
		}
		gotMat, err := parseMatrix(gotOut, size)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target output invalid on test %d: %v\noutput:\n%s\n", idx+1, err, gotOut)
			os.Exit(1)
		}
		for r := 0; r < size; r++ {
			if expMat[r] != gotMat[r] {
				fmt.Fprintf(os.Stderr, "mismatch on test %d row %d: expected %s got %s\nN=%d\n", idx+1, r+1, expMat[r], gotMat[r], tc.n)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
