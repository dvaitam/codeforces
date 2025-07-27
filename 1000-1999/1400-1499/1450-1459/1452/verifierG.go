package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type state struct {
	th  int
	rem int
	u   int
}

type pqState []state

func (pq pqState) Len() int            { return len(pq) }
func (pq pqState) Less(i, j int) bool  { return pq[i].th > pq[j].th }
func (pq pqState) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *pqState) Push(x interface{}) { *pq = append(*pq, x.(state)) }
func (pq *pqState) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[:n-1]
	return x
}

func solveG(n int, edges [][]int, bob []int) []int {
	adj := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	const INF = int(1e9)
	dBob := make([]int, n+1)
	for i := 1; i <= n; i++ {
		dBob[i] = INF
	}
	q := make([]int, 0, n)
	for _, v := range bob {
		dBob[v] = 0
		q = append(q, v)
	}
	for qi := 0; qi < len(q); qi++ {
		u := q[qi]
		for _, v := range adj[u] {
			if dBob[v] > dBob[u]+1 {
				dBob[v] = dBob[u] + 1
				q = append(q, v)
			}
		}
	}
	bestTh := make([]int, n+1)
	for i := range bestTh {
		bestTh[i] = -1
	}
	pqh := &pqState{}
	heap.Init(pqh)
	for u := 1; u <= n; u++ {
		th := dBob[u]
		heap.Push(pqh, state{th: th, rem: th, u: u})
	}
	for pqh.Len() > 0 {
		st := heap.Pop(pqh).(state)
		u, th, rem := st.u, st.th, st.rem
		if th <= bestTh[u] {
			continue
		}
		bestTh[u] = th
		if rem == 0 {
			continue
		}
		for _, v := range adj[u] {
			if th > bestTh[v] {
				heap.Push(pqh, state{th: th, rem: rem - 1, u: v})
			}
		}
	}
	ans := make([]int, n)
	for i := 1; i <= n; i++ {
		ans[i-1] = bestTh[i]
	}
	return ans
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("run error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		return
	}
	binary := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	const tnum = 100
	for ti := 0; ti < tnum; ti++ {
		n := rand.Intn(8) + 2
		edges := make([][]int, n-1)
		for i := 2; i <= n; i++ {
			p := rand.Intn(i-1) + 1
			edges[i-2] = []int{p, i}
		}
		k := rand.Intn(n) + 1
		perm := rand.Perm(n)
		bob := make([]int, k)
		for i := 0; i < k; i++ {
			bob[i] = perm[i] + 1
		}
		exp := solveG(n, edges, bob)
		var in strings.Builder
		fmt.Fprintf(&in, "%d\n", n)
		for _, e := range edges {
			fmt.Fprintf(&in, "%d %d\n", e[0], e[1])
		}
		fmt.Fprintf(&in, "%d\n", k)
		for i, v := range bob {
			if i > 0 {
				in.WriteByte(' ')
			}
			fmt.Fprintf(&in, "%d", v)
		}
		in.WriteByte('\n')
		out, err := runBinary(binary, in.String())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fields := strings.Fields(strings.TrimSpace(out))
		if len(fields) != n {
			fmt.Printf("test %d failed: expected %d values got %d\noutput:\n%s\n", ti+1, n, len(fields), out)
			os.Exit(1)
		}
		for i, f := range fields {
			v, err := strconv.Atoi(f)
			if err != nil || v != exp[i] {
				fmt.Printf("test %d failed at vertex %d: expected %d got %s\n", ti+1, i+1, exp[i], f)
				os.Exit(1)
			}
		}
	}
	fmt.Println("all tests passed")
}
