package main

import (
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

func run(bin, input string) (string, error) {
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type edge struct{ to, w int }

type item struct{ v, d int }

// priority queue for Dijkstra
type pq []item

func (p pq) Len() int            { return len(p) }
func (p pq) Less(i, j int) bool  { return p[i].d < p[j].d }
func (p pq) Swap(i, j int)       { p[i], p[j] = p[j], p[i] }
func (p *pq) Push(x interface{}) { *p = append(*p, x.(item)) }
func (p *pq) Pop() interface{} {
	old := *p
	v := old[len(old)-1]
	*p = old[:len(old)-1]
	return v
}

func solveC(a []int, r1, c1, r2, c2 int) int {
	n := len(a)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		b[i] = a[i] + 1
	}
	cols := []int{c1, c2}
	for _, v := range b {
		cols = append(cols, v)
	}
	sort.Ints(cols)
	uniq := cols[:0]
	for i, v := range cols {
		if i == 0 || v != cols[i-1] {
			uniq = append(uniq, v)
		}
	}
	cols = uniq
	m := len(cols)
	idx := make(map[int]int, m)
	for i, v := range cols {
		idx[v] = i
	}
	valid := make([][]int, n)
	for i := 0; i < n; i++ {
		for j, v := range cols {
			if v <= b[i] {
				valid[i] = append(valid[i], j)
			} else {
				break
			}
		}
	}
	N := n * m
	adj := make([][]edge, N)
	id := func(r, c int) int { return r*m + c }
	for i := 0; i < n; i++ {
		vs := valid[i]
		for k := 1; k < len(vs); k++ {
			u := vs[k-1]
			v := vs[k]
			w := cols[v] - cols[u]
			ui := id(i, u)
			vi := id(i, v)
			adj[ui] = append(adj[ui], edge{vi, w})
			adj[vi] = append(adj[vi], edge{ui, w})
		}
	}
	for i := 0; i < n; i++ {
		for _, u := range valid[i] {
			cu := cols[u]
			ui := id(i, u)
			if i > 0 {
				tc := cu
				if tc > b[i-1] {
					tc = b[i-1]
				}
				if v, ok := idx[tc]; ok {
					adj[ui] = append(adj[ui], edge{id(i-1, v), 1})
				}
			}
			if i+1 < n {
				tc := cu
				if tc > b[i+1] {
					tc = b[i+1]
				}
				if v, ok := idx[tc]; ok {
					adj[ui] = append(adj[ui], edge{id(i+1, v), 1})
				}
			}
		}
	}
	const INF = int(1e9)
	dist := make([]int, N)
	for i := range dist {
		dist[i] = INF
	}
	si := id(r1, idx[c1])
	ti := id(r2, idx[c2])
	dist[si] = 0
	pq := &pq{{si, 0}}
	heap.Init(pq)
	for pq.Len() > 0 {
		it := heap.Pop(pq).(item)
		v, d := it.v, it.d
		if d != dist[v] {
			continue
		}
		if v == ti {
			break
		}
		for _, e := range adj[v] {
			nd := d + e.w
			if nd < dist[e.to] {
				dist[e.to] = nd
				heap.Push(pq, item{e.to, nd})
			}
		}
	}
	return dist[ti]
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for tc := 1; tc <= 100; tc++ {
		n := rng.Intn(10) + 1
		a := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = rng.Intn(10)
		}
		r1 := rng.Intn(n) + 1
		r2 := rng.Intn(n) + 1
		c1 := rng.Intn(a[r1-1]+1) + 1
		c2 := rng.Intn(a[r2-1]+1) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i, v := range a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteString(fmt.Sprintf("\n%d %d %d %d\n", r1, c1, r2, c2))
		input := sb.String()
		expect := fmt.Sprintf("%d", solveC(append([]int(nil), a...), r1-1, c1, r2-1, c2))
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", tc, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", tc, expect, strings.TrimSpace(out), input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
