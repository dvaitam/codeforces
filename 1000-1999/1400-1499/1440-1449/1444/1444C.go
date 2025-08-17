package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type DSU struct {
	parent []int
	size   []int
	parity []int
}

func NewDSU(n int) *DSU {
	parent := make([]int, n)
	size := make([]int, n)
	parity := make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		size[i] = 1
	}
	return &DSU{parent: parent, size: size, parity: parity}
}

func (d *DSU) Find(x int) (int, int) {
	if d.parent[x] == x {
		return x, 0
	}
	r, p := d.Find(d.parent[x])
	d.parent[x] = r
	d.parity[x] ^= p
	return d.parent[x], d.parity[x]
}

func (d *DSU) Union(x, y, w int) bool {
	rx, px := d.Find(x)
	ry, py := d.Find(y)
	if rx == ry {
		return (px ^ py) == w
	}
	if d.size[rx] < d.size[ry] {
		rx, ry = ry, rx
		px, py = py, px
	}
	d.parent[ry] = rx
	d.parity[ry] = px ^ py ^ w
	d.size[rx] += d.size[ry]
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m, k int
	fmt.Fscan(reader, &n, &m, &k)
	group := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &group[i])
		group[i]--
	}

	dsu := NewDSU(n)
	bad := make([]bool, k)

	type Edge struct{ g1, g2, u, v int }
	cross := make([]Edge, 0, m)

	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		u--
		v--
		ga, gb := group[u], group[v]
		if ga == gb {
			if !bad[ga] && !dsu.Union(u, v, 1) {
				bad[ga] = true
			}
		} else {
			if ga > gb {
				ga, gb = gb, ga
				u, v = v, u
			}
			cross = append(cross, Edge{ga, gb, u, v})
		}
	}

	goodCnt := 0
	for i := 0; i < k; i++ {
		if !bad[i] {
			goodCnt++
		}
	}
	ans := int64(goodCnt) * int64(goodCnt-1) / 2

	sort.Slice(cross, func(i, j int) bool {
		if cross[i].g1 != cross[j].g1 {
			return cross[i].g1 < cross[j].g1
		}
		return cross[i].g2 < cross[j].g2
	})

	for i := 0; i < len(cross); {
		j := i
		g1, g2 := cross[i].g1, cross[i].g2
		for j < len(cross) && cross[j].g1 == g1 && cross[j].g2 == g2 {
			j++
		}
		if !bad[g1] && !bad[g2] {
			nodes := make([]int, 0, 2*(j-i))
			triples := make([]struct{ u, v, w int }, 0, j-i)
			for t := i; t < j; t++ {
				u, v := cross[t].u, cross[t].v
				ru, xu := dsu.Find(u)
				rv, xv := dsu.Find(v)
				nodes = append(nodes, ru, rv)
				triples = append(triples, struct{ u, v, w int }{ru, rv, 1 ^ xu ^ xv})
			}
			sort.Ints(nodes)
			// dedup
			uniq := nodes[:0]
			for _, x := range nodes {
				if len(uniq) == 0 || uniq[len(uniq)-1] != x {
					uniq = append(uniq, x)
				}
			}
			idx := make(map[int]int, len(uniq))
			for idxNum, v := range uniq {
				idx[v] = idxNum
			}
			local := NewDSU(len(uniq))
			ok := true
			for _, tr := range triples {
				iu := idx[tr.u]
				iv := idx[tr.v]
				if !local.Union(iu, iv, tr.w) {
					ok = false
					break
				}
			}
			if !ok {
				ans--
			}
		}
		i = j
	}

	fmt.Fprintln(writer, ans)
}
