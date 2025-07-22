package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type Edge struct {
	to, rk int
}

type Item struct {
	node, w int
}

type MaxHeap []Item

func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool  { return h[i].w > h[j].w }
func (h MaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(Item)) }
func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	it := old[n-1]
	*h = old[:n-1]
	return it
}

func widestPath(n int, adj [][]Edge) int {
	const INF = 1000000001
	width := make([]int, n+1)
	for i := 1; i <= n; i++ {
		width[i] = 0
	}
	width[1] = INF
	h := &MaxHeap{}
	heap.Init(h)
	heap.Push(h, Item{node: 1, w: INF})
	for h.Len() > 0 {
		it := heap.Pop(h).(Item)
		u, w := it.node, it.w
		if w < width[u] {
			continue
		}
		if u == n {
			break
		}
		for _, e := range adj[u] {
			nw := w
			if e.rk < nw {
				nw = e.rk
			}
			if nw > width[e.to] {
				width[e.to] = nw
				heap.Push(h, Item{node: e.to, w: nw})
			}
		}
	}
	return width[n]
}

func solve(input string) string {
	r := bufio.NewReader(strings.NewReader(input))
	var n, m int
	if _, err := fmt.Fscan(r, &n, &m); err != nil {
		return ""
	}
	type E struct{ u, v, l, r int }
	edges := make([]E, m)
	ls := make([]int, 0, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(r, &edges[i].u, &edges[i].v, &edges[i].l, &edges[i].r)
		ls = append(ls, edges[i].l)
	}
	if m == 0 {
		return "Nice work, Dima!\n"
	}
	sort.Ints(ls)
	us := ls[:0]
	for i, v := range ls {
		if i == 0 || v != ls[i-1] {
			us = append(us, v)
		}
	}
	answer := 0
	globalMaxR := 0
	for _, e := range edges {
		if e.r > globalMaxR {
			globalMaxR = e.r
		}
	}
	for _, L := range us {
		if globalMaxR-L+1 <= answer {
			break
		}
		adj := make([][]Edge, n+1)
		for _, e := range edges {
			if e.l <= L {
				adj[e.u] = append(adj[e.u], Edge{to: e.v, rk: e.r})
				adj[e.v] = append(adj[e.v], Edge{to: e.u, rk: e.r})
			}
		}
		W := widestPath(n, adj)
		if W >= L {
			loyalty := W - L + 1
			if loyalty > answer {
				answer = loyalty
			}
		}
	}
	if answer <= 0 {
		return "Nice work, Dima!\n"
	}
	return fmt.Sprintf("%d\n", answer)
}

func genTest() (string, string) {
	n := rand.Intn(4) + 2
	m := rand.Intn(5)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < m; i++ {
		u := rand.Intn(n) + 1
		v := rand.Intn(n) + 1
		l := rand.Intn(10) + 1
		r := l + rand.Intn(10)
		fmt.Fprintf(&sb, "%d %d %d %d\n", u, v, l, r)
	}
	inp := sb.String()
	out := solve(inp)
	return inp, out
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, exp := genTest()
		got, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\nInput:\n%sOutput:\n%s\n", i+1, err, in, got)
			return
		}
		if got != strings.TrimSpace(exp) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s\n", i+1, in, exp, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
