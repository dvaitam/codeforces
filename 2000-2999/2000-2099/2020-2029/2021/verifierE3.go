package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type edgeT struct {
	u, v int
	w    int64
}

type testCase struct {
	n, m, p int
	req     []int
	edges   []edgeT
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE3.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[len(os.Args)-1]

	tests := generateTests()
	input := buildInput(tests)

	refOut := solveAll(input)

	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\n%s", err, candOut)
		os.Exit(1)
	}

	if err := compareOutputs(refOut, candOut, tests); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

// solveAll implements the correct reference solver (Kruskal reconstruction tree approach).
func solveAll(input string) string {
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024*10)
	scanner.Split(bufio.ScanWords)

	readInt := func() int {
		scanner.Scan()
		b := scanner.Bytes()
		res := 0
		for _, v := range b {
			res = res*10 + int(v-'0')
		}
		return res
	}

	scanner.Scan()
	b := scanner.Bytes()
	t := 0
	for _, v := range b {
		t = t*10 + int(v-'0')
	}

	var out bytes.Buffer
	bw := bufio.NewWriter(&out)

	for tc := 0; tc < t; tc++ {
		n := readInt()
		m := readInt()
		p := readInt()

		S := make([]int, p)
		for i := 0; i < p; i++ {
			S[i] = readInt()
		}

		type Edge struct{ u, v, w int }
		edges := make([]Edge, m)
		for i := 0; i < m; i++ {
			edges[i].u = readInt()
			edges[i].v = readInt()
			edges[i].w = readInt()
		}

		sort.Slice(edges, func(i, j int) bool {
			return edges[i].w < edges[j].w
		})

		parent := make([]int, 2*n)
		for i := 1; i < 2*n; i++ {
			parent[i] = i
		}

		find := func(i int) int {
			root := i
			for root != parent[root] {
				root = parent[root]
			}
			curr := i
			for curr != root {
				nxt := parent[curr]
				parent[curr] = root
				curr = nxt
			}
			return root
		}

		krt_parent := make([]int, 2*n)
		left := make([]int, 2*n)
		right := make([]int, 2*n)
		W := make([]int64, 2*n)

		node_cnt := n
		for _, e := range edges {
			r1 := find(e.u)
			r2 := find(e.v)
			if r1 != r2 {
				node_cnt++
				W[node_cnt] = int64(e.w)
				krt_parent[r1] = node_cnt
				krt_parent[r2] = node_cnt
				left[node_cnt] = r1
				right[node_cnt] = r2
				parent[r1] = node_cnt
				parent[r2] = node_cnt
			}
		}

		c := make([]int64, 2*n)
		for _, s := range S {
			c[s] = 1
		}

		for i := 1; i < node_cnt; i++ {
			p_node := krt_parent[i]
			c[p_node] += c[i]
		}

		weight := make([]int64, 2*n)
		for i := 1; i < node_cnt; i++ {
			p_node := krt_parent[i]
			weight[i] = (W[p_node] - W[i]) * c[i]
		}

		dp := make([]int64, 2*n)
		heavy := make([]int, 2*n)

		for i := 1; i <= node_cnt; i++ {
			if left[i] != 0 {
				dp[i] = -1

				v := left[i]
				val := dp[v] + weight[v]
				if val > dp[i] {
					dp[i] = val
					heavy[i] = v
				}

				v = right[i]
				val = dp[v] + weight[v]
				if val > dp[i] {
					dp[i] = val
					heavy[i] = v
				}
			}
		}

		is_light := make([]bool, 2*n)
		is_light[node_cnt] = true

		for i := node_cnt; i >= 1; i-- {
			if left[i] != 0 {
				h := heavy[i]

				v := left[i]
				if v == h {
					is_light[v] = false
				} else {
					is_light[v] = true
				}

				v = right[i]
				if v == h {
					is_light[v] = false
				} else {
					is_light[v] = true
				}
			}
		}

		paths := make([]int64, 0, n)
		for i := 1; i <= node_cnt; i++ {
			if is_light[i] {
				var pw int64
				if i == node_cnt {
					pw = dp[i]
				} else {
					pw = dp[i] + weight[i]
				}
				paths = append(paths, pw)
			}
		}

		sort.Slice(paths, func(i, j int) bool {
			return paths[i] > paths[j]
		})

		ans := int64(p) * W[node_cnt]
		for k := 1; k <= n; k++ {
			if k-1 < len(paths) {
				ans -= paths[k-1]
			}
			bw.WriteString(strconv.FormatInt(ans, 10))
			if k == n {
				bw.WriteByte('\n')
			} else {
				bw.WriteByte(' ')
			}
		}
	}

	bw.Flush()
	return out.String()
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return errBuf.String(), err
	}
	if errBuf.Len() > 0 {
		return errBuf.String(), fmt.Errorf("unexpected stderr output")
	}
	return out.String(), nil
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const maxTotalN = 200000
	const maxTotalM = 200000
	var tests []testCase
	totalN, totalM := 0, 0

	add := func(tc testCase) {
		if totalN+tc.n > maxTotalN || totalM+tc.m > maxTotalM {
			return
		}
		tests = append(tests, tc)
		totalN += tc.n
		totalM += tc.m
	}

	// Small deterministic cases
	add(buildPathTest(2, []edgeT{{1, 2, 5}}, []int{1}))
	add(buildPathTest(3, []edgeT{{1, 2, 1}, {2, 3, 2}}, []int{1, 3}))
	add(buildPathTest(4, []edgeT{{1, 2, 7}, {2, 3, 3}, {3, 4, 4}}, []int{2, 4}))

	// Random connected graphs
	for len(tests) < 40 && totalN < 180000 && totalM < 180000 {
		n := rng.Intn(3000) + 2
		maxExtra := n
		m := n - 1 + rng.Intn(maxExtra+1)
		p := rng.Intn(n) + 1
		req := randSample(rng, n, p)
		edges := generateConnectedGraph(rng, n, m)
		add(testCase{n: n, m: m, p: p, req: req, edges: edges})
	}

	// One larger stress case if budget allows
	if totalN+15000 <= maxTotalN && totalM+20000 <= maxTotalM {
		n := 15000
		m := 20000
		p := n/3 + 1
		req := randSample(rng, n, p)
		edges := generateConnectedGraph(rng, n, m)
		add(testCase{n: n, m: m, p: p, req: req, edges: edges})
	}

	return tests
}

func buildPathTest(n int, edges []edgeT, req []int) testCase {
	return testCase{
		n:     n,
		m:     len(edges),
		p:     len(req),
		req:   append([]int(nil), req...),
		edges: append([]edgeT(nil), edges...),
	}
}

func generateConnectedGraph(rng *rand.Rand, n, m int) []edgeT {
	edges := make([]edgeT, 0, m)
	for i := 1; i < n; i++ {
		edges = append(edges, edgeT{u: i, v: i + 1, w: rngWeight(rng)})
	}
	seen := make(map[int]struct{}, m)
	for _, e := range edges {
		key := e.u*n + e.v
		seen[key] = struct{}{}
	}
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		key := u*n + v
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		edges = append(edges, edgeT{u: u, v: v, w: rngWeight(rng)})
	}
	return edges
}

func rngWeight(rng *rand.Rand) int64 {
	return int64(rng.Intn(1_000_000_000) + 1)
}

func randSample(rng *rand.Rand, n, k int) []int {
	perm := rng.Perm(n)
	res := make([]int, k)
	for i := 0; i < k; i++ {
		res[i] = perm[i] + 1
	}
	return res
}

func buildInput(tests []testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%d %d %d\n", tc.n, tc.m, tc.p)
		for i, v := range tc.req {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", v)
		}
		b.WriteByte('\n')
		for _, e := range tc.edges {
			fmt.Fprintf(&b, "%d %d %d\n", e.u, e.v, e.w)
		}
	}
	return b.String()
}

func compareOutputs(refOut, candOut string, tests []testCase) error {
	refTokens := strings.Fields(refOut)
	candTokens := strings.Fields(candOut)

	total := 0
	for _, tc := range tests {
		total += tc.n
	}
	if len(refTokens) != total {
		return fmt.Errorf("reference produced %d tokens, expected %d", len(refTokens), total)
	}
	if len(candTokens) != total {
		return fmt.Errorf("candidate produced %d tokens, expected %d", len(candTokens), total)
	}

	for i := 0; i < total; i++ {
		refVal, err := strconv.ParseInt(refTokens[i], 10, 64)
		if err != nil {
			return fmt.Errorf("reference token %q is not integer", refTokens[i])
		}
		candVal, err := strconv.ParseInt(candTokens[i], 10, 64)
		if err != nil {
			return fmt.Errorf("candidate token %q is not integer", candTokens[i])
		}
		if refVal != candVal {
			testIdx, pos := locateToken(tests, i)
			return fmt.Errorf("mismatch at test %d, position %d: expected %d got %d", testIdx+1, pos+1, refVal, candVal)
		}
	}
	return nil
}

func locateToken(tests []testCase, idx int) (testIndex int, pos int) {
	rem := idx
	for i, tc := range tests {
		if rem < tc.n {
			return i, rem
		}
		rem -= tc.n
	}
	return len(tests) - 1, tests[len(tests)-1].n - 1
}
