package main

import (
	"bytes"
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
const maxLog = 20

func referenceSolve(n int, edges [][2]int, queries [][2]int) (string, error) {
	adj := make([][]Edge, n)
	for _, e := range edges {
		u, v := e[0]-1, e[1]-1
		adj[u] = append(adj[u], Edge{v, 0})
		adj[v] = append(adj[v], Edge{u, 0})
	}

	depth := make([]int, n)
	parent := make([]int, n)
	var up [maxLog][]int
	for i := 0; i < maxLog; i++ {
		up[i] = make([]int, n)
	}

	// DFS stack to set parents/depth and order.
	stack := []int{0}
	parent[0] = 0
	order := make([]int, 0, n)
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, v)
		for _, to := range adj[v] {
			if to.to == parent[v] {
				continue
			}
			parent[to.to] = v
			depth[to.to] = depth[v] + 1
			stack = append(stack, to.to)
		}
	}
	up[0] = parent
	for k := 1; k < maxLog; k++ {
		for i := 0; i < n; i++ {
			up[k][i] = up[k-1][up[k-1][i]]
		}
	}

	lca := func(u, v int) int {
		if depth[u] < depth[v] {
			u, v = v, u
		}
		for k := maxLog - 1; k >= 0; k-- {
			if depth[u]-(1<<k) >= depth[v] {
				u = up[k][u]
			}
		}
		if u == v {
			return u
		}
		for k := maxLog - 1; k >= 0; k-- {
			if up[k][u] != up[k][v] {
				u = up[k][u]
				v = up[k][v]
			}
		}
		return up[0][u]
	}

	jump := func(u, k int) int {
		i := 0
		for k > 0 {
			if k&1 == 1 {
				u = up[i][u]
			}
			k >>= 1
			i++
		}
		return u
	}

	childOnPath := func(u, l int) int {
		if u == l {
			return -1
		}
		return jump(u, depth[u]-depth[l]-1)
	}

	nodeAdd := make([]int64, n)
	edgeAdd := make([]int64, n)
	pairCC := make([]map[uint64]int, n)
	sumPairChild := make([]map[int]int, n)
	endPair := make([]map[int]int, n)
	single := make([]int, n)

	for _, q := range queries {
		u, v := q[0]-1, q[1]-1
		w := lca(u, v)
		nodeAdd[u]++
		nodeAdd[v]++
		nodeAdd[w]--
		if w != 0 {
			nodeAdd[parent[w]]--
		}
		edgeAdd[u]++
		edgeAdd[v]++
		edgeAdd[w] -= 2
		a := childOnPath(u, w)
		b := childOnPath(v, w)
		if a == -1 && b == -1 {
			single[w]++
		} else if a == -1 || b == -1 {
			c := a
			if c == -1 {
				c = b
			}
			if endPair[w] == nil {
				endPair[w] = make(map[int]int)
			}
			endPair[w][c]++
		} else {
			if a > b {
				a, b = b, a
			}
			if pairCC[w] == nil {
				pairCC[w] = make(map[uint64]int)
			}
			key := uint64(a)<<32 | uint64(b)
			pairCC[w][key]++
			if sumPairChild[w] == nil {
				sumPairChild[w] = make(map[int]int)
			}
			sumPairChild[w][a]++
			sumPairChild[w][b]++
		}
	}

	for i := len(order) - 1; i > 0; i-- {
		v := order[i]
		p := parent[v]
		nodeAdd[p] += nodeAdd[v]
		edgeAdd[p] += edgeAdd[v]
	}

	nodeCnt := nodeAdd
	edgeCnt := make([]int64, n)
	for i := 1; i < n; i++ {
		edgeCnt[i] = edgeAdd[i]
	}

	var ans int64
	for x := 0; x < n; x++ {
		mcnt := nodeCnt[x]
		totalPairs := mcnt * (mcnt - 1) / 2
		var sumDir int64
		for _, c := range adj[x] {
			if c.to == parent[x] {
				continue
			}
			d := edgeCnt[c.to]
			sumDir += d * (d - 1) / 2
		}
		if x != 0 {
			d := edgeCnt[x]
			sumDir += d * (d - 1) / 2
		}
		var sumPair int64
		if mp := pairCC[x]; mp != nil {
			for _, cnt := range mp {
				c := int64(cnt)
				sumPair += c * (c - 1) / 2
			}
		}
		if x != 0 {
			for _, c := range adj[x] {
				if c.to == parent[x] {
					continue
				}
				cnt := edgeCnt[c.to]
				if mp := endPair[x]; mp != nil {
					cnt -= int64(mp[c.to])
				}
				if mp := sumPairChild[x]; mp != nil {
					cnt -= int64(mp[c.to])
				}
				sumPair += cnt * (cnt - 1) / 2
			}
		}
		ans += totalPairs - sumDir + sumPair
	}
	return strconv.FormatInt(ans, 10), nil
}

type testCase struct {
	input string
}

// Embedded testcases from testcasesF.txt.
const testcaseData = `
2 1 2 1 2 2
2 1 2 1 1 2
5 1 2 1 3 1 4 3 5 3 2 1 3 2 1 3
3 1 2 1 3 3 2 3 3 2 1 3
3 1 2 1 3 2 1 2 2 1
5 1 2 1 3 2 4 3 5 5 2 4 4 5 3 4 4 2 2 3
3 1 2 1 3 1 2 3
3 1 2 2 3 2 3 1 1 2
2 1 2 3 1 2 2 2 1 2
1 1 1 1
5 1 2 1 3 2 4 2 5 3 4 1 1 3 1 3
6 1 2 1 3 2 4 3 5 3 6 2 6 4 5 6
1 3 1 1 1 1 1 1
5 1 2 2 3 3 4 1 5 3 1 4 2 3 3 3
5 1 2 2 3 1 4 4 5 2 1 1 1 1
6 1 2 1 3 3 4 1 5 5 6 4 5 2 3 1 1 5 3 4
6 1 2 2 3 1 4 2 5 4 6 4 4 1 2 4 4 2 6 4
2 1 2 2 1 1 2 2
2 1 2 2 2 2 1 2
1 3 1 1 1 1 1 1
4 1 2 1 3 2 4 3 4 3 4 2 3 3
4 1 2 1 3 2 4 6 2 1 4 2 3 1 2 4 4 4 4 2
1 2 1 1 1 1
4 1 2 1 3 3 4 1 2 2
6 1 2 2 3 3 4 4 5 1 6 1 1 6
5 1 2 2 3 3 4 3 5 5 2 1 3 1 5 1 3 3 1 1
5 1 2 2 3 1 4 3 5 4 2 4 4 1 1 1 5 3
5 1 2 2 3 3 4 4 5 1 2 3
4 1 2 1 3 3 4 2 3 2 2 2
3 1 2 2 3 3 3 1 3 1 1 3
3 1 2 1 3 4 3 3 2 3 1 3 1 2
3 1 2 1 3 2 3 1 3 1
5 1 2 2 3 3 4 3 5 5 4 5 5 4 4 2 3 5 3 3
2 1 2 1 1 1
3 1 2 1 3 5 2 3 3 1 1 2 3 1 2 3
3 1 2 1 3 3 1 3 1 3 1 2
1 2 1 1 1 1
2 1 2 6 1 2 2 1 2 1 2 2 1 2 2 2
5 1 2 2 3 2 4 4 5 4 5 1 3 2 5 4 5 5
3 1 2 2 3 5 1 1 3 1 2 3 2 3 2 1
3 1 2 1 3 6 2 3 1 1 1 1 2 1 1 1 2 3
6 1 2 1 3 3 4 2 5 2 6 1 2 1
3 1 2 1 3 3 1 3 2 3 3 3
2 1 2 6 2 1 1 2 1 1 2 1 1 2 2 1
4 1 2 1 3 2 4 5 4 3 4 3 2 4 1 4 1 3
6 1 2 2 3 1 4 1 5 2 6 1 5 1
2 1 2 5 1 2 1 2 2 2 1 1 1 2
5 1 2 2 3 3 4 4 5 5 1 2 3 5 5 3 4 3 1 2
3 1 2 1 3 5 2 3 2 2 1 1 3 1 2 3
5 1 2 1 3 1 4 1 5 3 5 2 3 1 5 5
1 1 1 1
5 1 2 1 3 1 4 4 5 6 2 2 3 2 3 3 4 5 4 1 2 4
4 1 2 2 3 3 4 6 3 2 4 3 2 2 1 4 4 3 3 4
6 1 2 2 3 3 4 2 5 3 6 3 5 1 1 4 3 5
3 1 2 2 3 1 1 3
1 6 1 1 1 1 1 1 1 1 1 1 1 1
5 1 2 1 3 2 4 3 5 5 5 2 1 2 1 2 2 3 1 5
3 1 2 2 3 2 1 1 1 1
2 1 2 6 2 2 1 1 1 2 2 2 2 2 2 2
1 1 1 1
3 1 2 2 3 2 2 3 1 3
6 1 2 2 3 2 4 1 5 2 6 2 3 2 5 2
1 1 1 1
3 1 2 2 3 1 2 3
4 1 2 1 3 2 4 5 4 3 3 1 4 2 4 1 4 2
1 3 1 1 1 1 1 1
2 1 2 5 2 1 1 2 2 1 1 1 1 1
5 1 2 1 3 2 4 3 5 5 3 2 3 5 5 3 3 5 3 5
1 5 1 1 1 1 1 1 1 1 1 1
5 1 2 1 3 3 4 3 5 1 1 5
2 1 2 3 1 2 2 2 1 2
3 1 2 1 3 3 2 1 1 2 3 3
3 1 2 1 3 4 3 1 2 1 2 1 3 1
2 1 2 6 1 1 1 1 2 2 1 1 2 1 1 1
6 1 2 1 3 1 4 4 5 2 6 1 5 3
5 1 2 2 3 1 4 1 5 5 5 1 5 1 2 1 3 5 4 3
5 1 2 1 3 3 4 3 5 5 4 2 5 4 2 4 5 5 1 5
5 1 2 2 3 3 4 4 5 3 2 3 2 3 3 3
6 1 2 2 3 3 4 3 5 3 6 3 6 2 3 1 2 3
6 1 2 1 3 2 4 1 5 3 6 1 5 1
1 1 1 1
2 1 2 6 2 2 1 2 2 2 1 2 1 1 1 1
4 1 2 2 3 1 4 4 1 2 3 3 3 2 1 3
3 1 2 2 3 5 1 2 2 3 2 3 3 1 2 1
1 4 1 1 1 1 1 1 1 1
4 1 2 2 3 2 4 3 4 3 4 1 3 4
5 1 2 2 3 2 4 3 5 2 5 5 2 4
3 1 2 2 3 5 1 3 3 3 3 1 3 1 2 1
5 1 2 2 3 2 4 4 5 5 1 4 1 4 3 4 3 2 2 3
3 1 2 2 3 5 3 3 3 2 3 1 2 2 3 2
3 1 2 2 3 3 3 3 3 1 3 3
2 1 2 4 2 1 1 1 1 2 1 2
3 1 2 1 3 3 1 3 1 1 2 1
6 1 2 1 3 2 4 1 5 5 6 5 2 2 2 1 2 5 2 2 1 5
4 1 2 2 3 2 4 5 4 2 3 1 3 3 1 1 3 1
4 1 2 1 3 3 4 2 1 4 4 3
1 6 1 1 1 1 1 1 1 1 1 1 1 1
4 1 2 2 3 1 4 6 3 3 1 3 4 3 3 1 2 2 2 4
6 1 2 1 3 3 4 2 5 5 6 1 5 6
1 3 1 1 1 1 1 1
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
		if len(fields) < 2 {
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
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		expected, err := referenceSolveFromInput(tc.input)
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

func referenceSolveFromInput(input string) (string, error) {
	fields := strings.Fields(input)
	if len(fields) < 2 {
		return "", fmt.Errorf("invalid input")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return "", err
	}
	m, err := strconv.Atoi(fields[1])
	if err != nil {
		return "", err
	}
	if len(fields) != 2+3*m {
		return "", fmt.Errorf("expected %d tokens, got %d", 2+3*m, len(fields))
	}
	edges := make([][3]int, m)
	for i := 0; i < m; i++ {
		u, err := strconv.Atoi(fields[2+3*i])
		if err != nil {
			return "", err
		}
		v, err := strconv.Atoi(fields[2+3*i+1])
		if err != nil {
			return "", err
		}
		w, err := strconv.Atoi(fields[2+3*i+2])
		if err != nil {
			return "", err
		}
		edges[i] = [3]int{u, v, w}
	}
	return referenceSolve(n, edges)
}
