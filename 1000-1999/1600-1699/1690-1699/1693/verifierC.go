package main

import (
	"bytes"
	"container/heap"
	"context"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type item struct {
	node int
	dist int
}

type priorityQueue []item

func (pq priorityQueue) Len() int            { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool  { return pq[i].dist < pq[j].dist }
func (pq priorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *priorityQueue) Push(x interface{}) { *pq = append(*pq, x.(item)) }
func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[:n-1]
	return x
}

func solveC(n int, edges [][2]int) int {
	rev := make([][]int, n+1)
	outDeg := make([]int, n+1)
	for _, e := range edges {
		v, u := e[0], e[1]
		rev[u] = append(rev[u], v)
		outDeg[v]++
	}
	const inf = math.MaxInt32
	dist := make([]int, n+1)
	for i := range dist {
		dist[i] = inf
	}
	processed := make([]int, n+1)
	pq := &priorityQueue{}
	heap.Init(pq)
	dist[n] = 0
	heap.Push(pq, item{n, 0})
	for pq.Len() > 0 {
		cur := heap.Pop(pq).(item)
		u := cur.node
		if cur.dist != dist[u] {
			continue
		}
		for _, v := range rev[u] {
			processed[v]++
			cand := dist[u] + 1 + outDeg[v] - processed[v]
			if cand < dist[v] {
				dist[v] = cand
				heap.Push(pq, item{v, cand})
			}
		}
	}
	return dist[1]
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(44)
	const t = 100
	var sb strings.Builder
	var exp strings.Builder
	for i := 0; i < t; i++ {
		n := rand.Intn(5) + 2
		// create path edges
		edges := make([][2]int, 0)
		for j := 1; j < n; j++ {
			edges = append(edges, [2]int{j, j + 1})
		}
		// add random extra edges
		extra := rand.Intn(n)
		for k := 0; k < extra; k++ {
			v := rand.Intn(n-1) + 1
			u := rand.Intn(n-v) + v + 1
			edges = append(edges, [2]int{v, u})
		}
		sb.WriteString(fmt.Sprintf("%d %d\n", n, len(edges)))
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		exp.WriteString(fmt.Sprintf("%d\n", solveC(n, edges)))
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error running binary: %v\noutput:\n%s", err, out.String())
		os.Exit(1)
	}
	got := strings.TrimSpace(out.String())
	want := strings.TrimSpace(exp.String())
	if got != want {
		fmt.Fprintf(os.Stderr, "wrong answer\nexpected:\n%s\ngot:\n%s\n", want, got)
		os.Exit(1)
	}
	fmt.Println("all tests passed")
}
