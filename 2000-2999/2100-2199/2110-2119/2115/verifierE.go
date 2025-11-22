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

type edge struct {
	u int
	v int
}

type query struct {
	p int
	r int64
}

type testCase struct {
	n int
	m int
	c []int
	w []int64
	e []edge
	q []query
}

const (
	maxCoordR  = int64(1_000_000_000)
	randSeed   = 2115
	maxTests   = 80
	usageGuide = "usage: go run verifierE.go /path/to/candidate"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println(usageGuide)
		return
	}
	candidate := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Printf("failed to build oracle: %v\n", err)
		return
	}
	defer cleanup()

	tests := buildTests()
	for idx, tc := range tests {
		input := buildInput(tc)

		oracleOut, err := runBinary(oracle, input)
		if err != nil {
			fmt.Printf("oracle runtime error on test %d: %v\ninput:\n%s", idx+1, err, input)
			return
		}
		expected, err := parseOutput(oracleOut, len(tc.q))
		if err != nil {
			fmt.Printf("oracle output parse error on test %d: %v\noutput:\n%s", idx+1, err, oracleOut)
			return
		}

		candOut, err := runBinary(candidate, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\ninput:\n%s", idx+1, err, input)
			return
		}
		got, err := parseOutput(candOut, len(tc.q))
		if err != nil {
			fmt.Printf("candidate output parse error on test %d: %v\noutput:\n%s", idx+1, err, candOut)
			return
		}

		for i := range expected {
			if expected[i] != got[i] {
				fmt.Printf("Mismatch on test %d, query %d: expected %d, got %d\n", idx+1, i+1, expected[i], got[i])
				fmt.Println("Input used:")
				fmt.Print(input)
				return
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot locate verifier directory")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2115E-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracle2115E")
	cmd := exec.Command("go", "build", "-o", outPath, "2115E.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("%v\n%s", err, string(out))
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
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

func parseOutput(out string, q int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != q {
		return nil, fmt.Errorf("expected %d outputs, got %d", q, len(fields))
	}
	res := make([]int64, q)
	for i, tok := range fields {
		val, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("token %d: %v", i+1, err)
		}
		res[i] = val
	}
	return res, nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.Grow(tc.n*24 + tc.m*12 + len(tc.q)*16)
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte(' ')
	sb.WriteString(strconv.Itoa(tc.m))
	sb.WriteByte('\n')
	for i := 0; i < tc.n; i++ {
		sb.WriteString(strconv.Itoa(tc.c[i]))
		sb.WriteByte(' ')
		sb.WriteString(strconv.FormatInt(tc.w[i], 10))
		sb.WriteByte('\n')
	}
	for _, e := range tc.e {
		sb.WriteString(strconv.Itoa(e.u))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(e.v))
		sb.WriteByte('\n')
	}
	sb.WriteString(strconv.Itoa(len(tc.q)))
	sb.WriteByte('\n')
	for _, qq := range tc.q {
		sb.WriteString(strconv.Itoa(qq.p))
		sb.WriteByte(' ')
		sb.WriteString(strconv.FormatInt(qq.r, 10))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func buildTests() []testCase {
	var tests []testCase
	rng := rand.New(rand.NewSource(randSeed))

	add := func(tc testCase) {
		tests = append(tests, tc)
	}

	// Small deterministic cases
	add(testCase{
		n: 1,
		m: 0,
		c: []int{1},
		w: []int64{5},
		e: nil,
		q: []query{
			{p: 1, r: 1},
			{p: 1, r: 5},
			{p: 1, r: 10},
		},
	})

	add(simpleChain(3, []int{3, 2, 1}, []int64{9, 5, 2}, []int64{1, 4, 10}))
	add(simpleChain(5, []int{4, 2, 3, 1, 2}, []int64{2, 8, 5, 6, 3}, []int64{10, 20, 30, 40, 50}))

	// Random DAG cases
	for len(tests) < maxTests {
		n := rng.Intn(200) + 1
		m := rng.Intn(min(2000, n*(n-1)/2)-(n-1)+1) + (n - 1)
		tc := randomCase(n, m, rng)
		add(tc)
	}

	return tests
}

func simpleChain(n int, costs []int, weights []int64, coins []int64) testCase {
	e := make([]edge, 0, n-1)
	for i := 1; i < n; i++ {
		e = append(e, edge{u: i, v: i + 1})
	}
	q := make([]query, len(coins))
	for i, r := range coins {
		q[i] = query{p: n, r: r}
	}
	return testCase{
		n: n,
		m: len(e),
		c: costs,
		w: weights,
		e: e,
		q: q,
	}
}

func randomCase(n, m int, rng *rand.Rand) testCase {
	c := make([]int, n)
	w := make([]int64, n)
	for i := 0; i < n; i++ {
		c[i] = rng.Intn(200) + 1
		// spread weights to cover both small and big values
		if rng.Intn(5) == 0 {
			w[i] = int64(rng.Intn(1_000_000_000) + 1)
		} else {
			w[i] = int64(rng.Intn(1000) + 1)
		}
	}

	edges := make([]edge, 0, m)
	// ensure reachability via chain
	for i := 1; i < n; i++ {
		edges = append(edges, edge{u: i, v: i + 1})
	}
	edgeSet := make(map[int]struct{}, m)
	for _, e := range edges {
		edgeSet[(e.u-1)*n+(e.v-1)] = struct{}{}
	}

	// add remaining random edges
	for len(edges) < m {
		u := rng.Intn(n-1) + 1
		v := rng.Intn(n-u) + u + 1 // ensure u < v
		key := (u-1)*n + (v - 1)
		if _, ok := edgeSet[key]; ok {
			continue
		}
		edgeSet[key] = struct{}{}
		edges = append(edges, edge{u: u, v: v})
	}

	// build queries
	qCount := rng.Intn(300) + 1
	q := make([]query, qCount)
	for i := 0; i < qCount; i++ {
		p := rng.Intn(n) + 1
		var r int64
		switch rng.Intn(4) {
		case 0:
			r = int64(rng.Intn(50) + 1)
		case 1:
			r = int64(rng.Intn(50000) + 1)
		case 2:
			r = int64(rng.Intn(40000) + 1 + 100000)
		default:
			r = maxCoordR
		}
		q[i] = query{p: p, r: r}
	}

	return testCase{
		n: n,
		m: m,
		c: c,
		w: w,
		e: edges,
		q: q,
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
