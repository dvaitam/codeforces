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

type interval struct {
	l int
	r int
}

type testCase struct {
	planets []interval
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
	tmpDir, err := os.MkdirTemp("", "oracle-2030G2-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleG2")
	cmd := exec.Command("go", "build", "-o", outPath, "2030G2.go")
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

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.Grow(len(tc.planets)*16 + 32)
	sb.WriteString("1\n") // single test case
	sb.WriteString(strconv.Itoa(len(tc.planets)))
	sb.WriteByte('\n')
	for _, p := range tc.planets {
		sb.WriteString(strconv.Itoa(p.l))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(p.r))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseAnswer(out string) (int64, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	val, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q: %v", fields[0], err)
	}
	return val, nil
}

func deterministicTests() []testCase {
	return []testCase{
		{planets: []interval{{l: 1, r: 1}}},
		{planets: []interval{{l: 1, r: 1}, {l: 2, r: 3}, {l: 3, r: 3}}},
		{planets: []interval{{l: 1, r: 4}, {l: 2, r: 3}, {l: 2, r: 4}, {l: 1, r: 5}, {l: 1, r: 5}}},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 200)
	for len(tests) < cap(tests) {
		n := rng.Intn(50) + 1
		if rng.Intn(5) == 0 {
			n = rng.Intn(300) + 50
		}
		planets := make([]interval, n)
		for i := 0; i < n; i++ {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			planets[i] = interval{l: l, r: r}
		}
		tests = append(tests, testCase{planets: planets})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG2.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := append(deterministicTests(), randomTests()...)
	for idx, tc := range tests {
		input := buildInput(tc)
		expOut, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		gotOut, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		expAns, err := parseAnswer(expOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle output invalid on test %d: %v\n%s", idx+1, err, expOut)
			os.Exit(1)
		}
		gotAns, err := parseAnswer(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target output invalid on test %d: %v\n%s", idx+1, err, gotOut)
			os.Exit(1)
		}
		if expAns != gotAns {
			fmt.Fprintf(os.Stderr, "test %d mismatch: expected %d, got %d\ninput:\n%s", idx+1, expAns, gotAns, input)
			os.Exit(1)
		}
	}

	fmt.Println("All tests passed.")
}
