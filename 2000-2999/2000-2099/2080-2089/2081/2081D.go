package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const maxV = 500000 + 5
const inf = int(1e9)

type Edge struct {
	w   int
	u   int
	v   int
}

type DSU struct {
	parent []int
	sz     []int
}

func (d *DSU) init(maxVal int, touched []int) {
	for _, v := range touched {
		d.parent[v] = v
		d.sz[v] = 1
	}
}
func (d *DSU) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}
func (d *DSU) union(x, y int) bool {
	x = d.find(x)
	y = d.find(y)
	if x == y {
		return false
	}
	if d.sz[x] < d.sz[y] {
		x, y = y, x
	}
	d.parent[y] = x
	d.sz[x] += d.sz[y]
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	freq := make([]int, maxV)
	parent := make([]int, maxV)
	sz := make([]int, maxV)
	dsu := DSU{parent: parent, sz: sz}
	touched := make([]int, 0, 1<<10)

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		if n == 0 {
			fmt.Fprintln(out, 0)
			continue
		}

		vals := make([]int, n)
		minVal, maxVal := inf, 0
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &vals[i])
			v := vals[i]
			freq[v]++
			if freq[v] == 1 {
				touched = append(touched, v)
			}
			if v < minVal {
				minVal = v
			}
			if v > maxVal {
				maxVal = v
			}
		}

		// Initialize DSU for touched values
		dsu.init(maxVal, touched)

		// Next present value at or after i
		nxt := make([]int, maxVal+2)
		fillVal := maxVal + 1
		for i := range nxt {
			nxt[i] = fillVal
		}
		cur := fillVal
		for i := maxVal; i >= 1; i-- {
			if freq[i] > 0 {
				cur = i
			}
			nxt[i] = cur
		}

		// Union zero-weight edges (divisibility)
		for _, d := range touched {
			if d == 0 || freq[d] == 0 {
				continue
			}
			for m := d * 2; m <= maxVal; m += d {
				if freq[m] > 0 {
					dsu.union(d, m)
				}
			}
		}

		// Collect candidate edges
		edges := make([]Edge, 0)
		for _, d := range touched {
			if d > maxVal || freq[d] == 0 {
				continue
			}
			for m := d; m <= maxVal; m += d {
				v := nxt[m]
				if v > maxVal || v > m+d-1 {
					continue
				}
				if v == d {
					// Skip self; look for the next present value in this interval.
					v = nxt[v+1]
					if v > maxVal || v > m+d-1 {
						continue
					}
				}
				rem := v - m
				uComp := dsu.find(d)
				vComp := dsu.find(v)
				if uComp != vComp {
					edges = append(edges, Edge{w: rem, u: uComp, v: vComp})
				}
			}
		}

		sort.Slice(edges, func(i, j int) bool {
			return edges[i].w < edges[j].w
		})

		compCount := 0
		for _, v := range touched {
			if dsu.find(v) == v {
				compCount++
			}
		}

		var ans int64
		for _, e := range edges {
			if dsu.union(e.u, e.v) {
				ans += int64(e.w)
				compCount--
				if compCount == 1 {
					break
				}
			}
		}

		fmt.Fprintln(out, ans)

		// reset freq for touched values
		for _, v := range touched {
			freq[v] = 0
		}
		touched = touched[:0]
	}
}
