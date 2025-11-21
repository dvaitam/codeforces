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
	m int
	k int
	w []int
	a []int
	b []int
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-1240F-")
	if err != nil {
		return "", nil, err
	}
	path := filepath.Join(tmpDir, "oracleF")
	cmd := exec.Command("go", "build", "-o", path, "1240F.go")
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

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.m, tc.k))
	for i, v := range tc.w {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for i := 0; i < tc.m; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.a[i]+1, tc.b[i]+1))
	}
	return sb.String()
}

func parseOutput(out string, m, k int) ([]int, error) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != m {
		return nil, fmt.Errorf("expected %d lines, got %d", m, len(lines))
	}
	res := make([]int, m)
	for i := 0; i < m; i++ {
		v, err := strconv.Atoi(strings.TrimSpace(lines[i]))
		if err != nil {
			return nil, fmt.Errorf("invalid integer in line %d: %v", i+1, err)
		}
		if v < 0 || v > k {
			return nil, fmt.Errorf("line %d: stadium index %d out of range [0,%d]", i+1, v, k)
		}
		res[i] = v
	}
	return res, nil
}

func evaluate(tc testCase, assign []int) (int, error) {
	if len(assign) != tc.m {
		return 0, fmt.Errorf("assignment length mismatch")
	}
	games := make([][]int, tc.n)
	for i := 0; i < tc.n; i++ {
		games[i] = make([]int, tc.k+1)
	}
	total := 0
	for i := 0; i < tc.m; i++ {
		t := assign[i]
		if t == 0 {
			continue
		}
		if t < 1 || t > tc.k {
			return 0, fmt.Errorf("game %d assigned to invalid stadium %d", i+1, t)
		}
		u := tc.a[i]
		v := tc.b[i]
		games[u][t]++
		games[v][t]++
		total += tc.w[u] + tc.w[v]
	}
	for i := 0; i < tc.n; i++ {
		mi, ma := 0, 0
		for j := 1; j <= tc.k; j++ {
			if j == 1 || games[i][j] < mi {
				mi = games[i][j]
			}
			if games[i][j] > ma {
				ma = games[i][j]
			}
		}
		if ma-mi > 2 {
			return 0, fmt.Errorf("team %d violates balance: min=%d max=%d", i+1, mi, ma)
		}
	}
	return total, nil
}

func deterministicTests() []testCase {
	return []testCase{
		{
			n: 3, m: 3, k: 2,
			w: []int{5, 4, 3},
			a: []int{0, 0, 1},
			b: []int{1, 2, 2},
		},
		{
			n: 4, m: 5, k: 3,
			w: []int{1, 2, 3, 4},
			a: []int{0, 0, 1, 1, 2},
			b: []int{1, 2, 2, 3, 3},
		},
	}
}

func randomTest(rng *rand.Rand) testCase {
	n := rng.Intn(8) + 3
	m := rng.Intn(10)
	k := rng.Intn(5) + 1
	if k > n {
		k = n
	}
	w := make([]int, n)
	for i := 0; i < n; i++ {
		w[i] = rng.Intn(1000) + 1
	}
	type pair struct{ u, v int }
	used := make(map[pair]bool)
	a := make([]int, m)
	b := make([]int, m)
	for i := 0; i < m; i++ {
		for {
			u := rng.Intn(n)
			v := rng.Intn(n)
			if u == v {
				continue
			}
			if u > v {
				u, v = v, u
			}
			key := pair{u, v}
			if used[key] {
				continue
			}
			used[key] = true
			a[i] = u
			b[i] = v
			break
		}
	}
	return testCase{n: n, m: m, k: k, w: w, a: a, b: b}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
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
		expAssign, err := parseOutput(expOut, tc.m, tc.k)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid oracle output on test %d: %v\noutput:\n%s\n", idx+1, err, expOut)
			os.Exit(1)
		}
		expScore, err := evaluate(tc, expAssign)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle assignment invalid on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}

		gotOut, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		gotAssign, err := parseOutput(gotOut, tc.m, tc.k)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target output invalid on test %d: %v\noutput:\n%s\ninput:\n%s\n", idx+1, err, gotOut, input)
			os.Exit(1)
		}
		gotScore, err := evaluate(tc, gotAssign)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target assignment invalid on test %d: %v\ninput:\n%s\noutput:\n%s\n", idx+1, err, input, gotOut)
			os.Exit(1)
		}
		if gotScore != expScore {
			fmt.Fprintf(os.Stderr, "test %d: expected score %d got %d\ninput:\n%s\n", idx+1, expScore, gotScore, input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
