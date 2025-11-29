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

// Embedded copy of testcasesD.txt so the verifier is self-contained.
const testcasesRaw = `1 19 17 2 3 15
0 20
0 16
2 18 7 2 4 7 12 2 2 10
1 17 12 1 1 2
0 10
0 9
3 20 12 4 4 11 18 4 2 11 3 1 2 3
1 9 6 3 4 4
3 19 11 5 5 6 18 2 3 0 8 5 2 5
0 7
2 10 1 1 4 1 1 3 1 1
1 1 1 4 4 0
0 20
0 13
2 18 8 5 2 0 9 1 1 1
0 7
3 10 9 3 2 0 5 3 3 1 6 4 4 6
3 20 17 1 5 16 8 4 2 4 13 3 5 4
2 1 1 5 3 0 1 5 5 0
0 11
3 12 10 3 5 4 11 4 1 9 0 1 3 0
3 10 9 5 3 2 5 2 3 2 9 3 3 6
0 1
1 10 8 2 3 3
2 6 5 4 1 0 4 3 3 1
3 6 0 3 2 0 2 2 1 0 4 2 3 4
1 9 5 1 5 2
1 14 4 5 3 3
2 14 4 4 5 3 0 4 2 0
0 16
3 18 7 1 4 4 17 3 2 2 18 3 1 7
0 2
1 14 9 1 1 7
0 6
2 8 0 5 5 0 0 5 1 0
1 9 8 4 1 5
1 7 1 5 1 0
1 9 2 1 4 2
3 2 1 2 3 1 0 4 3 0 0 2 1 0
0 3
3 2 2 1 5 2 1 3 2 1 0 3 4 0
2 12 4 2 3 3 1 2 5 0
3 3 1 1 3 1 3 1 5 3 0 3 4 0
3 14 7 1 2 3 8 3 5 1 12 4 2 6
1 1 1 3 5 1
0 15
0 17
3 4 2 5 5 0 4 1 4 1 1 4 1 0
0 13
1 1 1 1 1 0
3 10 9 3 1 0 9 5 5 3 1 5 1 0
2 19 5 1 2 1 7 4 5 6
2 12 9 4 3 8 6 1 4 4
1 14 13 2 4 11
3 5 5 4 2 1 0 4 4 0 4 2 2 2
1 5 4 5 3 1
2 14 9 5 5 4 14 2 3 0
2 16 12 2 2 9 11 2 3 7
1 14 11 4 5 3
3 19 17 1 4 2 12 1 4 3 7 1 2 4
1 7 4 2 2 4
0 9
1 2 1 2 4 0
0 4
0 9
2 2 1 4 5 1 0 1 3 0
3 13 7 1 2 7 6 2 5 2 1 3 1 1
0 15
2 4 4 3 3 3 2 3 1 1
0 16
2 2 2 3 5 2 0 2 2 0
3 4 0 5 2 0 4 4 5 2 1 4 4 1
1 3 0 2 5 0
2 4 2 3 4 0 1 1 4 1
1 17 11 3 4 7
0 12
3 4 1 3 5 0 0 5 1 0 1 5 4 1
1 19 19 2 4 6
1 19 5 2 3 2
2 1 1 4 4 1 1 4 5 1
3 1 0 1 1 0 1 2 5 1 0 2 5 0
0 17
3 4 4 3 2 1 3 1 5 0 0 3 5 0
0 16
0 11
2 11 5 2 1 4 0 1 3 0
0 7
3 8 7 3 1 0 6 1 2 5 2 4 4 1
0 18
3 7 7 3 1 7 7 4 4 2 7 1 3 5
2 15 11 5 4 3 0 2 3 0
1 15 6 2 2 0
1 19 12 5 2 10
0 5
0 20
1 15 15 2 1 0
3 15 10 4 1 0 7 4 1 6 15 1 2 7
0 13
3 7 2 3 5 0 5 1 5 0 4 3 4 2`

type Edge struct {
	to   int
	cost int64
	ramp int
}

type Item struct {
	node int
	dist int64
}

type PriorityQueue []Item

func (pq PriorityQueue) Len() int            { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool  { return pq[i].dist < pq[j].dist }
func (pq PriorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

type Ramp struct {
	x, d, t, p int64
	idx        int
}

type testCase struct {
	n     int
	L     int64
	ramps []Ramp
}

func unique(a []int64) []int64 {
	j := 0
	for i := 0; i < len(a); i++ {
		if i == 0 || a[i] != a[i-1] {
			a[j] = a[i]
			j++
		}
	}
	return a[:j]
}

func solveCase(tc testCase) string {
	coords := make([]int64, 0, 3*tc.n+2)
	coords = append(coords, 0, tc.L)
	usable := make([]Ramp, 0, tc.n)
	for _, r := range tc.ramps {
		if r.x-r.p < 0 || r.x+r.d > tc.L {
			continue
		}
		usable = append(usable, r)
		coords = append(coords, r.x-r.p, r.x, r.x+r.d)
	}
	sort.Slice(coords, func(i, j int) bool { return coords[i] < coords[j] })
	coords = unique(coords)
	idxOf := func(v int64) int {
		return sort.Search(len(coords), func(i int) bool { return coords[i] >= v })
	}
	g := make([][]Edge, len(coords))
	for i := 0; i+1 < len(coords); i++ {
		w := coords[i+1] - coords[i]
		g[i] = append(g[i], Edge{i + 1, w, 0})
		g[i+1] = append(g[i+1], Edge{i, w, 0})
	}
	for _, r := range usable {
		u := idxOf(r.x - r.p)
		v := idxOf(r.x + r.d)
		cost := r.p + r.t
		g[u] = append(g[u], Edge{v, cost, r.idx})
	}

	const INF = int64(4e18)
	dist := make([]int64, len(coords))
	prevNode := make([]int, len(coords))
	prevRamp := make([]int, len(coords))
	for i := range dist {
		dist[i] = INF
		prevNode[i] = -1
	}
	src := idxOf(0)
	dst := idxOf(tc.L)
	dist[src] = 0
	pq := &PriorityQueue{{src, 0}}
	heap.Init(pq)
	for pq.Len() > 0 {
		it := heap.Pop(pq).(Item)
		u := it.node
		if it.dist != dist[u] {
			continue
		}
		if u == dst {
			break
		}
		for _, e := range g[u] {
			nd := it.dist + e.cost
			if nd < dist[e.to] {
				dist[e.to] = nd
				prevNode[e.to] = u
				prevRamp[e.to] = e.ramp
				heap.Push(pq, Item{e.to, nd})
			}
		}
	}
	path := make([]int, 0)
	for u := dst; u != src; u = prevNode[u] {
		if u < 0 {
			break
		}
		if prevRamp[u] != 0 {
			path = append(path, prevRamp[u])
		}
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", dist[dst]))
	sb.WriteString(fmt.Sprintf("%d\n", len(path)))
	if len(path) > 0 {
		for i, v := range path {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
	}
	return strings.TrimSpace(sb.String())
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(testcasesRaw, "\n")
	var cases []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid n: %v", idx+1, err)
		}
		if len(parts) != 2+4*n {
			return nil, fmt.Errorf("line %d: expected %d fields got %d", idx+1, 2+4*n, len(parts))
		}
		L, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid L: %v", idx+1, err)
		}
		ramps := make([]Ramp, n)
		pos := 2
		for i := 0; i < n; i++ {
			x, err1 := strconv.ParseInt(parts[pos], 10, 64)
			d, err2 := strconv.ParseInt(parts[pos+1], 10, 64)
			t, err3 := strconv.ParseInt(parts[pos+2], 10, 64)
			p, err4 := strconv.ParseInt(parts[pos+3], 10, 64)
			if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
				return nil, fmt.Errorf("line %d: parse ramp %d: %v %v %v %v", idx+1, i+1, err1, err2, err3, err4)
			}
			ramps[i] = Ramp{x: x, d: d, t: t, p: p, idx: i + 1}
			pos += 4
		}
		cases = append(cases, testCase{n: n, L: L, ramps: ramps})
	}
	if len(cases) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
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
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
		var input strings.Builder
		fmt.Fprintf(&input, "%d %d\n", tc.n, tc.L)
		for _, r := range tc.ramps {
			fmt.Fprintf(&input, "%d %d %d %d\n", r.x, r.d, r.t, r.p)
		}

		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Printf("test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		expect := solveCase(tc)
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed\nexpected:\n%s\n\ngot:\n%s\n", i+1, expect, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(cases))
}
