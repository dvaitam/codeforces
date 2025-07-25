package main

import (
	"bufio"
	"fmt"
	"os"
)

type DSU struct {
	parent []int
	rank   []int
	parity []int
}

func NewDSU(n int) *DSU {
	d := &DSU{
		parent: make([]int, n),
		rank:   make([]int, n),
		parity: make([]int, n),
	}
	for i := range d.parent {
		d.parent[i] = i
	}
	return d
}

func (d *DSU) find(x int) (int, int) {
	if d.parent[x] == x {
		return x, d.parity[x]
	}
	r, p := d.find(d.parent[x])
	d.parent[x] = r
	d.parity[x] ^= p
	return d.parent[x], d.parity[x]
}

func (d *DSU) union(x, y, rel int) bool {
	rx, px := d.find(x)
	ry, py := d.find(y)
	if rx == ry {
		return (px ^ py) == rel
	}
	if d.rank[rx] < d.rank[ry] {
		rx, ry = ry, rx
		px, py = py, px
	}
	d.parent[ry] = rx
	d.parity[ry] = px ^ py ^ rel
	if d.rank[rx] == d.rank[ry] {
		d.rank[rx]++
	}
	return true
}

func reverse(s string) string {
	b := []byte(s)
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m, k int
		if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
			return
		}
		s := make([]string, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &s[i])
		}
		r := make([]string, n)
		for i := 0; i < n; i++ {
			r[i] = reverse(s[i])
		}
		dsu := NewDSU(n)
		possible := true
		for i := 0; i < n && possible; i++ {
			for j := i + 1; j < n; j++ {
				same, diff := 0, 0
				si := s[i]
				sj := s[j]
				rj := r[j]
				for t2 := 0; t2 < m; t2++ {
					if si[t2] == sj[t2] {
						same++
					}
					if si[t2] == rj[t2] {
						diff++
					}
				}
				if same < k && diff < k {
					possible = false
					break
				} else if same >= k && diff < k {
					if !dsu.union(i, j, 0) {
						possible = false
						break
					}
				} else if same < k && diff >= k {
					if !dsu.union(i, j, 1) {
						possible = false
						break
					}
				}
			}
		}
		if !possible {
			fmt.Println(-1)
			continue
		}
		par := make([]int, n)
		comps := make(map[int][]int)
		for i := 0; i < n; i++ {
			root, p := dsu.find(i)
			par[i] = p
			comps[root] = append(comps[root], i)
		}
		var toReverse []int
		for _, list := range comps {
			count1 := 0
			for _, idx := range list {
				if par[idx] == 1 {
					count1++
				}
			}
			if count1 <= len(list)-count1 {
				for _, idx := range list {
					if par[idx] == 1 {
						toReverse = append(toReverse, idx+1)
					}
				}
			} else {
				for _, idx := range list {
					if par[idx] == 0 {
						toReverse = append(toReverse, idx+1)
					}
				}
			}
		}
		fmt.Println(len(toReverse))
		if len(toReverse) > 0 {
			for i, v := range toReverse {
				if i > 0 {
					fmt.Print(" ")
				}
				fmt.Print(v)
			}
		}
		fmt.Println()
	}
}
