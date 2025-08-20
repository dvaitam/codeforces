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
    for i := 1; i <= t; i++ {
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
        // Build single-case input
        var sb strings.Builder
        sb.WriteString(fmt.Sprintf("%d %d\n", n, len(edges)))
        for _, e := range edges {
            sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
        }
        expected := fmt.Sprintf("%d", solveC(n, edges))

        // Run binary for this single test
        ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
        cmd := exec.CommandContext(ctx, bin)
        cmd.Stdin = strings.NewReader(sb.String())
        var out bytes.Buffer
        cmd.Stdout = &out
        cmd.Stderr = &out
        err := cmd.Run()
        cancel()
        if err != nil {
            fmt.Fprintf(os.Stderr, "case %d runtime error: %v\ninput:\n%soutput:\n%s", i, err, sb.String(), out.String())
            os.Exit(1)
        }
        got := strings.TrimSpace(out.String())
        if got != expected {
            fmt.Fprintf(os.Stderr, "wrong answer on case %d\nexpected:\n%s\n\ngot:\n%s\ninput:\n%s", i, expected, got, sb.String())
            os.Exit(1)
        }
    }
    fmt.Println("all tests passed")
}
