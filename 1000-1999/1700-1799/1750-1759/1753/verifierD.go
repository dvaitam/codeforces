package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Edge struct {
	to   int
	cost int64
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
	x := old[n-1]
	*pq = old[:n-1]
	return x
}

type test struct {
	input    string
	expected string
}

func solveCase(n, m int, pCost, qCost int64, grid []string) string {
	idx := func(r, c int) int { return r*m + c }
	valid := func(r, c int) bool {
		return r >= 0 && r < n && c >= 0 && c < m && grid[r][c] != '#'
	}
	adj := make([][]Edge, n*m)
	dotCount := 0
	addEdge := func(fr, fc, toR, toC int, cost int64) {
		if !valid(fr, fc) || !valid(toR, toC) {
			return
		}
		from := idx(fr, fc)
		to := idx(toR, toC)
		adj[from] = append(adj[from], Edge{to: to, cost: cost})
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			switch grid[i][j] {
			case '.':
				dotCount++
			case 'L':
				addEdge(i, j-1, i, j+1, qCost)
				addEdge(i, j+2, i, j, qCost)
				addEdge(i-1, j, i, j+1, pCost)
				addEdge(i+1, j, i, j+1, pCost)
				addEdge(i-1, j+1, i, j, pCost)
				addEdge(i+1, j+1, i, j, pCost)
			case 'U':
				addEdge(i-1, j, i+1, j, qCost)
				addEdge(i+2, j, i, j, qCost)
				addEdge(i, j-1, i+1, j, pCost)
				addEdge(i, j+1, i+1, j, pCost)
				addEdge(i+1, j-1, i, j, pCost)
				addEdge(i+1, j+1, i, j, pCost)
			}
		}
	}
	if dotCount < 2 {
		return "-1"
	}
	const INF int64 = 1 << 60
	dist := make([]int64, n*m)
	for i := range dist {
		dist[i] = INF
	}
	pq := &PriorityQueue{}
	heap.Init(pq)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '.' {
				id := idx(i, j)
				dist[id] = 0
				heap.Push(pq, Item{node: id, dist: 0})
			}
		}
	}
	for pq.Len() > 0 {
		it := heap.Pop(pq).(Item)
		if it.dist != dist[it.node] {
			continue
		}
		v := it.node
		for _, e := range adj[v] {
			nd := it.dist + e.cost
			if nd < dist[e.to] {
				dist[e.to] = nd
				heap.Push(pq, Item{node: e.to, dist: nd})
			}
		}
	}
	ans := INF
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			id := idx(i, j)
			if j+1 < m {
				id2 := idx(i, j+1)
				sum := dist[id] + dist[id2]
				if sum < ans {
					ans = sum
				}
			}
			if i+1 < n {
				id2 := idx(i+1, j)
				sum := dist[id] + dist[id2]
				if sum < ans {
					ans = sum
				}
			}
		}
	}
	if ans >= INF {
		return "-1"
	}
	return fmt.Sprintf("%d", ans)
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(45))
	var tests []test
	base := []string{"..", ".."}
	tests = append(tests, test{input: "2 2\n1 1\n..\n..\n", expected: solveCase(2, 2, 1, 1, base)})
	for len(tests) < 100 {
		n := rng.Intn(3) + 1
		m := rng.Intn(3) + 1
		p := int64(rng.Intn(3) + 1)
		q := int64(rng.Intn(3) + 1)
		grid := make([]string, n)
		for i := 0; i < n; i++ {
			var sb strings.Builder
			for j := 0; j < m; j++ {
				chars := ".#LU"
				sb.WriteByte(chars[rng.Intn(len(chars))])
			}
			grid[i] = sb.String()
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		sb.WriteString(fmt.Sprintf("%d %d\n", p, q))
		for i := 0; i < n; i++ {
			sb.WriteString(grid[i])
			sb.WriteByte('\n')
		}
		tests = append(tests, test{input: sb.String(), expected: solveCase(n, m, p, q, grid)})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		out, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(t.expected) {
			fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected:%s\nGot:%s\n", i+1, t.input, t.expected, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
