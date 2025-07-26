package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Edge struct {
	to int
	w  int64
}

type Seg struct {
	zone  int
	start int64
	end   int64
}

var (
	n         int
	g         [][]Edge
	k         int
	zones     []int
	pass      []int64
	fine      []int64
	T         int64
	q         int
	dist      []int64
	tin       []int
	tout      []int
	timer     int
	zoneNodes [][]int   // per zone list of nodes ordered by tin
	counts    [][]int64 // [n+1][k]
)

func addEdge(u, v int, w int64) {
	g[u] = append(g[u], Edge{v, w})
	g[v] = append(g[v], Edge{u, w})
}

func dfsDist(u, p int, d int64) {
	timer++
	tin[u] = timer
	dist[u] = d
	zoneNodes[zones[u]] = append(zoneNodes[zones[u]], u)
	for _, e := range g[u] {
		if e.to == p {
			continue
		}
		dfsDist(e.to, u, d+e.w)
	}
	tout[u] = timer + 1
}

func countScans(dv, start, end int64) int64 {
	L := dv - end
	if L < T {
		L = T
	}
	R := dv - start
	if R <= L {
		return 0
	}
	return (R-1)/T - (L-1)/T
}

func dfsCounts(u, p int, segs []Seg) {
	dv := dist[u]
	for _, s := range segs {
		c := countScans(dv, s.start, s.end)
		counts[u][s.zone] += c
	}
	for _, e := range g[u] {
		if e.to == p {
			continue
		}
		zc := zones[e.to]
		if len(segs) > 0 && segs[len(segs)-1].zone == zc {
			oldEnd := segs[len(segs)-1].end
			segs[len(segs)-1].end = dist[e.to]
			dfsCounts(e.to, u, segs)
			segs[len(segs)-1].end = oldEnd
		} else {
			seg := Seg{zone: zc, start: segs[len(segs)-1].end, end: dist[e.to]}
			segs = append(segs, seg)
			dfsCounts(e.to, u, segs)
			segs = segs[:len(segs)-1]
		}
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	g = make([][]Edge, n+1)
	for i := 0; i < n-1; i++ {
		var v, u int
		var t int64
		fmt.Fscan(in, &v, &u, &t)
		addEdge(v, u, t)
	}
	fmt.Fscan(in, &k)
	var zoneStr string
	fmt.Fscan(in, &zoneStr)
	zones = make([]int, n+1)
	for i := 1; i <= n; i++ {
		zones[i] = int(zoneStr[i-1] - 'A')
	}
	pass = make([]int64, k)
	fine = make([]int64, k)
	for i := 0; i < k; i++ {
		fmt.Fscan(in, &pass[i])
	}
	for i := 0; i < k; i++ {
		fmt.Fscan(in, &fine[i])
	}
	fmt.Fscan(in, &T)
	fmt.Fscan(in, &q)

	dist = make([]int64, n+1)
	tin = make([]int, n+1)
	tout = make([]int, n+1)
	zoneNodes = make([][]int, k)
	timer = 0
	dfsDist(1, 0, 0)

	counts = make([][]int64, n+1)
	for i := 0; i <= n; i++ {
		counts[i] = make([]int64, k)
	}
	// initialize segments with root zone
	segs := []Seg{{zone: zones[1], start: 0, end: 0}}
	dfsCounts(1, 0, segs)

	for ; q > 0; q-- {
		var tp int
		fmt.Fscan(in, &tp)
		if tp == 1 {
			var ch string
			var c int64
			fmt.Fscan(in, &ch, &c)
			idx := int(ch[0] - 'A')
			pass[idx] = c
		} else if tp == 2 {
			var ch string
			var c int64
			fmt.Fscan(in, &ch, &c)
			idx := int(ch[0] - 'A')
			fine[idx] = c
		} else if tp == 3 {
			var u int
			fmt.Fscan(in, &u)
			z := zones[u]
			nodes := zoneNodes[z]
			l := sort.Search(len(nodes), func(i int) bool { return tin[nodes[i]] >= tin[u] })
			r := sort.Search(len(nodes), func(i int) bool { return tin[nodes[i]] >= tout[u] })
			best := int64(1<<63 - 1)
			for i := l; i < r; i++ {
				v := nodes[i]
				cost := int64(0)
				for j := 0; j < k; j++ {
					if j == z {
						continue
					}
					cnt := counts[v][j]
					p := pass[j]
					f := fine[j]
					if p < f*cnt {
						cost += p
					} else {
						cost += f * cnt
					}
				}
				if cost < best {
					best = cost
				}
			}
			if best == int64(1<<63-1) {
				best = 0
			}
			fmt.Fprintln(out, best)
		}
	}
}
