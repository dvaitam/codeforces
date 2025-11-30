package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const testcaseData = `
100
5 4
0
2
5 4
4 3
3
3 5
2 5
2 3
1
1 5
3
2 3 1
2 3
1
1 2
1
2 1
1
1 2
1
3
5 3
1
3 1
1
5 2
1
2 5
4
1 1 2 3
5 1
2
5 3
1 5
3
1 1 1
2 4
1
1 2
0
0
0
1
3
5 1
0
2
1 1
2 4
1
2 1
1
1 2
1
1 2
0
2
3 2
4 4
0
0
1
2 1
4
1 2
1 4
1 3
4 2
1
1
5 1
2
1 2
1 3
3
1 1 1
5 1
4
1 4
2 3
3 4
5 2
2
1 1
3 3
2
1 3
2 3
1
1 2
3
2 3
1 3
1 2
1
3
4 2
1
4 3
4
2 3
2 4
1 4
1 3
2
2 2
2 1
1
2 1
4
1 1 1 1
5 4
2
1 2
1 5
1
2 1
3
4 1
1 3
3 2
4
3 1
2 3
1 2
3 5
3
3 1 4
5 4
2
5 2
2 4
4
3 1
5 1
5 3
3 4
4
2 3
1 4
2 1
5 1
2
4 3
5 4
4
3 3 1 4
2 4
1
1 2
1
2 1
0
0
2
2 1
5 1
0
4
1 1 1 1
5 1
2
2 3
5 2
4
1 1 1 1
2 4
0
0
1
2 1
0
5
4 2 4 4 4
4 4
3
1 3
1 2
3 2
3
3 4
1 3
1 4
2
4 3
4 1
4
4 2
3 4
2 1
3 1
4
3 3 4 4
5 1
1
2 4
5
1 1 1 1 1
5 4
3
4 2
2 3
1 5
2
2 4
4 1
0
4
4 5
3 1
5 3
2 4
3
1 2 1
5 4
2
2 4
3 1
0
0
3
3 1
5 2
2 1
5
4 3 1 2 4
5 1
0
4
1 1 1 1
2 4
1
2 1
1
1 2
0
0
3
1 3 3
3 1
1
2 1
5
1 1 1 1 1
2 3
1
2 1
1
1 2
1
1 2
2
2 1
3 2
2
3 1
2 3
0
2
1 2
2 1
0
4
1 1 1 1
4 3
0
1
4 3
4
1 4
4 3
4 2
2 1
3
1 3 3
3 4
3
1 2
2 3
3 1
3
1 2
1 3
2 3
3
2 1
3 2
1 3
3
1 3
1 2
2 3
5
4 4 1 3 3
4 2
4
1 3
1 2
4 3
3 2
4
2 1
4 3
4 2
4 1
4
2 1 2 1
3 1
2
1 3
2 1
3
1 1 1
4 3
1
2 3
1
1 4
3
3 4
2 3
2 4
2
3 2
3 2
1
3 1
1
2 1
5
1 2 1 1 1
3 3
0
1
1 3
2
3 1
2 3
4
2 2 3 2
3 3
2
2 1
2 3
0
1
1 2
2
3 2
3 4
1
2 3
0
0
2
2 3
2 1
5
1 2 3 1 2
4 2
1
3 4
2
4 3
2 1
3
1 1 2
5 4
0
3
4 1
2 1
2 4
1
4 5
0
4
4 1 2 2
5 2
1
3 4
0
5
2 1 2 2 2
4 3
1
1 2
3
3 1
4 1
3 2
4
3 2
2 1
4 3
1 3
3
3 3 1
3 4
1
1 3
0
1
1 2
1
1 3
2
2 4
4 2
3
3 4
1 2
2 4
3
3 1
3 2
4 2
3
2 2 1
3 1
2
1 2
2 3
2
1 1
3 3
0
2
1 3
2 1
3
3 1
2 1
3 2
2
1 2
5 2
2
3 2
3 4
3
1 4
3 5
5 4
1
1
5 4
1
1 5
1
2 1
2
2 1
2 5
1
1 4
5
1 4 2 1 3
4 4
0
0
3
1 3
3 2
4 2
4
3 2
4 3
4 1
3 1
4
4 1 4 3
5 1
0
3
1 1 1
3 3
3
2 1
1 3
3 2
1
1 2
0
5
1 3 1 2 2
4 3
2
4 2
1 2
4
1 4
2 3
2 1
3 1
3
2 4
2 1
3 4
1
3
5 2
0
2
5 2
2 3
4
1 1 1 2
3 4
0
3
1 3
3 2
1 2
0
3
3 2
1 2
3 1
5
4 3 2 3 3
2 4
1
1 2
1
2 1
1
1 2
1
1 2
4
2 2 4 3
2 2
0
0
4
1 1 2 1
4 1
0
5
1 1 1 1 1
4 3
4
4 2
1 2
1 3
3 2
3
3 1
4 3
1 2
2
1 2
3 2
1
2
2 3
1
2 1
1
1 2
1
1 2
1
2
2 1
0
3
1 1 1
2 3
1
1 2
1
1 2
0
2
3 3
4 1
0
4
1 1 1 1
3 2
1
3 2
0
4
2 1 1 1
4 3
1
2 3
1
2 4
3
3 2
4 1
2 1
4
3 1 1 1
2 3
0
0
0
1
3
3 2
3
1 2
2 3
3 1
1
2 3
2
2 1
3 4
2
1 2
1 3
3
2 3
1 2
1 3
3
1 2
1 3
2 3
3
2 1
3 2
3 1
3
3 2 3
5 1
1
3 2
5
1 1 1 1 1
5 1
3
3 4
2 3
1 5
3
1 1 1
5 4
0
0
3
1 5
2 4
3 5
2
3 1
4 3
2
3 3
3 4
2
1 2
2 3
2
1 3
2 3
1
2 1
0
1
1
4 4
2
1 4
4 2
1
1 2
0
1
4 2
4
4 3 1 2
5 2
3
5 1
4 1
5 4
2
2 5
4 1
4
1 2 2 2
5 3
1
2 5
2
2 3
5 2
3
5 4
5 1
3 4
2
1 2
3 3
0
1
3 1
1
1 2
5
3 3 2 3 3
2 4
1
2 1
0
0
0
4
3 3 2 1
5 3
4
3 1
5 4
2 3
4 2
1
3 2
3
4 5
2 1
4 3
5
2 1 3 3 2
4 4
2
3 4
3 1
3
3 1
4 2
3 2
1
3 4
1
1 4
4
4 4 2 4
4 4
1
3 1
0
3
3 4
2 4
4 1
2
2 3
2 1
2
3 1
3 1
0
4
1 1 1 1
2 4
0
0
1
2 1
1
1 2
1
3
2 3
1
2 1
0
1
1 2
2
3 3
4 2
3
4 1
2 4
3 2
3
3 4
1 3
2 4
2
2 2
5 1
3
1 5
2 4
2 3
4
1 1 1 1
4 3
3
2 1
2 3
1 3
4
2 4
4 1
3 2
1 2
1
4 2
1
3
5 3
2
4 2
2 3
4
3 4
1 2
2 4
4 5
1
3 1
1
2
3 1
3
2 3
1 3
2 1
2
1 1
4 1
0
2
1 1
3 4
2
1 3
1 2
3
2 1
3 2
1 3
3
2 3
3 1
2 1
0
3
3 4 1
4 3
2
2 1
3 1
1
4 1
2
3 2
3 1
5
3 2 1 2 3
2 1
0
4
1 1 1 1
3 2
1
2 1
3
2 1
2 3
1 3
1
2
2 2
1
1 2
1
1 2
5
2 1 1 2 2
3 3
1
3 2
2
1 3
2 1
2
2 3
1 3
2
3 3
3 3
0
3
3 1
3 2
1 2
2
1 3
3 2
4
1 1 1 2
2 2
0
1
1 2
5
2 1 1 2 2
4 4
3
4 1
2 3
1 2
1
4 2
1
2 3
0
4
2 1 4 2
5 4
1
4 5
4
4 1
2 1
5 1
2 4
3
3 4
3 1
2 5
1
2 1
1
1
5 2
3
2 5
5 3
1 5
1
5 3
2
1 2
3 4
2
2 1
3 2
2
3 2
3 1
1
2 1
2
3 1
3 2
3
3 4 3
`

type edge struct {
	to, moment int
}

type item struct {
	step, city int
}

type minHeap []item

func (h minHeap) Len() int           { return len(h) }
func (h minHeap) Less(i, j int) bool { return h[i].step < h[j].step }
func (h minHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) {
	*h = append(*h, x.(item))
}
func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

type testCase struct {
	input   string
	n       int
	m       int
	adj     [][]edge
	visits  [][]int
	maxTime int
}

func solve(tc testCase) int {
	const INF = int(1e9)
	dist := make([]int, tc.n+1)
	for i := 1; i <= tc.n; i++ {
		dist[i] = INF
	}
	dist[1] = 0
	pq := &minHeap{{0, 1}}
	heap.Init(pq)
	for pq.Len() > 0 {
		it := heap.Pop(pq).(item)
		if it.step != dist[it.city] {
			continue
		}
		if it.city == tc.n {
			return it.step
		}
		if it.step >= len(tc.visits) {
			continue
		}
		for _, e := range tc.adj[it.city] {
			arr := tc.visits[e.moment]
			idx := sort.Search(len(arr), func(i int) bool { return arr[i] > it.step })
			if idx < len(arr) {
				step := arr[idx]
				if step < dist[e.to] {
					dist[e.to] = step
					heap.Push(pq, item{step, e.to})
				}
			}
		}
	}
	return -1
}

func loadCases() ([]testCase, error) {
	tokens := strings.Fields(testcaseData)
	if len(tokens) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	pos := 0
	t, err := strconv.Atoi(tokens[pos])
	if err != nil {
		return nil, fmt.Errorf("bad test count: %w", err)
	}
	pos++
	cases := make([]testCase, 0, t)
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if pos+2 > len(tokens) {
			return nil, fmt.Errorf("case %d: missing n/m", caseIdx+1)
		}
		n, _ := strconv.Atoi(tokens[pos])
		m, _ := strconv.Atoi(tokens[pos+1])
		pos += 2

		adj := make([][]edge, n+1)
		counts := make([]int, m+1)
		edgesByMoment := make([][][2]int, m+1)
		for moment := 1; moment <= m; moment++ {
			if pos >= len(tokens) {
				return nil, fmt.Errorf("case %d: missing cnt at moment %d", caseIdx+1, moment)
			}
			cnt, _ := strconv.Atoi(tokens[pos])
			pos++
			counts[moment] = cnt
			es := make([][2]int, cnt)
			for i := 0; i < cnt; i++ {
				if pos+1 >= len(tokens) {
					return nil, fmt.Errorf("case %d: missing edge at moment %d", caseIdx+1, moment)
				}
				u, _ := strconv.Atoi(tokens[pos])
				v, _ := strconv.Atoi(tokens[pos+1])
				pos += 2
				es[i] = [2]int{u, v}
				adj[u] = append(adj[u], edge{to: v, moment: moment})
				adj[v] = append(adj[v], edge{to: u, moment: moment})
			}
			edgesByMoment[moment] = es
		}
		if pos >= len(tokens) {
			return nil, fmt.Errorf("case %d: missing k", caseIdx+1)
		}
		k, _ := strconv.Atoi(tokens[pos])
		pos++
		times := make([]int, k)
		maxTime := 0
		for i := 0; i < k; i++ {
			if pos >= len(tokens) {
				return nil, fmt.Errorf("case %d: missing time value", caseIdx+1)
			}
			times[i], _ = strconv.Atoi(tokens[pos])
			pos++
			if times[i] > maxTime {
				maxTime = times[i]
			}
		}

		if maxTime < m {
			maxTime = m
		}
		visits := make([][]int, maxTime+1)
		for idx, tval := range times {
			if tval >= len(visits) {
				newVis := make([][]int, tval+1)
				copy(newVis, visits)
				visits = newVis
			}
			visits[tval] = append(visits[tval], idx+1)
		}

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for moment := 1; moment <= m; moment++ {
			sb.WriteString(fmt.Sprintf("%d\n", counts[moment]))
			for _, e := range edgesByMoment[moment] {
				sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
			}
		}
		sb.WriteString(fmt.Sprintf("%d\n", k))
		for i, v := range times {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')

		cases = append(cases, testCase{
			input:   sb.String(),
			n:       n,
			m:       m,
			adj:     adj,
			visits:  visits,
			maxTime: maxTime,
		})
	}
	if pos != len(tokens) {
		return nil, fmt.Errorf("unused tokens remain")
	}
	return cases, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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
		fmt.Println("usage: verifierB /path/to/binary")
		os.Exit(1)
	}
	cases, err := loadCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load cases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		expected := strconv.Itoa(solve(tc))
		got, err := run(os.Args[1], tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", idx+1, expected, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
