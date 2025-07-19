package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

type Person struct {
	id  int
	a   int
	b   int
	ans int
}

type Item struct {
	b   int
	idx int // index in array p
}

type PQ []Item

func (pq PQ) Len() int { return len(pq) }
func (pq PQ) Less(i, j int) bool {
	if pq[i].b == pq[j].b {
		return pq[i].idx < pq[j].idx
	}
	return pq[i].b < pq[j].b
}
func (pq PQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PQ) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PQ) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[:n-1]
	return x
}

type Pair struct {
	b  int
	id int
}

func pairMax(x, y Pair) Pair {
	if x.b > y.b {
		return x
	} else if x.b < y.b {
		return y
	}
	if x.id > y.id {
		return x
	}
	return y
}

var (
	n  int
	p  []Person
	ft []Pair
)

func update(u, b, id int) {
	for u > 0 {
		if pairMax(ft[u], Pair{b, id}) != ft[u] {
			ft[u] = pairMax(ft[u], Pair{b, id})
		}
		u -= u & -u
	}
}

func get(u int) Pair {
	res := Pair{0, 0}
	for u <= n {
		res = pairMax(res, ft[u])
		u += u & -u
	}
	return res
}

func output(out *bufio.Writer) {
	for i := 0; i < n; i++ {
		fmt.Fprintf(out, "%d", p[i].ans)
		if i+1 == n {
			fmt.Fprint(out, "\n")
		} else {
			fmt.Fprint(out, " ")
		}
	}
}

func solve(in *bufio.Reader, out *bufio.Writer) {
	fmt.Fscan(in, &n)
	p = make([]Person, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &p[i].a, &p[i].b)
		p[i].id = i + 1
	}
	sort.Slice(p, func(i, j int) bool { return p[i].a < p[j].a })
	pq := &PQ{}
	heap.Init(pq)
	j := 0
	for i := 1; i <= n; i++ {
		for j < n && p[j].a <= i {
			heap.Push(pq, Item{b: p[j].b, idx: j})
			j++
		}
		item := heap.Pop(pq).(Item)
		p[item.idx].ans = i
	}

	sort.Slice(p, func(i, j int) bool { return p[i].ans < p[j].ans })

	ft = make([]Pair, n+2)
	resU, resV := 0, 0
	for i := 0; i < n; i++ {
		mx := get(p[i].a)
		if mx.b >= p[i].ans {
			resU = p[i].id
			resV = mx.id
			break
		}
		update(p[i].ans, p[i].b, p[i].id)
	}

	sort.Slice(p, func(i, j int) bool { return p[i].id < p[j].id })

	if resU != 0 {
		fmt.Fprintln(out, "NO")
		output(out)
		for i := range p {
			if p[i].id == resU {
				for j := range p {
					if p[j].id == resV {
						p[i].ans, p[j].ans = p[j].ans, p[i].ans
						break
					}
				}
				break
			}
		}
	} else {
		fmt.Fprintln(out, "YES")
	}
	output(out)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	solve(in, out)
}
