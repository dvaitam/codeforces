package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const solveINF = 1000000000

type sEdge struct {
	to     int
	rev    int
	cap    int
	flow   int
	cost   int
	origID int
}

type sOrigEdge struct {
	u, v, c, w int
}

type sItem struct {
	vertex   int
	priority int
}

type sPQ []*sItem

func (pq sPQ) Len() int            { return len(pq) }
func (pq sPQ) Less(i, j int) bool  { return pq[i].priority < pq[j].priority }
func (pq sPQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *sPQ) Push(x interface{}) { *pq = append(*pq, x.(*sItem)) }
func (pq *sPQ) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

func solveF(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return ""
	}

	edges := make([]sOrigEdge, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &edges[i].u, &edges[i].v, &edges[i].c, &edges[i].w)
	}

	nNodes := n + 2
	S := 0
	T := n + 1
	graph := make([][]sEdge, nNodes)
	bal := make([]int, nNodes)

	addEdge := func(u, v, cap, cost, origID int) {
		graph[u] = append(graph[u], sEdge{v, len(graph[v]), cap, 0, cost, origID})
		graph[v] = append(graph[v], sEdge{u, len(graph[u]) - 1, 0, 0, -cost, -1})
	}

	for i := 0; i < m; i++ {
		u, v, c, w := edges[i].u, edges[i].v, edges[i].c, edges[i].w
		r := c % 2
		bal[u] -= r
		bal[v] += r
		C := c / 2
		W := 2 * w
		if W < 0 {
			bal[u] -= 2 * C
			bal[v] += 2 * C
			addEdge(v, u, C, -W, i)
		} else {
			addEdge(u, v, C, W, i)
		}
	}

	if bal[1]%2 != 0 {
		bal[1] += 1
		bal[n] -= 1
	}
	addEdge(n, 1, solveINF, 0, -1)

	sumReq := 0
	for i := 1; i <= n; i++ {
		if bal[i]%2 != 0 {
			return "Impossible"
		}
		req := bal[i] / 2
		if req > 0 {
			addEdge(S, i, req, 0, -1)
			sumReq += req
		} else if req < 0 {
			addEdge(i, T, -req, 0, -1)
		}
	}

	flowTotal := 0
	pot := make([]int, nNodes)

	for {
		dist := make([]int, nNodes)
		for i := range dist {
			dist[i] = solveINF
		}
		parentEdge := make([]int, nNodes)
		parentVertex := make([]int, nNodes)
		for i := range parentEdge {
			parentEdge[i] = -1
			parentVertex[i] = -1
		}

		dist[S] = 0
		pq := sPQ{}
		heap.Push(&pq, &sItem{vertex: S, priority: 0})

		for pq.Len() > 0 {
			curr := heap.Pop(&pq).(*sItem)
			u := curr.vertex
			if curr.priority > dist[u] {
				continue
			}

			for i, e := range graph[u] {
				if e.cap > e.flow {
					newDist := dist[u] + e.cost + pot[u] - pot[e.to]
					if newDist < dist[e.to] {
						dist[e.to] = newDist
						parentVertex[e.to] = u
						parentEdge[e.to] = i
						heap.Push(&pq, &sItem{vertex: e.to, priority: newDist})
					}
				}
			}
		}

		if dist[T] == solveINF {
			break
		}

		for i := 0; i < nNodes; i++ {
			if dist[i] != solveINF {
				pot[i] += dist[i]
			}
		}

		push := solveINF
		curr := T
		for curr != S {
			p := parentVertex[curr]
			idx := parentEdge[curr]
			rem := graph[p][idx].cap - graph[p][idx].flow
			if rem < push {
				push = rem
			}
			curr = p
		}

		flowTotal += push
		curr = T
		for curr != S {
			p := parentVertex[curr]
			idx := parentEdge[curr]
			rev := graph[p][idx].rev
			graph[p][idx].flow += push
			graph[curr][rev].flow -= push
			curr = p
		}
	}

	if flowTotal < sumReq {
		return "Impossible"
	}

	ans := make([]int, m)
	for u := 0; u < nNodes; u++ {
		for _, e := range graph[u] {
			if e.origID != -1 {
				id := e.origID
				origC := edges[id].c
				origW := edges[id].w
				r := origC % 2
				C := origC / 2

				if origW < 0 {
					ans[id] = r + 2*(C-e.flow)
				} else {
					ans[id] = r + 2*e.flow
				}
			}
		}
	}

	var buf strings.Builder
	buf.WriteString("Possible\n")
	for i := 0; i < m; i++ {
		if i > 0 {
			buf.WriteByte(' ')
		}
		fmt.Fprintf(&buf, "%d", ans[i])
	}
	return buf.String()
}

func genTestF(rng *rand.Rand) string {
	n := rng.Intn(4) + 2
	m := rng.Intn(6) + 1
	var buf strings.Builder
	fmt.Fprintf(&buf, "%d %d\n", n, m)
	for i := 0; i < m; i++ {
		x := rng.Intn(n-1) + 1
		y := rng.Intn(n-x) + x + 1
		if y > n {
			y = n
		}
		if y == 1 {
			y = n
		}
		c := rng.Intn(5) + 1
		w := rng.Intn(11) - 5
		fmt.Fprintf(&buf, "%d %d %d %d\n", x, y, c, w)
	}
	return buf.String()
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	for i := 1; i <= 100; i++ {
		in := genTestF(rng)
		expect := solveF(in)
		got, err := run(exe, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i, in, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
