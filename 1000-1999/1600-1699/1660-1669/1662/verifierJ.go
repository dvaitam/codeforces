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
	n int
	A [][]int
	C [][]int
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-1662J-")
	if err != nil {
		return "", nil, err
	}
	path := filepath.Join(tmpDir, "oracleJ")
	cmd := exec.Command("go", "build", "-o", path, "1662J.go")
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
	return strings.TrimSpace(stdout.String()), nil
}

func parseOutput(out string) (int, error) {
	val, err := strconv.Atoi(strings.TrimSpace(out))
	if err != nil {
		return 0, fmt.Errorf("invalid integer output: %v", err)
	}
	if val < 0 {
		return 0, fmt.Errorf("negative answer")
	}
	return val, nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for i := 0; i < tc.n; i++ {
		for j := 0; j < tc.n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(tc.A[i][j]))
		}
		sb.WriteByte('\n')
	}
	for i := 0; i < tc.n; i++ {
		for j := 0; j < tc.n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(tc.C[i][j]))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func deterministicTests() []testCase {
	return []testCase{
		{
			n: 1,
			A: [][]int{{1}},
			C: [][]int{{1}},
		},
		{
			n: 2,
			A: [][]int{
				{1, 2},
				{2, 1},
			},
			C: [][]int{
				{1, 0},
				{0, 1},
			},
		},
		{
			n: 3,
			A: [][]int{
				{1, 2, 3},
				{2, 3, 1},
				{3, 1, 2},
			},
			C: [][]int{
				{1, 0, 1},
				{0, 1, 0},
				{1, 0, 1},
			},
		},
	}
}

func randomLatinSquare(n int, rng *rand.Rand) [][]int {
	base := make([][]int, n)
	for i := 0; i < n; i++ {
		base[i] = make([]int, n)
		for j := 0; j < n; j++ {
			base[i][j] = (i + j) % n
		}
	}
	perm := rng.Perm(n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			base[i][j] = perm[(base[i][j])]
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			base[i][j]++
		}
	}
	return base
}

func randomTest(rng *rand.Rand) testCase {
	n := rng.Intn(15) + 1
	A := randomLatinSquare(n, rng)
	C := make([][]int, n)
	for i := 0; i < n; i++ {
		C[i] = make([]int, n)
		for j := 0; j < n; j++ {
			if rng.Intn(2) == 0 {
				C[i][j] = 0
			} else {
				C[i][j] = 1
			}
		}
	}
	return testCase{n: n, A: A, C: C}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierJ.go /path/to/binary")
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

		expOut, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		expVal, err := parseOutput(expOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid oracle output on test %d: %v\noutput:\n%s\n", idx+1, err, expOut)
			os.Exit(1)
		}

		gotOut, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		gotVal, err := parseOutput(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target output invalid on test %d: %v\noutput:\n%s\ninput:\n%s\n", idx+1, err, gotOut, input)
			os.Exit(1)
		}

		if gotVal != expVal {
			fmt.Fprintf(os.Stderr, "mismatch on test %d: expected %d got %d\ninput:\n%s\n", idx+1, expVal, gotVal, input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
