package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	n int
	m int
	a []int
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-1045B-")
	if err != nil {
		return "", nil, err
	}
	path := filepath.Join(tmpDir, "oracleB")
	cmd := exec.Command("go", "build", "-o", path, "1045B.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() { os.RemoveAll(tmpDir) }
	return path, cleanup, nil
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
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func parseOutput(out string, m int) ([]int, error) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty output")
	}
	k, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return nil, fmt.Errorf("invalid K: %v", err)
	}
	if k == 0 {
		return []int{}, nil
	}
	if len(lines) < 2 {
		return nil, fmt.Errorf("missing second line for residues")
	}
	fields := strings.Fields(lines[1])
	if len(fields) != k {
		return nil, fmt.Errorf("expected %d residues, got %d", k, len(fields))
	}
	res := make([]int, k)
	for i, f := range fields {
		val, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("invalid residue %q: %v", f, err)
		}
		if val < 0 || val >= m {
			return nil, fmt.Errorf("residue %d out of range", val)
		}
		res[i] = val
	}
	if !sort.IntsAreSorted(res) {
		return nil, fmt.Errorf("residues not sorted")
	}
	for i := 1; i < len(res); i++ {
		if res[i] == res[i-1] {
			return nil, fmt.Errorf("duplicate residue %d", res[i])
		}
	}
	return res, nil
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 2, m: 5, a: []int{3, 4}},
		{n: 4, m: 1_000_000_000, a: []int{5, 25, 125, 625}},
		{n: 1, m: 10, a: []int{7}},
		{n: 3, m: 10, a: []int{0, 2, 7}},
	}
}

func randomTest(rng *rand.Rand) testCase {
	n := rng.Intn(15) + 1
	if rng.Intn(4) == 0 {
		n = rng.Intn(2000) + 1
	}
	m := n + 1 + rng.Intn(1000)
	if rng.Intn(4) == 0 {
		m = n + 1 + rng.Intn(1_000_000_000-n)
	}
	set := make(map[int]struct{})
	for len(set) < n {
		val := rng.Intn(m)
		set[val] = struct{}{}
	}
	a := make([]int, 0, n)
	for v := range set {
		a = append(a, v)
	}
	sort.Ints(a)
	return testCase{n: n, m: m, a: a}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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
		expRes, err := parseOutput(expOut, tc.m)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid oracle output on test %d: %v\noutput:\n%s\n", idx+1, err, expOut)
			os.Exit(1)
		}

		gotOut, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		gotRes, err := parseOutput(gotOut, tc.m)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target output invalid on test %d: %v\noutput:\n%s\ninput:\n%s\n", idx+1, err, gotOut, input)
			os.Exit(1)
		}

		if len(expRes) != len(gotRes) {
			fmt.Fprintf(os.Stderr, "test %d: mismatch in count: expected %d got %d\ninput:\n%s\n", idx+1, len(expRes), len(gotRes), input)
			os.Exit(1)
		}
		for i := range expRes {
			if expRes[i] != gotRes[i] {
				fmt.Fprintf(os.Stderr, "test %d: mismatch at position %d: expected %d got %d\ninput:\n%s\n", idx+1, i+1, expRes[i], gotRes[i], input)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
