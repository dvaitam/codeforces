package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	refSource = "1000-1999/1900-1999/1970-1979/1970/1970G2.go"
	maxNTotal = 300
	maxMTotal = 300
)

type testCase struct {
	name  string
	input string
}

type graphCase struct {
	n     int
	edges [][2]int
	c     int64
}

type edge struct {
	to  int
	idx int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG2.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()
	for idx, tc := range tests {
		expected, err := solveInput(tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "internal solver error on test %d (%s): %v\n", idx+1, tc.name, err)
			os.Exit(1)
		}
		t := len(expected)

		if _, err := runAndParse(refBin, tc.input, t); err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime/format error on test %d (%s): %v\n", idx+1, tc.name, err)
			os.Exit(1)
		}

		got, err := runAndParse(candidate, tc.input, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime/format error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		for i := range expected {
			if got[i] != expected[i] {
				fmt.Fprintf(os.Stderr, "test %d (%s) failed on case %d: expected %d got %d\ninput:\n%s", idx+1, tc.name, i+1, expected[i], got[i], tc.input)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-1970G2-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref1970G2.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		_ = os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func runAndParse(bin, input string, expectedCount int) ([]int64, error) {
	raw, err := runProgram(bin, input)
	if err != nil {
		return nil, fmt.Errorf("%v\nprogram output:\n%s", err, raw)
	}
	return parseOutputs(raw, expectedCount)
}

func parseOutputs(output string, expectedCount int) ([]int64, error) {
	fields := strings.Fields(output)
	if len(fields) != expectedCount {
		return nil, fmt.Errorf("expected %d integers, got %d (output: %q)", expectedCount, len(fields), output)
	}
	res := make([]int64, expectedCount)
	for i, token := range fields {
		val, err := strconv.ParseInt(token, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", token)
		}
		res[i] = val
	}
	return res, nil
}

func solveInput(input string) ([]int64, error) {
	reader := strings.NewReader(input)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, fmt.Errorf("failed to read t: %v", err)
	}
	results := make([]int64, t)
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		var n, m int
		var c int64
		if _, err := fmt.Fscan(reader, &n, &m, &c); err != nil {
			return nil, fmt.Errorf("failed to read case %d header: %v", caseIdx+1, err)
		}
		edges := make([][2]int, m)
		for i := 0; i < m; i++ {
			var u, v int
			if _, err := fmt.Fscan(reader, &u, &v); err != nil {
				return nil, fmt.Errorf("failed to read edge %d of case %d: %v", i+1, caseIdx+1, err)
			}
			edges[i] = [2]int{u - 1, v - 1}
		}
		results[caseIdx] = computeMinFunding(n, edges, c)
	}
	return results, nil
}

func computeMinFunding(n int, edges [][2]int, c int64) int64 {
	adj := make([][]edge, n)
	for idx, e := range edges {
		u := e[0]
		v := e[1]
		adj[u] = append(adj[u], edge{to: v, idx: idx})
		adj[v] = append(adj[v], edge{to: u, idx: idx})
	}

	compID, compSizes, k := connectedComponentInfo(n, adj)
	const inf int64 = 1 << 62
	best := inf

	otherDP := make([][]bool, k)
	for compIdx := 0; compIdx < k; compIdx++ {
		otherTotal := n - compSizes[compIdx]
		dp := make([]bool, otherTotal+1)
		dp[0] = true
		for j := 0; j < k; j++ {
			if j == compIdx {
				continue
			}
			sz := compSizes[j]
			for s := otherTotal; s >= sz; s-- {
				if dp[s-sz] {
					dp[s] = true
				}
			}
		}
		otherDP[compIdx] = dp
	}

	bridgeExtraCost := int64(k-1) * c

	if k >= 2 {
		dp := make([]bool, n+1)
		dp[0] = true
		for _, sz := range compSizes {
			for s := n; s >= sz; s-- {
				if dp[s-sz] {
					dp[s] = true
				}
			}
		}
		baseCost := bridgeExtraCost
		for s := 1; s <= n-1; s++ {
			if dp[s] {
				ss := int64(s)
				tt := int64(n - s)
				cost := ss*ss + tt*tt + baseCost
				if cost < best {
					best = cost
				}
			}
		}
	}

	tin := make([]int, n)
	low := make([]int, n)
	sub := make([]int, n)
	for i := range tin {
		tin[i] = -1
		low[i] = -1
	}
	timer := 0

	var dfs func(v int, parentEdge int)
	dfs = func(v int, parentEdge int) {
		tin[v] = timer
		low[v] = timer
		timer++
		sub[v] = 1
		for _, e := range adj[v] {
			if e.idx == parentEdge {
				continue
			}
			u := e.to
			if tin[u] != -1 {
				if tin[u] < low[v] {
					low[v] = tin[u]
				}
				continue
			}
			dfs(u, e.idx)
			sub[v] += sub[u]
			if low[u] < low[v] {
				low[v] = low[u]
			}
			if low[u] > tin[v] {
				s := sub[u]
				compIdx := compID[v]
				otherTotal := len(otherDP[compIdx]) - 1
				for extra := 0; extra <= otherTotal; extra++ {
					if !otherDP[compIdx][extra] {
						continue
					}
					ss := int64(s + extra)
					tt := int64(n) - ss
					cost := ss*ss + tt*tt + bridgeExtraCost
					if cost < best {
						best = cost
					}
				}
			}
		}
	}

	for i := 0; i < n; i++ {
		if tin[i] == -1 {
			dfs(i, -1)
		}
	}

	if best == inf {
		return -1
	}
	return best
}

func connectedComponentInfo(n int, adj [][]edge) ([]int, []int, int) {
	compID := make([]int, n)
	for i := range compID {
		compID[i] = -1
	}
	var sizes []int
	comp := 0
	queue := make([]int, 0, n)
	for i := 0; i < n; i++ {
		if compID[i] != -1 {
			continue
		}
		queue = queue[:0]
		comp++
		compID[i] = comp - 1
		size := 0
		queue = append(queue, i)
		for len(queue) > 0 {
			v := queue[0]
			queue = queue[1:]
			size++
			for _, e := range adj[v] {
				if compID[e.to] == -1 {
					compID[e.to] = comp - 1
					queue = append(queue, e.to)
				}
			}
		}
		sizes = append(sizes, size)
	}
	return compID, sizes, comp
}

func buildTests() []testCase {
	tests := []testCase{
		testManualCycle(),
		testManualMixed(),
		testManualComponents(),
		testManualDense(),
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 40; i++ {
		tests = append(tests, randomTestCase(rng, i+1))
	}
	return tests
}

func testManualCycle() testCase {
	cases := []graphCase{
		{
			n: 3,
			edges: [][2]int{
				{1, 2}, {2, 3}, {3, 1},
			},
			c: 5,
		},
	}
	return testCase{name: "manual_cycle", input: buildInput(cases)}
}

func testManualMixed() testCase {
	cases := []graphCase{
		{
			n: 4,
			edges: [][2]int{
				{1, 2}, {2, 3}, {3, 4},
			},
			c: 2,
		},
		{
			n: 5,
			edges: [][2]int{
				{1, 2}, {3, 4},
			},
			c: 3,
		},
	}
	return testCase{name: "manual_mix", input: buildInput(cases)}
}

func testManualComponents() testCase {
	cases := []graphCase{
		{
			n: 6,
			edges: [][2]int{
				{1, 2}, {3, 4}, {5, 6},
			},
			c: 10,
		},
		{
			n: 6,
			edges: [][2]int{
				{1, 2}, {2, 3}, {3, 4}, {4, 1}, {4, 5},
			},
			c: 7,
		},
	}
	return testCase{name: "manual_components", input: buildInput(cases)}
}

func testManualDense() testCase {
	cases := []graphCase{
		{
			n: 4,
			edges: [][2]int{
				{1, 2}, {2, 3}, {3, 4}, {4, 1}, {1, 3},
			},
			c: 4,
		},
	}
	return testCase{name: "manual_dense", input: buildInput(cases)}
}

func buildInput(cases []graphCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, cs := range cases {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", cs.n, len(cs.edges), cs.c))
		for _, e := range cs.edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
	}
	return sb.String()
}

func randomTestCase(rng *rand.Rand, idx int) testCase {
	leftN := maxNTotal
	leftM := maxMTotal
	numCases := rng.Intn(4) + 1
	cases := make([]graphCase, 0, numCases)
	for i := 0; i < numCases; i++ {
		casesLeft := numCases - i
		minN := 2
		maxN := leftN - 2*(casesLeft-1)
		if maxN < minN {
			maxN = minN
		}
		n := minN + rng.Intn(maxN-minN+1)

		maxEdgesPossible := n * (n - 1) / 2
		minM := 1
		maxMAllowed := leftM - (casesLeft - 1)
		if maxMAllowed > maxEdgesPossible {
			maxMAllowed = maxEdgesPossible
		}
		if maxMAllowed < minM {
			maxMAllowed = minM
		}
		m := minM + rng.Intn(maxMAllowed-minM+1)

		edges := randomEdges(rng, n, m)
		c := rng.Int63n(1_000_000_000) + 1
		cases = append(cases, graphCase{n: n, edges: edges, c: c})
		leftN -= n
		leftM -= m
	}
	name := fmt.Sprintf("random_%d", idx)
	return testCase{name: name, input: buildInput(cases)}
}

func randomEdges(rng *rand.Rand, n, m int) [][2]int {
	set := make(map[int]struct{})
	res := make([][2]int, 0, m)
	for len(res) < m {
		u := rng.Intn(n)
		v := rng.Intn(n - 1)
		if v >= u {
			v++
		}
		a, b := u, v
		if a > b {
			a, b = b, a
		}
		key := a*n + b
		if _, ok := set[key]; ok {
			continue
		}
		set[key] = struct{}{}
		res = append(res, [2]int{a + 1, b + 1})
	}
	return res
}
