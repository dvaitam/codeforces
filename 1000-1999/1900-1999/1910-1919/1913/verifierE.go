package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `
4 3 0 0 1 1 0 1 1 1 1 0 1 0 2 2 1 2 2 2 0
4 2 0 1 1 0 0 1 1 1 2 1 1 0 1 3
2 2 1 1 0 1 1 2 2 1
3 2 0 0 1 0 0 1 2 1 0 3 0
3 3 1 0 1 0 1 1 1 1 1 0 0 1 3 3 0
2 4 0 1 1 1 0 0 0 1 0 4 2 2 0 0
4 4 1 0 1 1 1 1 0 0 1 0 0 1 1 1 1 1 0 2 4 4 0 3 0 0
2 2 1 1 0 0 1 1 0 1
3 2 1 0 1 0 1 1 0 2 0 1 1
4 2 0 1 0 1 1 0 1 1 2 0 1 0 2 1
4 3 1 0 1 1 0 1 0 1 0 0 0 0 1 3 0 0 2 3 2
3 2 0 1 0 1 1 1 2 0 0 0 2
2 2 1 0 0 1 2 1 1 2
4 3 1 1 0 0 1 1 1 0 0 0 1 0 1 2 1 2 4 2 0
4 4 0 0 1 0 0 1 0 0 1 0 1 1 1 1 0 1 0 0 1 2 2 3 0 3
2 4 0 0 0 0 1 0 0 1 1 2 0 2 0 1
4 4 0 1 0 0 1 1 0 1 0 1 1 0 1 0 1 0 2 2 4 2 4 4 1 4
3 4 1 1 0 0 0 0 0 1 0 1 0 1 3 4 1 2 2 3 2
2 3 0 1 1 1 1 0 2 0 0 0 2
3 4 0 0 0 1 0 0 0 1 0 1 0 0 2 2 2 3 2 0 1
2 2 0 0 0 0 1 0 1 1
2 2 0 0 1 0 0 0 1 0
4 2 0 1 1 0 0 1 1 0 1 0 1 0 3 1
2 4 0 1 0 0 1 1 1 0 2 4 2 2 1 1
2 4 1 1 0 0 0 0 0 0 4 2 2 0 2 2
3 4 1 0 0 0 0 1 0 1 0 0 1 0 3 3 3 3 3 3 0
2 3 0 1 1 0 1 0 2 0 2 0 2
4 2 1 0 1 0 1 0 0 1 0 0 2 2 2 1
4 4 1 0 0 0 0 1 1 1 1 1 0 1 1 0 1 1 4 0 2 0 0 0 3 3
4 3 0 1 1 0 0 1 1 1 0 1 0 0 0 3 3 2 3 4 4
4 4 1 0 1 0 1 0 0 1 0 1 0 1 1 0 1 1 1 2 2 0 1 0 0 1
4 4 1 0 0 0 1 0 0 0 0 0 1 0 0 0 1 1 1 0 4 4 2 1 4 2
2 2 1 0 1 0 0 2 1 1
2 2 1 0 0 1 0 2 2 2
2 2 1 0 0 0 1 1 1 1
2 3 1 0 0 1 0 1 2 0 1 1 0
3 4 1 0 0 0 1 1 1 0 0 0 0 1 4 4 4 1 0 2 2
3 3 0 1 0 0 0 1 0 1 1 2 0 1 0 1 2
3 4 0 1 0 0 1 1 0 0 0 0 1 1 0 0 1 1 0 0 3
2 4 1 0 0 1 1 1 1 1 2 0 0 1 1 2
3 4 1 1 1 1 0 1 1 0 0 0 1 0 2 4 0 2 2 2 0
3 4 0 1 1 0 0 0 1 1 0 0 0 0 1 4 4 3 3 1 2
2 4 1 0 1 1 0 1 1 0 4 3 2 2 2 1
3 4 0 0 0 0 0 1 1 0 0 0 1 1 3 2 4 2 3 1 1
4 4 0 1 0 1 1 0 1 1 1 0 0 0 0 1 0 0 3 2 1 1 1 3 2 4
2 4 1 0 1 0 0 0 0 1 1 0 0 0 0 1
2 2 0 0 1 0 2 2 0 2
4 2 1 0 1 0 1 1 0 0 1 1 0 1 2 1
3 4 1 1 0 1 1 1 1 1 0 0 0 0 0 4 4 0 2 3 2
3 2 1 0 0 0 0 0 1 1 2 1 3
4 3 1 1 0 1 0 1 1 0 0 1 1 0 3 3 3 3 4 4 4
3 4 1 1 1 1 0 1 1 1 1 0 0 1 3 3 1 3 2 0 3
3 3 1 1 1 0 1 1 1 1 1 0 3 0 1 0 0
4 2 0 0 1 0 0 0 0 0 2 2 0 2 4 1
4 3 1 0 1 0 1 0 1 1 0 1 1 1 3 2 3 1 0 0 3
3 4 0 1 1 1 1 1 1 1 0 0 0 0 1 4 4 2 2 2 3
3 2 1 1 1 1 0 1 0 1 0 0 3
4 4 0 0 1 1 1 1 1 1 1 0 1 1 1 0 0 1 4 2 3 0 3 1 3 2
2 2 0 0 0 1 2 1 1 2
4 3 1 0 0 0 0 1 0 0 1 0 0 0 3 3 0 1 3 4 0
3 3 1 1 0 0 1 0 1 1 0 2 1 0 3 0 0
2 2 0 0 1 0 2 1 0 1
2 3 0 1 0 1 0 0 0 0 2 2 0
3 2 1 1 1 1 0 1 1 2 2 0 1
2 2 1 0 0 0 0 2 1 2
2 3 0 1 1 0 1 1 3 1 2 2 0
4 2 0 1 0 1 1 0 1 0 1 0 2 2 1 1
2 3 1 0 0 1 1 1 3 0 2 0 1
2 3 0 1 1 0 0 0 3 3 2 2 2
4 3 0 1 1 1 0 1 1 1 1 0 0 1 1 3 0 3 3 3 1
3 3 0 1 1 1 0 0 1 0 0 3 3 1 3 3 1
2 2 1 0 0 1 2 0 1 1
3 3 0 0 0 1 1 1 1 1 0 2 3 0 3 0 0
2 3 1 1 0 0 1 0 1 3 2 2 0
4 4 1 1 1 0 0 0 0 0 0 1 0 0 0 0 0 1 3 3 4 0 2 2 3 4
4 4 1 1 0 1 0 0 1 1 0 1 0 1 1 0 0 0 4 0 4 1 3 0 3 3
4 4 0 0 0 0 0 1 1 0 1 0 1 1 1 1 1 0 4 2 0 1 2 2 0 3
4 2 1 1 1 1 0 0 0 0 2 1 0 1 1 4
3 3 1 0 1 1 1 1 0 0 1 2 3 3 3 1 1
3 4 1 1 1 1 1 1 0 1 0 1 1 1 3 4 3 0 1 2 3
3 4 0 1 1 1 0 1 0 0 1 0 1 1 4 2 3 3 3 3 0
4 2 1 1 0 0 0 1 1 1 1 0 1 1 2 4
4 3 1 1 1 1 0 1 1 1 0 1 0 0 0 1 0 3 3 0 1
2 4 1 1 1 0 0 1 1 0 4 0 0 2 1 1
4 2 1 0 1 0 1 1 1 0 2 1 2 1 4 2
4 4 1 1 0 0 0 1 1 1 1 1 1 0 1 1 0 1 4 4 2 1 2 2 3 4
4 4 1 1 1 1 1 0 0 1 0 0 0 0 0 0 0 1 4 0 0 0 3 4 2 4
2 2 1 0 0 0 1 2 1 2
3 2 1 1 0 0 0 1 0 1 0 3 2
4 3 0 1 0 0 0 0 1 1 1 1 1 1 0 3 0 3 3 0 3
4 4 1 1 1 0 1 1 1 1 0 1 1 0 1 0 1 0 4 1 0 2 4 1 3 0
4 3 0 0 0 1 1 0 1 0 1 0 1 1 2 3 1 2 3 2 3
3 4 1 1 1 0 0 1 0 1 1 1 1 0 3 2 0 0 3 2 0
2 3 0 1 1 0 0 0 1 3 0 1 2
3 3 1 1 1 1 0 1 1 1 1 1 1 2 2 2 0
3 4 0 0 0 0 0 1 0 1 0 0 0 1 4 4 3 3 2 3 3
3 4 1 1 0 0 1 1 0 0 0 0 0 1 0 3 1 2 2 3 2
4 2 0 0 0 1 0 0 1 0 2 2 2 2 3 3
3 4 1 0 0 1 1 0 0 1 0 1 1 0 3 0 2 1 0 3 0
2 2 0 1 1 0 1 1 0 2
`

type edge struct {
	to   int
	rev  int
	cap  int
	cost int
}

type MCMF struct {
	n     int
	graph [][]edge
	dist  []int
	prevv []int
	preve []int
}

func NewMCMF(n int) *MCMF {
	g := make([][]edge, n)
	return &MCMF{n: n, graph: g, dist: make([]int, n), prevv: make([]int, n), preve: make([]int, n)}
}

func (f *MCMF) AddEdge(u, v, cap, cost int) {
	f.graph[u] = append(f.graph[u], edge{to: v, rev: len(f.graph[v]), cap: cap, cost: cost})
	f.graph[v] = append(f.graph[v], edge{to: u, rev: len(f.graph[u]) - 1, cap: 0, cost: -cost})
}

const INF = int(1e9)

func (f *MCMF) MinCostFlow(s, t, maxf int) (int, int) {
	res := 0
	flow := 0
	for maxf > 0 {
		for i := 0; i < f.n; i++ {
			f.dist[i] = INF
		}
		inq := make([]bool, f.n)
		q := make([]int, 0)
		f.dist[s] = 0
		inq[s] = true
		q = append(q, s)
		for idx := 0; idx < len(q); idx++ {
			v := q[idx]
			inq[v] = false
			for i, e := range f.graph[v] {
				if e.cap > 0 && f.dist[e.to] > f.dist[v]+e.cost {
					f.dist[e.to] = f.dist[v] + e.cost
					f.prevv[e.to] = v
					f.preve[e.to] = i
					if !inq[e.to] {
						q = append(q, e.to)
						inq[e.to] = true
					}
				}
			}
		}
		if f.dist[t] == INF {
			break
		}
		d := maxf
		for v := t; v != s; v = f.prevv[v] {
			if d > f.graph[f.prevv[v]][f.preve[v]].cap {
				d = f.graph[f.prevv[v]][f.preve[v]].cap
			}
		}
		maxf -= d
		flow += d
		for v := t; v != s; v = f.prevv[v] {
			e := &f.graph[f.prevv[v]][f.preve[v]]
			e.cap -= d
			rev := &f.graph[v][e.rev]
			rev.cap += d
		}
		res += d * f.dist[t]
	}
	return res, flow
}

type testCase struct {
	n, m int
	grid [][]int
	row  []int
	col  []int
}

func parseTestcases(raw string) ([]testCase, error) {
	sc := bufio.NewScanner(strings.NewReader(raw))
	sc.Buffer(make([]byte, 1024), 1<<20)
	cases := make([]testCase, 0)
	lineNo := 0
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		lineNo++
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("line %d: not enough fields", lineNo)
		}
		idx := 0
		n, err := strconv.Atoi(fields[idx])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad n", lineNo)
		}
		idx++
		m, err := strconv.Atoi(fields[idx])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad m", lineNo)
		}
		idx++
		need := n*m + n + m
		if len(fields) != 2+need {
			return nil, fmt.Errorf("line %d: expected %d numbers got %d", lineNo, 2+need, len(fields))
		}
		grid := make([][]int, n)
		for i := 0; i < n; i++ {
			grid[i] = make([]int, m)
			for j := 0; j < m; j++ {
				v, err := strconv.Atoi(fields[idx])
				if err != nil {
					return nil, fmt.Errorf("line %d: bad grid value", lineNo)
				}
				grid[i][j] = v
				idx++
			}
		}
		row := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[idx])
			if err != nil {
				return nil, fmt.Errorf("line %d: bad row value", lineNo)
			}
			row[i] = v
			idx++
		}
		col := make([]int, m)
		for j := 0; j < m; j++ {
			v, err := strconv.Atoi(fields[idx])
			if err != nil {
				return nil, fmt.Errorf("line %d: bad col value", lineNo)
			}
			col[j] = v
			idx++
		}
		cases = append(cases, testCase{n: n, m: m, grid: grid, row: row, col: col})
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("scan error: %w", err)
	}
	if len(cases) == 0 {
		return nil, fmt.Errorf("no testcases parsed")
	}
	return cases, nil
}

func solve(tc testCase) int {
	sumA, sumB := 0, 0
	for _, v := range tc.row {
		sumA += v
	}
	for _, v := range tc.col {
		sumB += v
	}
	if sumA != sumB {
		return -1
	}
	total1 := 0
	for i := 0; i < tc.n; i++ {
		for j := 0; j < tc.m; j++ {
			if tc.grid[i][j] == 1 {
				total1++
			}
		}
	}
	s := tc.n + tc.m
	t := s + 1
	flow := NewMCMF(t + 1)
	for i := 0; i < tc.n; i++ {
		flow.AddEdge(s, i, tc.row[i], 0)
	}
	for j := 0; j < tc.m; j++ {
		flow.AddEdge(tc.n+j, t, tc.col[j], 0)
	}
	for i := 0; i < tc.n; i++ {
		for j := 0; j < tc.m; j++ {
			cost := 1
			if tc.grid[i][j] == 1 {
				cost = -1
			}
			flow.AddEdge(i, tc.n+j, 1, cost)
		}
	}
	cost, got := flow.MinCostFlow(s, t, sumA)
	if got < sumA {
		return -1
	}
	return total1 + cost
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for i := 0; i < tc.n; i++ {
		for j := 0; j < tc.m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(tc.grid[i][j]))
		}
		sb.WriteByte('\n')
	}
	for i := 0; i < tc.n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(tc.row[i]))
	}
	sb.WriteByte('\n')
	for j := 0; j < tc.m; j++ {
		if j > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(tc.col[j]))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseTestcases(testcasesRaw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		input := buildInput(tc)
		expected := solve(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strconv.Itoa(expected) {
			fmt.Printf("case %d failed: expected %d got %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
