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
)

const refSource = "2000-2999/2000-2099/2050-2059/2057/2057E2.go"

type edge struct {
	u, v int
	w    int
}

type query struct {
	a, b int
	k    int
}

type testCase struct {
	n       int
	edges   []edge
	queries []query
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE2.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[len(os.Args)-1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	input := renderInput(tests)

	refOut, err := runWithInput(exec.Command(refBin), input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference solution failed: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}

	candOut, err := runWithInput(commandFor(candidate), input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}

	expect, err := parseOutputs(refOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse reference output: %v\n", err)
		os.Exit(1)
	}
	got, err := parseOutputs(candOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse candidate output: %v\n", err)
		os.Exit(1)
	}

	for i := range tests {
		if len(expect[i]) != len(got[i]) {
			fmt.Fprintf(os.Stderr, "test %d: expected %d answers, got %d\ninput:\n%s", i+1, len(expect[i]), len(got[i]), formatSingleInput(tests[i]))
			os.Exit(1)
		}
		for j := range expect[i] {
			if expect[i][j] != got[i][j] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d query %d: expected %d got %d\ninput:\n%s", i+1, j+1, expect[i][j], got[i][j], formatSingleInput(tests[i]))
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2057E2-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
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

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutputs(output string, tests []testCase) ([][]int, error) {
	fields := strings.Fields(output)
	total := 0
	for _, tc := range tests {
		total += len(tc.queries)
	}
	if len(fields) < total {
		return nil, fmt.Errorf("expected %d numbers, got %d", total, len(fields))
	}
	if len(fields) > total {
		return nil, fmt.Errorf("extra output detected after %d numbers", total)
	}
	res := make([][]int, len(tests))
	pos := 0
	for i, tc := range tests {
		ans := make([]int, len(tc.queries))
		for j := range ans {
			val, err := strconv.Atoi(fields[pos])
			if err != nil {
				return nil, fmt.Errorf("failed to parse number %d (%q): %v", pos+1, fields[pos], err)
			}
			ans[j] = val
			pos++
		}
		res[i] = ans
	}
	return res, nil
}

func renderInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		n := tc.n
		m := len(tc.edges)
		q := len(tc.queries)
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, q))
		for _, e := range tc.edges {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", e.u, e.v, e.w))
		}
		for _, qu := range tc.queries {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", qu.a, qu.b, qu.k))
		}
	}
	return sb.String()
}

func formatSingleInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	n := tc.n
	m := len(tc.edges)
	q := len(tc.queries)
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, q))
	for _, e := range tc.edges {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", e.u, e.v, e.w))
	}
	for _, qu := range tc.queries {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", qu.a, qu.b, qu.k))
	}
	return sb.String()
}

func generateTests() []testCase {
	// Sample-like basic test
	tests := []testCase{
		{
			n: 4,
			edges: []edge{
				{1, 2, 2}, {2, 4, 2}, {1, 3, 4}, {3, 4, 1},
			},
			queries: []query{{1, 4, 2}, {2, 3, 1}},
		},
	}

	totalN := 0
	for _, tc := range tests {
		totalN += tc.n
	}

	const maxTotalN = 200
	rng := rand.New(rand.NewSource(2057))

	for totalN < maxTotalN {
		n := rng.Intn(40) + 2
		if totalN+n > maxTotalN {
			n = maxTotalN - totalN
		}
		maxEdges := n * (n - 1) / 2
		m := n - 1 + rng.Intn(minInt(maxEdges-(n-1), 3*n)+1)
		edges := randomGraph(rng, n, m)
		dist := allPairsDist(n, edges)

		q := rng.Intn(50) + 1
		queries := make([]query, q)
		for i := 0; i < q; i++ {
			a := rng.Intn(n) + 1
			b := rng.Intn(n) + 1
			for a == b {
				b = rng.Intn(n) + 1
			}
			d := dist[a-1][b-1]
			k := 1
			if d > 0 {
				k = rng.Intn(d) + 1
			}
			queries[i] = query{a: a, b: b, k: k}
		}

		tests = append(tests, testCase{n: n, edges: edges, queries: queries})
		totalN += n
		if len(tests) > 60 {
			break
		}
	}

	return tests
}

func randomGraph(rng *rand.Rand, n, m int) []edge {
	type pair struct{ u, v int }
	seen := make(map[pair]struct{})
	edges := make([]edge, 0, m)

	// start with a tree
	for v := 2; v <= n; v++ {
		p := rng.Intn(v-1) + 1
		e := pair{u: p, v: v}
		if e.u > e.v {
			e.u, e.v = e.v, e.u
		}
		seen[e] = struct{}{}
		edges = append(edges, edge{u: e.u, v: e.v, w: rng.Intn(1_000_000_000) + 1})
	}

	maxEdges := n * (n - 1) / 2
	for len(edges) < m && len(seen) < maxEdges {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		e := pair{u: u, v: v}
		if _, ok := seen[e]; ok {
			continue
		}
		seen[e] = struct{}{}
		edges = append(edges, edge{u: u, v: v, w: rng.Intn(1_000_000_000) + 1})
	}
	return edges
}

func allPairsDist(n int, edges []edge) [][]int {
	const inf = 1 << 30
	dist := make([][]int, n)
	for i := 0; i < n; i++ {
		dist[i] = make([]int, n)
		for j := 0; j < n; j++ {
			if i == j {
				dist[i][j] = 0
			} else {
				dist[i][j] = inf
			}
		}
	}
	for _, e := range edges {
		u, v := e.u-1, e.v-1
		dist[u][v], dist[v][u] = 1, 1
	}
	// Floyd-Warshall for small n.
	for k := 0; k < n; k++ {
		for i := 0; i < n; i++ {
			if dist[i][k] == inf {
				continue
			}
			ik := dist[i][k]
			for j := 0; j < n; j++ {
				if dist[k][j] == inf {
					continue
				}
				if ik+dist[k][j] < dist[i][j] {
					dist[i][j] = ik + dist[k][j]
				}
			}
		}
	}
	return dist
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
