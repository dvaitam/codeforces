package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Edge struct {
	to int
	w  int
}

type Item struct {
	d    int64
	node int
	last int
	idx  int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].d < pq[j].d }
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].idx = i
	pq[j].idx = j
}
func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*Item)
	item.idx = len(*pq)
	*pq = append(*pq, item)
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*pq = old[:n-1]
	return item
}

const INF int64 = 1 << 60

// referenceSolve mirrors 1486E.go for a single test case input.
func referenceSolve(input string) (string, error) {
	fields := strings.Fields(input)
	if len(fields) < 2 {
		return "", fmt.Errorf("invalid input")
	}
	pos := 0
	nextInt := func() (int, error) {
		if pos >= len(fields) {
			return 0, fmt.Errorf("unexpected EOF")
		}
		val, err := strconv.Atoi(fields[pos])
		pos++
		return val, err
	}
	n, err := nextInt()
	if err != nil {
		return "", fmt.Errorf("bad n: %v", err)
	}
	m, err := nextInt()
	if err != nil {
		return "", fmt.Errorf("bad m: %v", err)
	}
	if pos+3*m > len(fields) {
		return "", fmt.Errorf("not enough edge data")
	}
	g := make([][]Edge, n+1)
	for i := 0; i < m; i++ {
		u, err := nextInt()
		if err != nil {
			return "", fmt.Errorf("edge %d bad u: %v", i+1, err)
		}
		v, err := nextInt()
		if err != nil {
			return "", fmt.Errorf("edge %d bad v: %v", i+1, err)
		}
		w, err := nextInt()
		if err != nil {
			return "", fmt.Errorf("edge %d bad w: %v", i+1, err)
		}
		g[u] = append(g[u], Edge{v, w})
		g[v] = append(g[v], Edge{u, w})
	}

	dist := make([][]int64, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = make([]int64, 51)
		for j := 0; j <= 50; j++ {
			dist[i][j] = INF
		}
	}
	dist[1][0] = 0
	pq := &PriorityQueue{}
	heap.Push(pq, &Item{d: 0, node: 1, last: 0})
	for pq.Len() > 0 {
		cur := heap.Pop(pq).(*Item)
		if cur.d != dist[cur.node][cur.last] {
			continue
		}
		if cur.last == 0 {
			for _, e := range g[cur.node] {
				if dist[e.to][e.w] > cur.d {
					dist[e.to][e.w] = cur.d
					heap.Push(pq, &Item{d: cur.d, node: e.to, last: e.w})
				}
			}
		} else {
			for _, e := range g[cur.node] {
				nd := cur.d + int64(cur.last+e.w)*int64(cur.last+e.w)
				if dist[e.to][0] > nd {
					dist[e.to][0] = nd
					heap.Push(pq, &Item{d: nd, node: e.to, last: 0})
				}
			}
		}
	}

	out := make([]string, n)
	for i := 1; i <= n; i++ {
		if dist[i][0] == INF {
			out[i-1] = "-1"
		} else {
			out[i-1] = strconv.FormatInt(dist[i][0], 10)
		}
	}
	return strings.Join(out, " "), nil
}

type testCase struct {
	input string
}

// Embedded testcases from testcasesE.txt.
const testcaseData = `
3 3 1 3 6 1 3 1 2 3 10
5 5 1 3 5 2 5 1 2 5 7 2 5 4 4 5 9
5 7 4 5 7 1 3 1 3 5 9 4 5 10 2 4 4 4 5 3 3 4 5
3 2 1 3 6 1 3 5
6 2 4 6 8 2 6 10
2 1 1 2 3
2 1 1 2 10
4 5 2 4 5 1 3 2 1 4 10 2 3 1 2 3 7
6 5 3 5 9 1 3 7 3 4 5 4 6 4 3 5 1
5 3 1 5 1 3 5 8 1 3 8
2 1 1 2 3
4 3 2 3 9 2 3 4 1 4 1
4 2 2 4 3 1 4 2
2 1 1 2 5
3 1 1 3 4
4 5 2 4 5 3 4 6 2 3 6 1 3 10 3 4 7
6 7 4 5 4 2 3 2 3 5 9 1 4 3 3 5 2 1 6 8 1 2 8
3 1 1 3 9
3 2 2 3 2 1 3 1
3 3 1 2 1 1 2 3 1 2 10
4 5 1 2 1 1 4 6 1 3 4 3 4 1 2 3 1
2 1 1 2 1
2 1 1 2 6
2 1 1 2 6
5 2 1 4 10 2 5 1
3 1 2 3 1
6 7 1 4 4 1 3 8 3 6 7 4 6 8 2 4 1 1 2 4 3 5 10
4 3 1 3 8 1 3 10 1 4 2
2 1 1 2 9
2 1 1 2 1
4 1 3 4 10
4 1 1 2 2
6 2 1 5 9 3 5 3
2 1 1 2 7
6 4 5 6 9 4 6 3 4 6 10 1 2 8
5 8 3 4 10 2 3 3 2 3 9 1 3 8 2 3 5 2 4 3 2 5 3 3 5 4
4 4 2 4 10 1 4 8 2 4 8 1 4 2
3 1 1 3 4
4 2 2 3 3 1 2 5
3 1 1 2 7
2 1 1 2 8
6 6 1 4 4 1 3 5 5 6 8 1 4 9 2 4 9 1 6 7
4 1 3 4 5
2 1 1 2 9
4 1 2 3 3
3 2 2 3 2 1 3 3
4 6 3 4 3 2 3 2 1 2 9 3 4 2 1 4 3 1 4 9
4 2 1 3 8 1 2 5
6 2 1 6 10 1 6 3
3 3 2 3 7 1 3 7 1 3 10
3 3 1 2 8 2 3 10 1 3 10
4 6 1 2 1 2 3 3 3 4 8 2 4 4 1 4 2 1 2 8
5 2 1 5 1 3 5 4
6 2 1 2 10 4 5 1
2 1 1 2 2
3 2 1 2 1 1 3 8
5 2 1 4 9 2 4 8
4 1 2 4 1
4 3 3 4 9 3 4 4 1 2 5
4 2 1 2 3 2 4 3
6 7 1 2 1 1 2 7 1 2 10 2 4 8 2 5 1 3 4 1 1 4 8
2 1 1 2 6
2 1 1 2 9
4 1 1 3 1
5 2 1 2 9 1 4 1
3 3 2 3 7 1 3 7 2 3 4
2 1 1 2 2
6 7 4 6 6 2 3 9 2 4 8 3 5 9 4 5 6 2 3 7 4 6 3
6 7 2 6 1 1 6 10 1 4 3 1 4 2 2 3 10 2 4 2 1 4 8
4 5 1 3 8 1 2 3 2 3 4 2 4 9 1 2 4
6 6 2 5 2 2 6 10 1 4 7 2 5 8 2 3 2 1 5 2
2 1 1 2 5
6 8 2 4 8 4 6 2 2 4 2 1 3 4 3 6 2 1 3 6 1 5 8 5 6 6
4 2 1 3 1 2 3 10
3 1 1 3 3
4 5 2 3 3 1 3 8 1 4 1 1 4 5 3 4 7
5 8 1 3 8 2 3 9 1 2 3 2 3 5 3 4 8 2 3 7 2 5 3 3 4 7
6 1 1 2 1
3 2 1 2 7 1 3 8
4 1 1 4 5
4 2 1 4 3 3 4 9
3 2 2 3 1 1 2 3
2 1 1 2 6
5 8 3 4 9 3 5 10 2 3 5 2 5 1 2 4 10 3 5 1 1 5 8 2 5 6
2 1 1 2 1
6 8 2 6 7 1 4 4 4 5 9 4 6 4 3 4 1 1 3 3 2 3 1 3 4 7
3 3 1 3 2 1 3 9 1 3 5
2 1 1 2 3
3 2 2 3 7 1 3 3
2 1 1 2 3
3 2 1 3 9 1 3 3
3 2 1 2 7 1 2 3
2 1 1 2 7
4 5 1 2 1 1 2 6 1 4 6 3 4 8 2 3 8
3 3 1 3 1 1 2 3 1 3 5
5 3 1 2 7 1 5 1 3 5 9
4 1 2 3 10
6 4 3 5 6 3 4 9 3 6 6 3 5 10
4 3 2 3 2 3 4 5 2 3 4
5 3 3 4 9 2 5 9 4 5 10
`

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	res := make([]testCase, 0, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			return nil, fmt.Errorf("line %d: not enough data", i+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad n: %v", i+1, err)
		}
		m, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad m: %v", i+1, err)
		}
		if len(fields) != 2+3*m {
			return nil, fmt.Errorf("line %d: expected %d edge tokens, got %d", i+1, 2+3*m, len(fields))
		}
		var input bytes.Buffer
		fmt.Fprintln(&input, n, m)
		for j := 0; j < m; j++ {
			u := fields[2+3*j]
			v := fields[2+3*j+1]
			w := fields[2+3*j+2]
			fmt.Fprintf(&input, "%s %s %s\n", u, v, w)
		}
		res = append(res, testCase{input: input.String()})
	}
	return res, nil
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		expected, err := referenceSolve(tc.input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := runCandidate(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx+1, tc.input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
