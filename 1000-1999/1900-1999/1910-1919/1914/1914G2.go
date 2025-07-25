package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const MOD int = 998244353

type DSU struct {
	parent []int
}

func NewDSU(n int) *DSU {
	d := &DSU{parent: make([]int, n)}
	for i := 0; i < n; i++ {
		d.parent[i] = i
	}
	return d
}

func (d *DSU) Find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.Find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) Union(a, b int) {
	ra, rb := d.Find(a), d.Find(b)
	if ra != rb {
		d.parent[rb] = ra
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		m := 2 * n
		colors := make([]int, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &colors[i])
		}

		first := make([]int, n+1)
		for i := range first {
			first[i] = -1
		}
		intervals := make([][3]int, n) // l,r,index
		for i, c := range colors {
			if first[c] == -1 {
				first[c] = i
			} else {
				intervals[c-1] = [3]int{first[c], i, c - 1}
			}
		}

		// DSU1: connected components by overlap
		idx := make([]int, n)
		for i := 0; i < n; i++ {
			idx[i] = i
		}
		sort.Slice(idx, func(i, j int) bool {
			return intervals[idx[i]][0] < intervals[idx[j]][0]
		})
		dsu1 := NewDSU(n)
		curR := intervals[idx[0]][1]
		last := idx[0]
		for k := 1; k < n; k++ {
			i := idx[k]
			l, r := intervals[i][0], intervals[i][1]
			if l > curR {
				curR = r
				last = i
			} else {
				dsu1.Union(last, i)
				if r > curR {
					curR = r
					last = i
				}
			}
		}

		// DSU2: strongly connected components via crossing
		dsu2 := NewDSU(n)
		sort.Slice(intervals, func(i, j int) bool { return intervals[i][0] < intervals[j][0] })
		stack := make([][3]int, 0)
		for _, it := range intervals {
			l, r, idx := it[0], it[1], it[2]
			for len(stack) > 0 && stack[len(stack)-1][1] < l {
				stack = stack[:len(stack)-1]
			}
			for len(stack) > 0 && stack[len(stack)-1][1] < r {
				dsu2.Union(idx, stack[len(stack)-1][2])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, [3]int{l, r, idx})
		}

		// detect incoming edges between SCCs using containment
		rootIn := make([]bool, n)
		byLR := make([][3]int, n)
		copy(byLR, intervals)
		sort.Slice(byLR, func(i, j int) bool {
			if byLR[i][0] == byLR[j][0] {
				return byLR[i][1] > byLR[j][1]
			}
			return byLR[i][0] < byLR[j][0]
		})
		maxR := -1
		maxIdx := -1
		for _, it := range byLR {
			l, r, idx := it[0], it[1], it[2]
			if r <= maxR {
				ri := dsu2.Find(idx)
				rj := dsu2.Find(maxIdx)
				if ri != rj {
					rootIn[ri] = true
				}
			} else {
				maxR = r
				maxIdx = idx
			}
			_ = l
		}
		byRL := make([][3]int, n)
		copy(byRL, intervals)
		sort.Slice(byRL, func(i, j int) bool {
			if byRL[i][1] == byRL[j][1] {
				return byRL[i][0] < byRL[j][0]
			}
			return byRL[i][1] > byRL[j][1]
		})
		minL := 1 << 30
		minIdx := -1
		for _, it := range byRL {
			l, r, idx := it[0], it[1], it[2]
			if l >= minL {
				ri := dsu2.Find(idx)
				rj := dsu2.Find(minIdx)
				if ri != rj {
					rootIn[ri] = true
				}
			} else {
				minL = l
				minIdx = idx
			}
			_ = r
		}

		// group by connected components
		compMap := map[int][]int{}
		for i := 0; i < n; i++ {
			root := dsu1.Find(i)
			compMap[root] = append(compMap[root], i)
		}
		minSize := len(compMap)
		ways := 1
		for _, members := range compMap {
			sccMap := map[int][]int{}
			for _, i := range members {
				root := dsu2.Find(i)
				sccMap[root] = append(sccMap[root], i)
			}
			var source int
			for r := range sccMap {
				if !rootIn[r] {
					source = r
					break
				}
			}
			ways = ways * (2 * len(sccMap[source])) % MOD
		}
		fmt.Fprintln(out, minSize, ways)
	}
}
