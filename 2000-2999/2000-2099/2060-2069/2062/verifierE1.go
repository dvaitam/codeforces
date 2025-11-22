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
)

type testCase struct {
	n     int
	w     []int
	edges [][2]int
}

type pair struct {
	first  int
	second int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE1.go /path/to/candidate")
		return
	}
	candPath := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Printf("failed to build oracle: %v\n", err)
		return
	}
	defer cleanup()

	tests := buildTests()
	input := buildInput(tests)

	expectedSets := make([]map[int]struct{}, len(tests))
	for i, tc := range tests {
		expectedSets[i] = winningSet(tc)
	}

	oracleOut, err := runBinary(oracle, input)
	if err != nil {
		fmt.Printf("oracle runtime error: %v\n", err)
		return
	}
	oracleAns, err := parseOutput(oracleOut, len(tests))
	if err != nil {
		fmt.Printf("oracle output parse error: %v\noutput:\n%s", err, oracleOut)
		return
	}
	if err := validateAnswers(oracleAns, expectedSets); err != nil {
		fmt.Printf("oracle produced invalid answer: %v\ninput used:\n%s", err, input)
		return
	}

	candOut, err := runBinary(candPath, input)
	if err != nil {
		fmt.Printf("candidate runtime error: %v\n", err)
		return
	}
	candAns, err := parseOutput(candOut, len(tests))
	if err != nil {
		fmt.Printf("candidate output parse error: %v\noutput:\n%s", err, candOut)
		return
	}
	if err := validateAnswers(candAns, expectedSets); err != nil {
		fmt.Printf("candidate produced wrong answer: %v\ninput used:\n%s", err, input)
		return
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot locate verifier directory")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2062E1-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleE1")
	cmd := exec.Command("go", "build", "-o", outPath, "2062E1.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("%v\n%s", err, string(out))
	}
	cleanup := func() { os.RemoveAll(tmpDir) }
	return outPath, cleanup, nil
}

func runBinary(path string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseOutput(out string, t int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d tokens, got %d", t, len(fields))
	}
	res := make([]int, t)
	for i, tok := range fields {
		val, err := strconv.Atoi(tok)
		if err != nil {
			return nil, fmt.Errorf("token %d: %v", i+1, err)
		}
		res[i] = val
	}
	return res, nil
}

func validateAnswers(ans []int, expected []map[int]struct{}) error {
	if len(ans) != len(expected) {
		return fmt.Errorf("answer length mismatch: got %d, need %d", len(ans), len(expected))
	}
	for i, set := range expected {
		if len(set) == 0 {
			if ans[i] != 0 {
				return fmt.Errorf("test %d: expected 0, got %d", i+1, ans[i])
			}
			continue
		}
		if _, ok := set[ans[i]]; !ok {
			return fmt.Errorf("test %d: answer %d is not a winning first move", i+1, ans[i])
		}
	}
	return nil
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.Grow(len(tests) * 256)
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for i, v := range tc.w {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for _, e := range tc.edges {
			sb.WriteString(strconv.Itoa(e[0]))
			sb.WriteByte(' ')
			sb.WriteString(strconv.Itoa(e[1]))
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func buildTests() []testCase {
	var tests []testCase
	rng := rand.New(rand.NewSource(2062))
	totalN := 0
	const limitN = 400000

	add := func(tc testCase) {
		tests = append(tests, tc)
		totalN += tc.n
	}

	// Deterministic small cases
	add(testCase{
		n:     1,
		w:     []int{1},
		edges: nil,
	})
	add(testCase{
		n:     2,
		w:     []int{1, 2},
		edges: [][2]int{{1, 2}},
	})
	add(testCase{
		n:     3,
		w:     []int{2, 4, 3},
		edges: [][2]int{{1, 2}, {2, 3}},
	})
	add(fromSample())

	// Randomized cases
	for totalN < limitN && len(tests) < 200 {
		remain := limitN - totalN
		maxN := 10000
		if remain < maxN {
			maxN = remain
		}
		n := rng.Intn(maxN) + 1

		w := make([]int, n)
		for i := 0; i < n; i++ {
			w[i] = rng.Intn(n) + 1
		}

		edges := make([][2]int, 0, n-1)
		for v := 2; v <= n; v++ {
			u := rng.Intn(v-1) + 1
			edges = append(edges, [2]int{u, v})
		}
		add(testCase{n: n, w: w, edges: edges})
	}
	return tests
}

func fromSample() testCase {
	// The sample contains five test cases; pick one with multiple answers.
	return testCase{
		n: 5,
		w: []int{2, 2, 4, 3, 1},
		edges: [][2]int{
			{1, 2},
			{1, 3},
			{2, 4},
			{5, 1},
		},
	}
}

func winningSet(tc testCase) map[int]struct{} {
	g := make([][]int, tc.n)
	for _, e := range tc.edges {
		u := e[0] - 1
		v := e[1] - 1
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}

	tin := make([]int, tc.n)
	tout := make([]int, tc.n)
	order := make([]int, 0, tc.n)

	type item struct {
		v         int
		parent    int
		processed bool
	}
	st := []item{{0, -1, false}}
	for len(st) > 0 {
		cur := st[len(st)-1]
		st = st[:len(st)-1]
		if cur.processed {
			tout[cur.v] = len(order)
			continue
		}
		tin[cur.v] = len(order)
		order = append(order, cur.v)
		st = append(st, item{cur.v, cur.parent, true})
		for _, to := range g[cur.v] {
			if to == cur.parent {
				continue
			}
			st = append(st, item{to, cur.v, false})
		}
	}

	vals := make([]int, tc.n)
	for idx, v := range order {
		vals[idx] = tc.w[v]
	}

	pref := make([]pair, tc.n)
	for i := 0; i < tc.n; i++ {
		if i == 0 {
			pref[i] = pair{vals[i], -1}
		} else {
			pref[i] = merge(pref[i-1], pair{vals[i], -1})
		}
	}

	suff := make([]pair, tc.n)
	for i := tc.n - 1; i >= 0; i-- {
		if i == tc.n-1 {
			suff[i] = pair{vals[i], -1}
		} else {
			suff[i] = merge(suff[i+1], pair{vals[i], -1})
		}
	}

	res := make(map[int]struct{})
	for v := 0; v < tc.n; v++ {
		l, r := tin[v], tout[v]
		outside := pair{-1, -1}
		if l > 0 {
			outside = merge(outside, pref[l-1])
		}
		if r < tc.n {
			outside = merge(outside, suff[r])
		}
		top1, top2 := outside.first, outside.second
		if top1 == -1 {
			continue
		}
		if top1 > tc.w[v] && (top2 == -1 || top2 <= tc.w[v]) {
			res[v+1] = struct{}{}
		}
	}
	return res
}

func merge(a, b pair) pair {
	x1, x2 := a.first, a.second
	candidates := []int{b.first, b.second}
	for _, v := range candidates {
		if v == -1 || v == x1 || v == x2 {
			continue
		}
		if v > x1 {
			x2 = x1
			x1 = v
		} else if v > x2 {
			x2 = v
		}
	}
	return pair{x1, x2}
}
