package main

import (
	"bufio"
	"fmt"
	"os"
)

type FastScanner struct {
	r *bufio.Reader
}

func NewFastScanner() *FastScanner {
	return &FastScanner{r: bufio.NewReader(os.Stdin)}
}

func (fs *FastScanner) NextInt() int {
	sign := 1
	val := 0
	c, _ := fs.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		c, _ = fs.r.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, _ = fs.r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int(c-'0')
		c, _ = fs.r.ReadByte()
	}
	return sign * val
}

func (fs *FastScanner) NextInt64() int64 {
	sign := int64(1)
	var val int64
	c, _ := fs.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		c, _ = fs.r.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, _ = fs.r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int64(c-'0')
		c, _ = fs.r.ReadByte()
	}
	return sign * val
}

type Edge struct {
	to int
	w  int64
}

func main() {
	fs := NewFastScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t := fs.NextInt()
	for ; t > 0; t-- {
		n := fs.NextInt()
		adj := make([][]Edge, n)
		for i := 0; i < n-1; i++ {
			u := fs.NextInt() - 1
			v := fs.NextInt() - 1
			w := fs.NextInt64()
			adj[u] = append(adj[u], Edge{to: v, w: w})
			adj[v] = append(adj[v], Edge{to: u, w: w})
		}

		parent := make([]int, n)
		for i := range parent {
			parent[i] = -1
		}
		parent[0] = -2
		parentEdge := make([]int64, n)
		order := make([]int, 0, n)
		stack := []int{0}
		for len(stack) > 0 {
			v := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			order = append(order, v)
			for _, e := range adj[v] {
				if parent[e.to] == -1 {
					parent[e.to] = v
					parentEdge[e.to] = e.w
					stack = append(stack, e.to)
				}
			}
		}
		parent[0] = -1
		root := order[0]

		dpFree := make([]int64, n)
		dpAttach := make([]int64, n)
		dpBest := make([]int64, n)

		q := fs.NextInt()
		for ; q > 0; q-- {
			l := fs.NextInt64()
			for idx := len(order) - 1; idx >= 0; idx-- {
				v := order[idx]
				total := int64(0)
				attach := int64(0)
				for _, e := range adj[v] {
					if e.to == parent[v] {
						continue
					}
					child := e.to
					total += dpBest[child]
					include := dpAttach[child]
					cut := dpBest[child] - e.w
					if cut > include {
						attach += cut
					} else {
						attach += include
					}
				}
				dpFree[v] = total
				parentW := int64(0)
				if parent[v] != -1 {
					parentW = parentEdge[v]
				}
				take := l - parentW + attach
				if take < 0 {
					take = 0
				}
				if take < dpFree[v] {
					dpBest[v] = dpFree[v]
				} else {
					dpBest[v] = take
				}
				dpAttach[v] = attach
			}
			fmt.Fprintln(out, dpBest[root])
		}
	}
}
