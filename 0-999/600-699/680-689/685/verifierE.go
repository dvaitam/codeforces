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
	to  int
	idx int
}

type Item struct {
	node int
	time int
}

type PQ []Item

func (pq PQ) Len() int            { return len(pq) }
func (pq PQ) Less(i, j int) bool  { return pq[i].time < pq[j].time }
func (pq PQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PQ) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PQ) Pop() interface{} {
	old := *pq
	n := len(old)
	it := old[n-1]
	*pq = old[:n-1]
	return it
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
	return out.String(), err
}

func reachable(adj [][]Edge, l, r, s, t int) bool {
	const INF = int(1<<31 - 1)
	dist := make([]int, len(adj))
	for i := range dist {
		dist[i] = INF
	}
	pq := &PQ{}
	heap.Init(pq)
	dist[s] = l
	heap.Push(pq, Item{s, l})
	for pq.Len() > 0 {
		it := heap.Pop(pq).(Item)
		if it.time != dist[it.node] {
			continue
		}
		if it.time > r {
			continue
		}
		if it.node == t {
			break
		}
		for _, e := range adj[it.node] {
			if e.idx < l || e.idx > r || e.idx < it.time {
				continue
			}
			if e.idx < dist[e.to] {
				dist[e.to] = e.idx
				heap.Push(pq, Item{e.to, e.idx})
			}
		}
	}
	return dist[t] <= r
}

func buildGraph(n, m int) ([][]Edge, [][2]int) {
	g := make([][]Edge, n+1)
	edges := make([][2]int, m)
	for i := 0; i < m; i++ {
		v := rand.Intn(n) + 1
		u := rand.Intn(n) + 1
		for u == v {
			u = rand.Intn(n) + 1
		}
		idx := i + 1
		g[v] = append(g[v], Edge{u, idx})
		g[u] = append(g[u], Edge{v, idx})
		edges[i] = [2]int{v, u}
	}
	return g, edges
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go <binary>")
		return
	}
	bin := os.Args[1]
	rand.Seed(42)
	for tcase := 0; tcase < 100; tcase++ {
		n := rand.Intn(4) + 2
		m := rand.Intn(6) + 1
		q := rand.Intn(4) + 1
		adj, edges := buildGraph(n, m)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, q))
		for i := 0; i < m; i++ {
			sb.WriteString(fmt.Sprintf("%d %d\n", edges[i][0], edges[i][1]))
		}
		answers := make([]string, q)
		for i := 0; i < q; i++ {
			l := rand.Intn(m) + 1
			r := rand.Intn(m-l+1) + l
			s := rand.Intn(n) + 1
			t := rand.Intn(n) + 1
			for t == s {
				t = rand.Intn(n) + 1
			}
			sb.WriteString(fmt.Sprintf("%d %d %d %d\n", l, r, s, t))
			if reachable(adj, l, r, s, t) {
				answers[i] = "Yes"
			} else {
				answers[i] = "No"
			}
		}
		out, err := runBinary(bin, sb.String())
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\noutput:%s", tcase+1, err, out)
			return
		}
		outLines := strings.Fields(out)
		if len(outLines) != q {
			fmt.Printf("invalid output on test %d\ninput:%soutput:%s", tcase+1, sb.String(), out)
			return
		}
		for i := 0; i < q; i++ {
			if outLines[i] != answers[i] {
				fmt.Printf("wrong answer on test %d\ninput:%sexpected:%v\noutput:%s", tcase+1, sb.String(), answers, out)
				return
			}
		}
	}
	fmt.Println("All tests passed")
}
