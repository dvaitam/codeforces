package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type edge struct {
	u, v int
	a, b int64
	idx  int
}

type dsu struct {
	parent []int
	rank   []int
}

func newDSU(n int) *dsu {
	d := &dsu{parent: make([]int, n), rank: make([]int, n)}
	for i := 0; i < n; i++ {
		d.parent[i] = i
	}
	return d
}

func (d *dsu) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *dsu) union(a, b int) bool {
	ra := d.find(a)
	rb := d.find(b)
	if ra == rb {
		return false
	}
	if d.rank[ra] < d.rank[rb] {
		ra, rb = rb, ra
	}
	d.parent[rb] = ra
	if d.rank[ra] == d.rank[rb] {
		d.rank[ra]++
	}
	return true
}

const testcasesI = `3 3 4 1 3 1 1 2 3 2 2 1 2 5 2
5 9 1 1 4 4 2 1 2 1 3 2 4 5 4 1 3 4 4 2 3 1 3 2 5 1 1 1 5 4 5 3 5 4 1 4 5 1 3
3 2 4 1 3 5 2 1 2 5 1
2 1 2 1 2 2 5
5 9 5 4 5 5 3 1 2 3 4 1 4 5 5 3 5 2 3 1 3 2 3 1 5 5 3 2 5 5 1 2 4 5 5 2 3 2 3
3 2 3 2 3 3 1 1 2 4 3
5 6 4 4 5 1 1 1 5 4 3 2 3 2 5 3 5 4 2 1 2 2 4 2 4 5 4
5 7 1 3 5 5 5 1 5 3 2 1 4 5 4 3 4 2 5 2 5 1 2 1 2 1 2 2 4 3 4
2 1 3 1 2 3 3
2 1 2 1 2 1 5
2 1 2 1 2 3 4
5 4 1 2 3 5 4 1 4 2 4 1 2 5 1 2 4 4 5
5 10 3 2 5 1 5 1 4 5 1 2 3 3 4 2 4 5 1 3 5 5 5 4 5 5 2 1 3 5 2 1 2 3 5 1 5 4 5 3 4 2 1
4 5 1 2 4 1 5 1 3 4 3 1 2 3 5 2 3 3 3 1 4 3 1
3 3 3 1 2 2 3 2 3 4 5 1 3 4 5
4 5 4 1 4 1 4 2 4 3 4 3 4 1 5 1 2 3 3 2 3 3 5
4 5 4 2 4 2 2 2 3 5 5 1 2 2 3 3 4 4 1 1 4 1 4
2 1 5 1 2 4 1
3 2 5 2 3 1 1 1 3 4 4
4 3 2 3 4 1 2 2 3 3 1 1 3 2 1
4 4 3 1 3 5 3 2 3 4 2 1 2 2 3 2 4 4 2
5 9 1 2 3 1 1 2 5 1 5 3 4 2 2 2 4 5 3 4 5 3 2 1 4 5 4 1 2 1 5 1 5 1 2 3 5 2 2
4 4 4 1 4 4 5 1 2 5 3 2 3 4 4 2 4 5 3
5 10 3 1 2 5 5 2 5 5 3 3 4 2 5 1 3 5 4 1 5 2 4 2 4 4 1 4 5 3 4 2 3 1 1 1 4 3 5 3 5 3 3
5 9 1 3 5 3 5 2 5 1 5 2 3 3 3 2 4 1 3 1 5 1 4 1 3 2 2 3 4 4 3 1 2 5 3 1 4 3 4
4 3 2 1 3 5 3 2 4 4 3 2 3 4 5
3 3 2 1 2 1 4 1 3 3 5 2 3 3 4
4 5 4 2 4 3 2 1 4 5 1 1 3 3 5 1 2 4 3 3 4 4 5
5 5 5 1 3 4 3 1 5 2 1 2 5 1 1 3 4 4 2 2 4 3 3
2 1 1 1 2 2 5
4 5 1 2 3 4 2 3 4 1 1 2 4 5 5 1 4 2 4 1 3 4 1
5 8 2 3 4 1 3 2 5 1 4 4 5 3 2 3 5 2 1 1 2 4 3 1 5 3 3 2 3 5 3 2 4 5 1
5 10 2 1 2 2 1 2 3 4 3 1 4 5 5 3 4 4 4 3 5 5 3 2 4 5 2 2 5 5 2 1 3 4 1 4 5 2 5 1 5 3 5
3 2 2 2 3 3 2 1 2 5 3
2 1 3 1 2 5 2
3 3 1 1 2 3 3 2 3 1 1 1 3 4 4
5 8 1 3 5 5 1 4 5 5 2 1 3 1 1 1 4 4 3 2 3 1 5 2 5 1 2 1 2 2 1 1 5 1 3
4 6 3 1 2 1 4 1 4 4 5 3 4 2 5 1 3 1 5 2 4 5 4 2 3 4 5
5 10 2 4 5 4 5 1 2 2 3 2 3 2 3 1 5 3 5 2 5 1 4 2 4 5 3 3 4 1 2 1 4 5 4 3 5 4 3 1 3 2 4
2 1 3 1 2 1 2
5 4 5 2 5 4 2 1 2 4 1 2 3 5 3 4 5 3 1
5 8 1 3 5 2 4 1 5 4 5 1 2 4 1 2 3 2 1 2 5 3 2 1 4 4 4 2 4 3 5 3 4 4 5
2 1 4 1 2 4 2
2 1 2 1 2 4 2
4 4 2 2 3 2 5 3 4 3 5 1 3 4 3 1 4 2 2
5 4 5 2 5 1 4 3 5 3 1 2 3 4 4 3 4 4 4
4 6 2 3 4 4 1 1 3 3 3 2 3 1 2 2 4 5 1 1 2 5 4 1 4 1 5
2 1 1 1 2 3 2
3 3 3 1 2 1 2 1 3 2 5 2 3 2 3
2 1 1 1 2 2 1
2 1 3 1 2 2 5
4 4 1 1 3 5 3 1 2 3 2 2 4 1 3 2 3 1 5
2 1 5 1 2 2 3
3 2 2 1 3 5 5 1 2 2 3
4 6 3 1 2 4 5 2 4 2 5 2 3 5 3 1 3 1 5 3 4 3 5 1 4 1 3
4 3 3 2 4 5 5 1 3 2 1 1 4 1 2
2 1 3 1 2 5 2
5 8 5 2 4 4 5 1 2 3 4 1 3 4 5 2 3 4 5 4 5 1 5 1 4 2 2 2 5 5 2 3 5 4 4
4 4 3 3 4 4 4 2 3 5 3 1 3 4 3 1 2 1 4
3 2 3 1 3 2 2 2 3 1 4
3 2 2 1 2 3 1 1 3 2 4
2 1 1 1 2 5 5
4 3 5 2 4 3 3 2 3 3 3 1 4 2 5
2 1 4 1 2 1 1
5 7 3 3 4 2 3 2 4 4 1 1 3 1 2 4 5 3 2 3 5 5 4 2 3 1 4 1 4 3 2
4 5 3 1 4 1 2 1 3 2 5 3 4 1 1 2 4 4 5 1 2 5 5
3 2 4 2 3 1 5 1 2 3 3
2 1 1 1 2 3 3
4 6 5 2 4 4 3 1 4 1 1 2 3 4 5 1 2 3 4 3 4 1 2 1 3 3 4
4 6 4 2 4 4 3 1 4 4 3 1 3 3 4 1 2 4 1 3 4 2 1 2 3 3 1
2 1 1 1 2 1 1
3 3 4 2 3 1 3 1 2 5 3 1 3 3 4
3 3 3 1 2 2 4 2 3 3 1 1 3 1 1
3 3 5 1 3 2 5 2 3 1 2 1 2 3 4
5 10 2 1 3 2 5 1 2 1 4 3 4 1 2 2 4 3 4 2 3 1 4 3 5 3 1 1 4 1 1 2 5 3 2 1 5 4 3 4 5 1 4
4 3 5 1 4 4 2 3 4 3 4 2 3 5 4
4 5 3 1 3 2 4 2 3 3 2 1 2 1 3 3 4 3 1 1 4 2 4
4 4 5 1 3 5 4 1 4 3 3 1 2 1 3 3 4 5 2
5 8 2 1 4 5 2 3 5 1 4 4 5 2 4 2 3 3 5 2 4 3 4 1 2 2 1 2 5 4 4 1 5 4 5
4 3 4 2 4 4 1 1 2 1 5 3 4 4 5
3 3 2 2 3 2 3 1 2 3 2 1 3 1 1
4 6 1 1 4 2 4 2 3 2 2 1 3 4 2 1 2 5 3 3 4 2 3 2 4 3 1
4 3 2 1 4 2 5 1 2 5 5 1 3 4 4
3 3 1 1 2 3 2 1 3 1 5 2 3 5 1
5 9 1 2 4 1 5 2 3 1 4 3 5 2 3 3 4 5 1 2 5 1 4 4 5 5 2 1 3 4 3 1 4 3 4 1 2 4 5
4 5 1 1 2 4 5 1 4 2 1 3 4 1 4 1 3 4 2 2 3 4 5
4 4 1 1 4 5 3 2 3 1 5 2 4 2 5 1 2 4 2
4 3 5 1 2 3 2 2 3 2 4 1 4 1 3
3 3 1 1 2 2 1 1 3 5 4 2 3 1 2
4 4 3 1 4 1 3 1 2 5 4 3 4 4 5 2 3 2 2
4 4 1 2 3 4 4 1 3 1 3 2 4 4 2 1 2 5 2
5 9 1 2 5 5 1 3 4 4 3 2 4 4 5 1 2 4 2 3 5 4 3 4 5 2 5 1 3 4 4 1 5 3 4 1 4 5 1
3 3 3 1 2 5 2 1 3 5 3 2 3 3 5
2 1 2 1 2 2 3
5 10 5 1 4 2 4 2 5 1 5 3 5 2 1 1 5 2 5 3 4 4 2 2 4 2 4 1 2 2 3 1 3 1 3 2 3 3 1 4 5 3 4
4 5 3 2 3 2 3 3 4 2 3 1 4 1 2 2 4 5 3 1 2 2 1
2 1 2 1 2 2 2
3 3 4 2 3 1 4 1 2 5 1 1 3 5 4
2 1 1 1 2 4 1
5 8 3 2 5 1 5 1 2 2 3 1 5 5 2 1 4 1 3 1 3 2 5 2 4 4 4 3 5 4 2 3 4 1 1`

func solve(n, m int, k int64, edges []edge) int64 {
	used := make([]int64, m)
	var total int64
	for step := int64(0); step < k; step++ {
		weights := make([]struct {
			w int64
			e edge
		}, m)
		for i, e := range edges {
			w := e.a*(2*used[i]+1) + e.b
			weights[i] = struct {
				w int64
				e edge
			}{w, e}
		}
		sort.Slice(weights, func(i, j int) bool {
			return weights[i].w < weights[j].w
		})
		d := newDSU(n)
		cnt := 0
		for _, item := range weights {
			if d.union(item.e.u, item.e.v) {
				used[item.e.idx]++
				total += item.w
				cnt++
				if cnt == n-1 {
					break
				}
			}
		}
	}
	return total
}

type testCase struct {
	input    string
	expected string
}

func parseCase(line string) (testCase, error) {
	fields := strings.Fields(line)
	if len(fields) < 3 {
		return testCase{}, fmt.Errorf("not enough fields")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return testCase{}, fmt.Errorf("bad n: %w", err)
	}
	m, err := strconv.Atoi(fields[1])
	if err != nil {
		return testCase{}, fmt.Errorf("bad m: %w", err)
	}
	k, err := strconv.ParseInt(fields[2], 10, 64)
	if err != nil {
		return testCase{}, fmt.Errorf("bad k: %w", err)
	}
	if len(fields) != 3+4*m {
		return testCase{}, fmt.Errorf("expected %d edge values, got %d", 4*m, len(fields)-3)
	}
	edges := make([]edge, m)
	pos := 3
	for i := 0; i < m; i++ {
		u, err := strconv.Atoi(fields[pos])
		if err != nil {
			return testCase{}, fmt.Errorf("bad u: %w", err)
		}
		v, err := strconv.Atoi(fields[pos+1])
		if err != nil {
			return testCase{}, fmt.Errorf("bad v: %w", err)
		}
		aVal, err := strconv.ParseInt(fields[pos+2], 10, 64)
		if err != nil {
			return testCase{}, fmt.Errorf("bad a: %w", err)
		}
		bVal, err := strconv.ParseInt(fields[pos+3], 10, 64)
		if err != nil {
			return testCase{}, fmt.Errorf("bad b: %w", err)
		}
		edges[i] = edge{u - 1, v - 1, aVal, bVal, i}
		pos += 4
	}
	exp := solve(n, m, k, edges)

	var input strings.Builder
	input.WriteString("1\n")
	fmt.Fprintf(&input, "%d %d %d\n", n, m, k)
	for _, e := range edges {
		fmt.Fprintf(&input, "%d %d %d %d\n", e.u+1, e.v+1, e.a, e.b)
	}
	return testCase{input: input.String(), expected: strconv.FormatInt(exp, 10)}, nil
}

func loadCases() ([]testCase, error) {
	lines := strings.Split(testcasesI, "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		tc, err := parseCase(line)
		if err != nil {
			return nil, fmt.Errorf("case %d: %w", idx+1, err)
		}
		cases = append(cases, tc)
	}
	return cases, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, string(out))
	}
	return strings.TrimSpace(string(out)), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierI /path/to/binary")
		os.Exit(1)
	}
	cases, err := loadCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load cases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		got, err := runCandidate(os.Args[1], tc.input)
		if err != nil {
			fmt.Printf("case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, tc.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
