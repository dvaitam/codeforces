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

type edge struct {
	h int
	w int
	r int
}

type testCase struct {
	n     int
	edges []edge
	t     int64
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-177F1-")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(tmpDir, "oracleF1")
	cmd := exec.Command("go", "build", "-o", bin, "177F1.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return bin, cleanup, nil
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
	return strings.TrimSpace(stdout.String()), nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, len(tc.edges), tc.t))
	for _, e := range tc.edges {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", e.h, e.w, e.r))
	}
	return sb.String()
}

func parseAnswer(out string) (int, error) {
	fields := strings.Fields(out)
	if len(fields) != 1 {
		return 0, fmt.Errorf("expected 1 token, got %d: %q", len(fields), out)
	}
	val, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q: %v", fields[0], err)
	}
	return val, nil
}

func countMatchings(tc testCase) int64 {
	k := len(tc.edges)
	masksMen := make([]int, k)
	masksWomen := make([]int, k)
	for i, e := range tc.edges {
		masksMen[i] = 1 << (e.h - 1)
		masksWomen[i] = 1 << (e.w - 1)
	}
	var dfs func(idx, menMask, womMask int) int64
	dfs = func(idx, menMask, womMask int) int64 {
		if idx == k {
			return 1
		}
		res := dfs(idx+1, menMask, womMask)
		if (menMask&masksMen[idx]) == 0 && (womMask&masksWomen[idx]) == 0 {
			res += dfs(idx+1, menMask|masksMen[idx], womMask|masksWomen[idx])
		}
		return res
	}
	return dfs(0, 0, 0)
}

func deterministicTests() []testCase {
	return []testCase{
		{
			n: 1,
			edges: []edge{
				{h: 1, w: 1, r: 5},
			},
			t: 1, // empty set
		},
		{
			n: 2,
			edges: []edge{
				{h: 1, w: 1, r: 3},
				{h: 2, w: 2, r: 7},
				{h: 1, w: 2, r: 4},
			},
			t: 4,
		},
		{
			n: 5,
			edges: []edge{
				{1, 1, 10},
				{2, 2, 20},
				{3, 3, 30},
				{4, 4, 40},
				{5, 5, 50},
			},
			t: 17,
		},
		// stress t near upper constraint while keeping edges disjoint
		{
			n: 20,
			edges: []edge{
				{1, 1, 1}, {2, 2, 2}, {3, 3, 3}, {4, 4, 4}, {5, 5, 5},
				{6, 6, 6}, {7, 7, 7}, {8, 8, 8}, {9, 9, 9}, {10, 10, 10},
				{11, 11, 11}, {12, 12, 12}, {13, 13, 13}, {14, 14, 14},
				{15, 15, 15}, {16, 16, 16}, {17, 17, 17}, {18, 18, 18},
			},
			t: 200000,
		},
	}
}

func randomTest(rng *rand.Rand) testCase {
	n := rng.Intn(6) + 1
	maxEdges := n * n
	k := rng.Intn(minInt(maxEdges, 12)-1) + 1 // ensure at least 1 edge
	used := make(map[[2]int]bool)
	edges := make([]edge, 0, k)
	for len(edges) < k {
		h := rng.Intn(n) + 1
		w := rng.Intn(n) + 1
		if used[[2]int{h, w}] {
			continue
		}
		used[[2]int{h, w}] = true
		r := rng.Intn(1000) + 1
		edges = append(edges, edge{h: h, w: w, r: r})
	}
	tc := testCase{n: n, edges: edges}
	total := countMatchings(tc)
	if total <= 0 {
		total = 1
	}
	if total > 200000 {
		total = 200000
	}
	tc.t = rng.Int63n(total) + 1
	return tc
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF1.go /path/to/binary")
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
		expStr, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		exp, err := parseAnswer(expStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle output invalid on test %d: %v\n%s\n", idx+1, err, expStr)
			os.Exit(1)
		}

		gotStr, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		got, err := parseAnswer(gotStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target output invalid on test %d: %v\noutput:\n%s\ninput:\n%s\n", idx+1, err, gotStr, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "mismatch on test %d: expected %d got %d\ninput:\n%s\n", idx+1, exp, got, input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
