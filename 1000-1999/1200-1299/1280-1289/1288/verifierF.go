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

const testcasesRaw = `5 1 4 3 1
UUBBR
B
3 1
1 1
4 1
4 1
2 5 5 3 5
BU
RRURR
2 4
1 2
2 1
1 5
1 3
4 2 5 5 5
UBRU
BR
3 2
1 1
2 1
3 2
2 2
4 2 2 5 2
BBUU
UU
3 2
3 2
5 4 3 2 4
UURBB
RURU
1 4
4 4
1 4
4 3 4 1 2
BBRR
RBU
1 1
1 2
3 1
3 3
4 3 5 3 4
RBRR
BRB
4 2
2 3
4 2
3 3
4 2
3 2 5 4 5
RRB
RR
3 1
2 1
3 2
1 1
1 1
3 3 4 5 1
UUR
BUR
3 2
1 2
1 1
2 2
2 2 5 3 2
UB
BU
1 1
2 2
1 1
1 1
1 2
4 4 1 2 2
RRRR
BRRB
1 3
3 3 3 1 1
BRR
RRR
2 2
1 3
1 1
5 2 5 1 2
RRUUB
UB
4 2
1 1
1 2
5 1
4 2
1 2 1 3 5
U
UB
1 1
3 2 5 5 4
UUB
BB
1 2
3 1
2 1
3 1
3 2
2 3 4 5 4
UR
BBB
1 3
1 3
2 2
1 1
2 2 4 1 1
UU
RU
1 1
2 1
1 1
2 1
5 3 1 4 1
BBUUB
RUU
2 1
1 4 3 3 4
U
UBBB
1 3
1 3
1 1
5 5 1 3 3
BUUUR
BRRBU
3 2
1 1 3 1 4
U
R
1 1
1 1
1 1
4 5 2 2 4
BRBB
BRRRU
3 1
2 4
1 2 5 1 1
B
UU
1 2
1 2
1 1
1 1
1 1
3 4 3 3 2
RRB
BRUU
2 2
1 4
3 2
1 2 4 4 2
U
BU
1 2
1 2
1 1
1 1
1 1 5 5 5
B
U
1 1
1 1
1 1
1 1
1 1
2 2 2 4 4
BU
BB
1 2
1 2
2 3 3 1 3
UU
BUR
2 2
1 1
2 3
1 4 5 3 1
B
RUUR
1 4
1 2
1 4
1 1
1 4
3 1 4 1 4
BUB
U
2 1
2 1
2 1
2 1
5 2 5 1 5
RBBRU
UR
1 1
5 2
2 1
3 1
3 2
3 2 2 5 1
URB
UR
1 2
3 1
4 4 5 5 3
UBBU
RUBB
2 4
3 2
1 4
3 3
3 3
4 3 1 5 4
BURU
BBR
1 1
1 1 4 4 5
U
U
1 1
1 1
1 1
1 1
4 5 1 1 3
RRBR
RUUUB
4 2
2 3 1 1 3
UB
UUR
2 2
2 1 3 2 2
UU
U
1 1
2 1
2 1
3 5 4 2 2
BUU
BURUB
2 4
3 3
3 5
1 5
2 1 4 2 4
BB
U
1 1
2 1
2 1
2 1
3 5 5 3 5
URB
BURRB
3 4
1 1
2 5
1 3
3 4
4 3 1 4 4
RUUB
URR
1 3
1 5 5 3 3
U
URBUR
1 1
1 2
1 5
1 4
1 1
5 3 4 3 2
RRUUB
UUU
1 2
1 3
4 2
3 2
2 2 3 2 1
UU
UR
2 1
2 2
2 1
3 3 2 5 1
BUU
RBU
3 2
3 1
1 1 1 5 2
B
U
1 1
2 4 2 1 2
RU
BBRU
1 1
2 1
4 5 2 2 4
BBBB
BRURB
2 2
2 3
4 3 4 1 4
UURU
BRR
1 3
3 2
3 2
4 2
5 4 3 2 5
BBRBU
BRBU
5 1
2 1
2 3
4 5 1 2 5
RUBU
RRURR
2 3
5 1 5 3 2
RBRUB
B
5 1
1 1
2 1
5 1
2 1
4 1 5 3 3
UBRR
B
1 1
2 1
2 1
2 1
2 1
3 3 5 2 2
UUB
BRU
1 1
3 1
2 2
2 2
2 3
2 4 5 5 4
RB
RRBB
1 1
1 4
2 1
2 3
1 1
1 1 2 2 4
B
B
1 1
1 1
1 1 4 1 3
R
B
1 1
1 1
1 1
1 1
2 2 4 4 4
BU
RB
1 1
2 2
1 1
2 2
5 1 5 3 2
BRBRR
R
4 1
3 1
3 1
1 1
3 1
2 3 5 1 1
BU
UUB
2 1
2 1
1 1
1 1
2 1
2 2 1 5 2
BB
UB
2 2
1 1 3 5 2
R
U
1 1
1 1
1 1
5 5 5 1 3
RUUUR
UUBUB
3 3
1 5
2 3
4 3
1 2
3 1 4 5 4
URR
B
1 1
2 1
3 1
1 1
1 1 1 3 1
R
U
1 1
3 3 2 2 2
RBB
RUU
1 1
2 1
5 4 5 4 4
RUBUU
RUUU
3 1
4 3
3 2
4 2
1 3
2 2 2 2 1
UR
BU
2 1
2 2
5 1 1 3 1
BRBBR
R
2 1
4 4 1 4 2
URUU
UUUU
2 2
5 2 4 2 1
UBBBB
BB
3 1
1 2
1 2
2 1
1 5 4 3 2
B
BRUUR
1 3
1 5
1 1
1 1
1 4 5 1 4
B
URBR
1 1
1 1
1 1
1 3
1 4
1 3 1 3 5
U
RRB
1 1
5 3 3 3 1
BRBBR
URR
1 1
3 2
3 1
4 3 2 5 3
UURU
UBR
2 2
2 2
5 3 4 3 2
UBUUU
BBB
4 3
2 1
4 3
3 2
2 5 4 2 2
UU
URRRB
1 2
2 1
2 3
1 4
2 1 5 2 1
BU
R
2 1
2 1
1 1
2 1
1 1
5 4 5 1 5
BRRBU
RURU
3 3
3 2
4 2
3 1
3 4
3 3 2 2 3
RBU
UUB
1 3
2 2
4 1 2 3 1
BUBU
R
1 1
4 1
5 5 3 1 5
RBRUB
BUBUR
1 5
1 4
2 2
4 2 3 2 5
BURR
UU
2 1
1 2
2 1
1 5 2 4 2
B
BURUB
1 3
1 5
5 2 5 3 2
BRUBB
RU
3 2
2 1
1 2
3 1
4 2
2 4 4 2 4
UB
BRBU
2 1
2 1
1 4
2 4
3 4 2 1 5
URR
BRRU
2 4
3 4
2 3 2 3 4
BB
BRU
1 1
1 2
2 2 4 1 3
RB
RU
2 1
2 2
2 1
2 1
4 1 2 2 4
BUUB
U
1 1
4 1
3 4 3 1 3
UUU
BRRU
1 2
3 4
1 4
2 2 1 5 4
UU
RU
2 2
2 2 2 1 5
RU
BU
1 2
1 1
5 1 2 4 1
RRURB
R
2 1
5 1
3 1 3 3 3
BRU
B
2 1
2 1
2 1
4 2 1 2 5
BRUU
RB
2 1
4 1 1 1 4
BRBB
U
1 1
1 5 5 1 3
B
UUBBU
1 2
1 1
1 4
1 1
1 2`

type edge struct {
	to, rev int
	cap     int
	cost    int
}

type testCase struct {
	n1, n2, m int
	R, B      int
	s1, s2    string
	edges     [][2]int
}

func parseTestcases(raw string) ([]testCase, error) {
	scan := bufio.NewScanner(strings.NewReader(raw))
	scan.Split(bufio.ScanWords)
	tokens := make([]string, 0)
	for scan.Scan() {
		tokens = append(tokens, scan.Text())
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}
	pos := 0
	var tests []testCase
	for pos < len(tokens) {
		if pos+5 > len(tokens) {
			return nil, fmt.Errorf("incomplete header at token %d", pos+1)
		}
		n1, err1 := strconv.Atoi(tokens[pos])
		n2, err2 := strconv.Atoi(tokens[pos+1])
		m, err3 := strconv.Atoi(tokens[pos+2])
		R, err4 := strconv.Atoi(tokens[pos+3])
		B, err5 := strconv.Atoi(tokens[pos+4])
		if err := firstErr(err1, err2, err3, err4, err5); err != nil {
			return nil, fmt.Errorf("bad header at case %d: %w", len(tests)+1, err)
		}
		pos += 5
		if pos+2 > len(tokens) {
			return nil, fmt.Errorf("missing s1/s2 at case %d", len(tests)+1)
		}
		s1 := tokens[pos]
		s2 := tokens[pos+1]
		pos += 2
		if pos+2*m > len(tokens) {
			return nil, fmt.Errorf("missing edges at case %d", len(tests)+1)
		}
		edges := make([][2]int, m)
		for i := 0; i < m; i++ {
			u, errU := strconv.Atoi(tokens[pos])
			v, errV := strconv.Atoi(tokens[pos+1])
			if err := firstErr(errU, errV); err != nil {
				return nil, fmt.Errorf("invalid edge in case %d", len(tests)+1)
			}
			edges[i] = [2]int{u, v}
			pos += 2
		}
		tests = append(tests, testCase{n1: n1, n2: n2, m: m, R: R, B: B, s1: s1, s2: s2, edges: edges})
	}
	return tests, nil
}

func firstErr(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

func minCostSolution(tc testCase) (string, error) {
	n1, n2, m, R, B := tc.n1, tc.n2, tc.m, tc.R, tc.B
	s1, s2 := tc.s1, tc.s2

	N := n1 + n2 + 4
	s := n1 + n2
	t := s + 1
	S := t + 1
	T := S + 1
	graph := make([][]edge, N)
	deg := make([]int, N)

	addEdge := func(u, v, cap, cost int) {
		graph[u] = append(graph[u], edge{to: v, rev: len(graph[v]), cap: cap, cost: cost})
		graph[v] = append(graph[v], edge{to: u, rev: len(graph[u]) - 1, cap: 0, cost: -cost})
	}
	addLower := func(u, v, l, r, cost int) {
		deg[u] -= l
		deg[v] += l
		addEdge(u, v, r-l, cost)
	}

	INF_CAP := m + n1 + n2 + 5
	for i := 0; i < n1; i++ {
		c := s1[i]
		if c == 'R' {
			addLower(s, i, 1, INF_CAP, 0)
		} else if c == 'B' {
			addLower(i, t, 1, INF_CAP, 0)
		} else {
			addLower(s, i, 0, INF_CAP, 0)
			addLower(i, t, 0, INF_CAP, 0)
		}
	}
	for i := 0; i < n2; i++ {
		c := s2[i]
		node := n1 + i
		if c == 'R' {
			addLower(node, t, 1, INF_CAP, 0)
		} else if c == 'B' {
			addLower(s, node, 1, INF_CAP, 0)
		} else {
			addLower(s, node, 0, INF_CAP, 0)
			addLower(node, t, 0, INF_CAP, 0)
		}
	}

	us := make([]int, m)
	vs := make([]int, m)
	reIdx := make([]int, m)
	beIdx := make([]int, m)
	for i, e := range tc.edges {
		u := e[0] - 1
		v := e[1] - 1 + n1
		us[i] = u
		vs[i] = v
		addEdge(u, v, 1, R)
		reIdx[i] = len(graph[u]) - 1
		addEdge(v, u, 1, B)
		beIdx[i] = len(graph[v]) - 1
	}

	addEdge(t, s, INF_CAP, 0)
	sum := 0
	for i := 0; i < N; i++ {
		if deg[i] > 0 {
			addEdge(S, i, deg[i], 0)
			sum += deg[i]
		} else if deg[i] < 0 {
			addEdge(i, T, -deg[i], 0)
		}
	}

	const INF_COST = int(1e18)
	flow, cost := 0, 0
	dist := make([]int, N)
	inQ := make([]bool, N)

	for {
		for i := 0; i < N; i++ {
			dist[i] = INF_COST
			inQ[i] = false
		}
		queue := make([]int, 0, N)
		dist[S] = 0
		queue = append(queue, S)
		inQ[S] = true
		for qi := 0; qi < len(queue); qi++ {
			u := queue[qi]
			inQ[u] = false
			for _, e := range graph[u] {
				if e.cap > 0 && dist[e.to] > dist[u]+e.cost {
					dist[e.to] = dist[u] + e.cost
					if !inQ[e.to] {
						inQ[e.to] = true
						queue = append(queue, e.to)
					}
				}
			}
		}
		if dist[T] == INF_COST {
			break
		}

		cur := make([]int, N)
		var dfs func(int, int) int
		visited := make([]bool, N)
		dfs = func(u, f int) int {
			if u == T {
				flow += f
				cost += f * dist[T]
				return f
			}
			visited[u] = true
			for i := cur[u]; i < len(graph[u]); i++ {
				e := &graph[u][i]
				if e.cap > 0 && !visited[e.to] && dist[e.to] == dist[u]+e.cost {
					pushed := dfs(e.to, min(f, e.cap))
					if pushed > 0 {
						e.cap -= pushed
						graph[e.to][e.rev].cap += pushed
						return pushed
					}
				}
				cur[u]++
			}
			return 0
		}
		for {
			for i := range visited {
				visited[i] = false
			}
			if dfs(S, INF_CAP) == 0 {
				break
			}
		}
	}

	if flow < sum {
		return "-1", nil
	}

	res := make([]byte, m)
	for i := 0; i < m; i++ {
		if graph[us[i]][reIdx[i]].cap == 0 {
			res[i] = 'R'
		} else if graph[vs[i]][beIdx[i]].cap == 0 {
			res[i] = 'B'
		} else {
			res[i] = 'U'
		}
	}

	return fmt.Sprintf("%d\n%s", cost, string(res)), nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases(testcasesRaw)
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid test data:", err)
		os.Exit(1)
	}

	for i, tc := range tests {
		expected, err := minCostSolution(tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i+1, err)
			os.Exit(1)
		}

		var input strings.Builder
		fmt.Fprintf(&input, "%d %d %d %d %d\n", tc.n1, tc.n2, tc.m, tc.R, tc.B)
		input.WriteString(tc.s1)
		input.WriteByte('\n')
		input.WriteString(tc.s2)
		input.WriteByte('\n')
		for _, e := range tc.edges {
			fmt.Fprintf(&input, "%d %d\n", e[0], e[1])
		}

		got, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Printf("case %d failed\nexpected:\n%s\ngot:\n%s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
