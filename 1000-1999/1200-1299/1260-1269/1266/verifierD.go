package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	refSource = "1000-1999/1200-1299/1260-1269/1266/1266D.go"
	maxEdges  = 300000
)

type testInfo struct {
	n        int
	minTotal int64
	balance  []int64
}

type testCase struct {
	input string
	info  testInfo
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests, err := buildTests()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build tests:", err)
		os.Exit(1)
	}

	for i, tc := range tests {
		refOut, err := runExecutable(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := checkOutput(refOut, tc.info); err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d: %v\n", i+1, err)
			os.Exit(1)
		}

		candOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := checkOutput(candOut, tc.info); err != nil {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: %v\nInput:\n%sCandidate output:\n%s\n", i+1, err, tc.input, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1266D-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var combined bytes.Buffer
	cmd.Stdout = &combined
	cmd.Stderr = &combined
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, combined.String())
	}
	return tmp.Name(), nil
}

func runExecutable(path, input string) (string, error) {
	cmd := exec.Command(path)
	return execute(cmd, input)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return execute(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func execute(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildTests() ([]testCase, error) {
	rawTests := []string{
		"3 2\n1 2 5\n2 3 5\n",
		"3 3\n1 2 3\n2 3 2\n3 1 1\n",
		"4 4\n1 2 10\n2 3 5\n3 4 5\n4 1 5\n",
		"1 0\n",
		"5 0\n",
	}

	tests := make([]testCase, 0, len(rawTests)+8)
	for _, input := range rawTests {
		info, err := prepareTest(input)
		if err != nil {
			return nil, err
		}
		tests = append(tests, testCase{input: input, info: info})
	}

	randomConfigs := []struct {
		n, m int
		seed int64
	}{
		{2, 1, 1},
		{5, 7, 2},
		{10, 30, 3},
		{50, 200, 4},
		{100, 500, 5},
		{500, 2000, 6},
		{1000, 4000, 7},
		{2000, 6000, 8},
		{5000, 12000, time.Now().UnixNano()},
	}

	for _, cfg := range randomConfigs {
		input := randomTest(cfg.n, cfg.m, cfg.seed)
		info, err := prepareTest(input)
		if err != nil {
			return nil, err
		}
		tests = append(tests, testCase{input: input, info: info})
	}
	return tests, nil
}

func randomTest(n, m int, seed int64) string {
	if n < 1 {
		n = 1
	}
	if n == 1 {
		m = 0
	}
	r := rand.New(rand.NewSource(seed))
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < m; i++ {
		u := r.Intn(n) + 1
		v := r.Intn(n-1) + 1
		if v >= u {
			v++
		}
		w := int64(r.Intn(1_000_000_000) + 1)
		sb.WriteString(fmt.Sprintf("%d %d %d\n", u, v, w))
	}
	return sb.String()
}

func prepareTest(input string) (testInfo, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return testInfo{}, fmt.Errorf("failed to read n and m: %v", err)
	}
	balance := make([]int64, n+1)
	for i := 0; i < m; i++ {
		var u, v int
		var d int64
		if _, err := fmt.Fscan(reader, &u, &v, &d); err != nil {
			return testInfo{}, fmt.Errorf("failed to read debt %d: %v", i+1, err)
		}
		if u < 1 || u > n || v < 1 || v > n || u == v {
			return testInfo{}, fmt.Errorf("invalid edge in test input")
		}
		balance[u] -= d
		balance[v] += d
	}
	var minTotal int64
	for i := 1; i <= n; i++ {
		if balance[i] > 0 {
			minTotal += balance[i]
		}
	}
	return testInfo{
		n:        n,
		minTotal: minTotal,
		balance:  balance,
	}, nil
}

func checkOutput(output string, info testInfo) error {
	reader := bufio.NewReader(strings.NewReader(output))
	var k int
	if _, err := fmt.Fscan(reader, &k); err != nil {
		return fmt.Errorf("failed to read number of debts: %v", err)
	}
	if k < 0 || k > maxEdges {
		return fmt.Errorf("invalid number of debts %d", k)
	}

	diff := make([]int64, info.n+1)
	used := make(map[int64]struct{}, k*2)
	var total int64
	for i := 0; i < k; i++ {
		var u, v int
		var w int64
		if _, err := fmt.Fscan(reader, &u, &v, &w); err != nil {
			return fmt.Errorf("failed to read debt %d: %v", i+1, err)
		}
		if u < 1 || u > info.n || v < 1 || v > info.n {
			return fmt.Errorf("debt %d has node outside range", i+1)
		}
		if u == v {
			return fmt.Errorf("debt %d has self-loop", i+1)
		}
		if w <= 0 || w > 1_000_000_000_000_000_000 {
			return fmt.Errorf("debt %d has invalid weight %d", i+1, w)
		}
		key := int64(u)*(int64(info.n)+1) + int64(v)
		if _, ok := used[key]; ok {
			return fmt.Errorf("duplicate debt pair (%d,%d)", u, v)
		}
		used[key] = struct{}{}
		diff[u] -= w
		diff[v] += w
		total += w
	}

	for i := 1; i <= info.n; i++ {
		if diff[i] != info.balance[i] {
			return fmt.Errorf("balance mismatch at node %d: expected %d got %d", i, info.balance[i], diff[i])
		}
	}

	if total != info.minTotal {
		return fmt.Errorf("total debt mismatch: expected %d got %d", info.minTotal, total)
	}

	return nil
}
