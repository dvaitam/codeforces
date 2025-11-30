package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

// Embedded copy of testcasesC.txt so the verifier is self-contained.
const testcasesRaw = `
4 2 2 1 4 5 1 9 10 2 10 16 3 2
3 1 3 3 9 9 2 9 14 2 6
4 3 4 2 2 8 1 10 10 3 10 15 3 10
4 1 1 1 2 7 2 2 14 3 5 21 1 9
5 2 1 3 3 8 2 7 12 3 3 19 2 10 24 1 10
3 1 4 1 8 7 3 8 16 2 3
5 3 2 2 10 8 1 2 15 3 7 20 2 4 24 1 6
2 1 2 3 9 6 1 7
1 1 3 2 2
2 1 4 2 5 7 2 6
1 1 3 3 7
3 2 2 1 5 7 3 6 15 3 3
4 3 4 1 6 10 2 6 17 1 8 20 2 3
2 1 3 2 4 8 2 8
5 5 2 3 7 10 1 6 11 2 6 17 2 1 22 1 5
3 1 1 1 6 3 1 1 7 1 2
1 1 2 1 9
1 1 2 1 7
4 1 4 1 7 7 1 10 12 2 8 15 2 2
3 1 2 3 3 9 3 6 16 3 9
5 1 4 1 9 11 1 8 15 2 5 20 3 9 26 2 4
4 4 1 2 1 9 1 7 13 2 8 17 2 5
4 2 1 1 5 2 1 6 5 3 3 14 3 4
4 2 2 2 1 7 3 6 14 1 10 19 3 8
2 2 3 3 4 6 1 3
4 2 1 3 7 8 2 4 13 3 6 22 2 6
1 1 2 1 3
4 2 3 3 1 9 1 2 16 3 4 20 3 8
4 1 1 3 8 5 2 10 11 1 1 13 3 9
2 2 3 1 10 9 1 7
5 3 2 3 5 5 1 6 9 3 8 17 2 1 22 3 3
1 1 4 1 1
3 2 1 2 5 6 3 9 11 3 9
2 2 1 3 5 7 2 2
1 1 1 3 2
4 1 2 2 2 6 3 4 13 3 8 19 2 6
2 2 3 1 5 4 3 6
3 1 1 3 1 6 1 5 10 2 2
5 5 4 1 2 10 1 5 15 3 5 18 1 5 21 2 9
1 1 2 2 10
4 4 3 3 4 10 3 8 16 2 6 19 1 8
1 1 1 2 2
2 2 1 2 9 6 3 4
3 2 4 2 8 7 3 9 13 1 3
3 1 2 3 5 6 1 7 9 3 3
2 2 1 3 8 5 2 3
5 4 2 2 9 7 1 2 8 1 4 15 3 5 23 3 7
1 1 2 3 4
1 1 4 1 8
5 5 3 1 10 7 3 6 14 3 5 18 2 4 25 2 5
2 2 3 3 6 11 2 10
5 3 2 1 7 6 3 3 13 2 5 20 2 9 23 2 8
1 1 1 2 1
2 1 3 3 5 6 1 10
3 1 4 2 6 8 3 7 11 1 1
3 3 2 2 4 7 1 8 12 3 4
1 1 1 3 3
4 3 3 3 10 7 2 3 13 2 8 19 3 2
2 1 3 3 7 12 3 3
2 1 1 1 8 4 3 7
3 1 3 3 9 6 1 10 8 1 7
4 1 3 3 3 9 3 10 13 1 7 14 3 7
5 3 4 2 9 9 3 3 14 3 3 20 2 5 25 2 1
5 1 1 3 5 6 3 3 11 1 6 15 3 2 20 1 5
3 3 4 2 9 6 2 8 14 3 8
5 1 2 1 7 6 1 3 8 2 10 10 3 5 17 3 8
5 3 4 3 6 8 3 7 15 3 1 19 2 2 26 3 4
5 5 1 1 10 4 2 3 10 2 8 15 3 4 23 3 3
4 2 4 3 5 11 1 8 16 2 10 23 3 3
4 2 3 2 5 7 2 4 14 1 10 18 2 4
3 3 2 3 10 8 2 4 13 3 4
5 4 2 2 7 7 2 7 11 1 4 14 3 8 19 3 5
4 3 4 1 2 7 1 5 13 3 7 21 1 2
1 1 3 3 10
5 2 1 1 2 4 3 7 9 2 5 13 3 4 18 1 5
3 1 3 1 2 8 1 4 11 3 8
3 1 3 3 4 7 1 10 11 1 3
5 1 3 1 3 9 3 8 15 3 4 20 2 5 26 1 5
1 1 4 1 10
5 2 4 1 7 8 1 2 14 2 7 18 2 7 22 3 1
5 4 4 1 8 9 1 2 13 3 8 17 1 4 19 3 7
1 1 2 1 8
2 1 2 3 4 11 3 5
5 1 3 2 4 6 2 4 9 1 4 16 1 10 21 2 6
1 1 4 2 10
2 1 4 1 1 10 1 9
2 1 3 2 2 8 1 5
4 4 1 1 8 4 3 3 11 2 3 17 2 4
2 1 4 3 7 8 2 8
3 3 1 1 8 3 1 10 7 2 10
5 4 2 1 2 7 3 4 11 1 10 14 3 10 18 1 5
3 3 4 3 5 10 3 2 17 1 6
3 3 4 2 1 9 1 2 12 3 4
3 1 3 2 10 8 3 10 15 2 2
3 1 3 2 2 8 2 1 16 1 5
4 2 2 1 5 6 3 5 13 2 5 19 3 8
2 1 2 2 5 4 2 5
5 2 3 2 3 6 2 8 12 3 6 21 2 6 23 3 8
1 1 4 3 8
4 1 4 1 5 8 2 9 12 1 3 13 2 6
`

type Edge struct {
	to, rev, cap int
	cost, flow   int64
}

type MCMF struct {
	n     int
	graph [][]Edge
	dist  []int64
	prevV []int
	prevE []int
}

func NewMCMF(n int) *MCMF {
	return &MCMF{
		n:     n,
		graph: make([][]Edge, n),
		dist:  make([]int64, n),
		prevV: make([]int, n),
		prevE: make([]int, n),
	}
}

func (m *MCMF) AddEdge(u, v, cap int, cost int64) {
	m.graph[u] = append(m.graph[u], Edge{to: v, rev: len(m.graph[v]), cap: cap, cost: cost})
	m.graph[v] = append(m.graph[v], Edge{to: u, rev: len(m.graph[u]) - 1, cap: 0, cost: -cost})
}

func (m *MCMF) minCostFlow(s, t, maxf int) (int, int64) {
	const INF = int64(4e18)
	flow := 0
	cost := int64(0)
	for flow < maxf {
		inq := make([]bool, m.n)
		for i := 0; i < m.n; i++ {
			m.dist[i] = INF
		}
		m.dist[s] = 0
		queue := make([]int, 0, m.n)
		queue = append(queue, s)
		inq[s] = true
		head := 0
		for head < len(queue) {
			u := queue[head]
			head++
			inq[u] = false
			for ei, e := range m.graph[u] {
				if e.flow < int64(e.cap) && m.dist[u]+e.cost < m.dist[e.to] {
					m.dist[e.to] = m.dist[u] + e.cost
					m.prevV[e.to] = u
					m.prevE[e.to] = ei
					if !inq[e.to] {
						inq[e.to] = true
						queue = append(queue, e.to)
					}
				}
			}
		}
		if m.dist[t] == INF {
			break
		}
		df := maxf - flow
		v := t
		for v != s {
			e := m.graph[m.prevV[v]][m.prevE[v]]
			if df > e.cap-int(e.flow) {
				df = e.cap - int(e.flow)
			}
			v = m.prevV[v]
		}
		if df <= 0 {
			break
		}
		flow += df
		v = t
		for v != s {
			pe := &m.graph[m.prevV[v]][m.prevE[v]]
			pe.flow += int64(df)
			rev := &m.graph[v][pe.rev]
			rev.flow -= int64(df)
			cost += int64(df) * pe.cost
			v = m.prevV[v]
		}
	}
	return flow, cost
}

type testCase struct {
	n int
	k int
	s []int64
	t []int64
	c []int64
}

func solveCase(tc testCase) string {
	n, k := tc.n, tc.k
	s, tArr, c := tc.s, tc.t, tc.c
	times := make([]int64, 0, 2*n)
	for i := 0; i < n; i++ {
		times = append(times, s[i])
		times = append(times, s[i]+tArr[i])
	}
	sort.Slice(times, func(i, j int) bool { return times[i] < times[j] })
	uniq := times[:0]
	for i, v := range times {
		if i == 0 || v != times[i-1] {
			uniq = append(uniq, v)
		}
	}
	times = uniq
	idx := func(x int64) int {
		return sort.Search(len(times), func(i int) bool { return times[i] >= x })
	}
	m := len(times)
	V := m + 2
	S := m
	T := m + 1
	mf := NewMCMF(V)
	mf.AddEdge(S, 0, k, 0)
	mf.AddEdge(m-1, T, k, 0)
	for i := 0; i < m-1; i++ {
		mf.AddEdge(i, i+1, k, 0)
	}
	taskEdge := make([]struct {
		u  int
		ei int
	}, n)
	for i := 0; i < n; i++ {
		u := idx(s[i])
		v := idx(s[i] + tArr[i])
		uEdges := len(mf.graph[u])
		mf.AddEdge(u, v, 1, -c[i])
		taskEdge[i] = struct {
			u  int
			ei int
		}{u: u, ei: uEdges}
	}
	mf.minCostFlow(S, T, k)
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		e := mf.graph[taskEdge[i].u][taskEdge[i].ei]
		if e.flow > 0 {
			sb.WriteByte('1')
		} else {
			sb.WriteByte('0')
		}
	}
	return sb.String()
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %v", idx+1, err)
		}
		k, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse k: %v", idx+1, err)
		}
		expected := 2 + 3*n
		if len(fields) != expected {
			return nil, fmt.Errorf("line %d: expected %d numbers got %d", idx+1, expected, len(fields))
		}
		tc := testCase{n: n, k: k, s: make([]int64, n), t: make([]int64, n), c: make([]int64, n)}
		for i := 0; i < n; i++ {
			val, err := strconv.ParseInt(fields[2+3*i], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: parse s[%d]: %v", idx+1, i, err)
			}
			tc.s[i] = val
			val, err = strconv.ParseInt(fields[2+3*i+1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: parse t[%d]: %v", idx+1, i, err)
			}
			tc.t[i] = val
			val, err = strconv.ParseInt(fields[2+3*i+2], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: parse c[%d]: %v", idx+1, i, err)
			}
			tc.c[i] = val
		}
		cases = append(cases, tc)
	}
	return cases, nil
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	for i, tc := range cases {
		expected := solveCase(tc)

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
		for j := 0; j < tc.n; j++ {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.s[j], tc.t[j], tc.c[j]))
		}

		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("case %d failed\nexpected: %s\ngot: %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
