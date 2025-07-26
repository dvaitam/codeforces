package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type dsu struct {
	parent []int
	rank   []byte
}

func newDSU(n int) *dsu {
	d := &dsu{parent: make([]int, n), rank: make([]byte, n)}
	for i := range d.parent {
		d.parent[i] = i
	}
	return d
}

func (d *dsu) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *dsu) union(a, b int) {
	a = d.find(a)
	b = d.find(b)
	if a == b {
		return
	}
	if d.rank[a] < d.rank[b] {
		a, b = b, a
	}
	d.parent[b] = a
	if d.rank[a] == d.rank[b] {
		d.rank[a]++
	}
}

type mine struct {
	x, y int
	t    int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		mines := make([]mine, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &mines[i].x, &mines[i].y, &mines[i].t)
		}

		d := newDSU(n)

		type pair struct{ val, idx int }
		// group by x
		mpX := make(map[int][]pair)
		for i, m := range mines {
			mpX[m.x] = append(mpX[m.x], pair{m.y, i})
		}
		for _, arr := range mpX {
			sort.Slice(arr, func(i, j int) bool { return arr[i].val < arr[j].val })
			for j := 1; j < len(arr); j++ {
				if arr[j].val-arr[j-1].val <= k {
					d.union(arr[j].idx, arr[j-1].idx)
				}
			}
		}
		// group by y
		mpY := make(map[int][]pair)
		for i, m := range mines {
			mpY[m.y] = append(mpY[m.y], pair{m.x, i})
		}
		for _, arr := range mpY {
			sort.Slice(arr, func(i, j int) bool { return arr[i].val < arr[j].val })
			for j := 1; j < len(arr); j++ {
				if arr[j].val-arr[j-1].val <= k {
					d.union(arr[j].idx, arr[j-1].idx)
				}
			}
		}

		compMin := make(map[int]int)
		maxTime := 0
		for i, m := range mines {
			root := d.find(i)
			if val, ok := compMin[root]; !ok || m.t < val {
				compMin[root] = m.t
			}
			if m.t > maxTime {
				maxTime = m.t
			}
		}

		times := make([]int, 0, len(compMin))
		for _, t := range compMin {
			times = append(times, t)
		}
		sort.Ints(times)

		hi := maxTime
		if hi < len(times) {
			hi = len(times)
		}
		l, r := -1, hi
		for r-l > 1 {
			mid := (l + r) / 2
			idx := sort.Search(len(times), func(i int) bool { return times[i] > mid })
			cnt := len(times) - idx
			if cnt <= mid+1 {
				r = mid
			} else {
				l = mid
			}
		}
		fmt.Fprintln(writer, r)
	}
}
