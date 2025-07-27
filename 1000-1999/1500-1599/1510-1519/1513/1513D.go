package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type dsu struct {
	parent []int
	size   []int
}

func newDSU(n int) *dsu {
	p := make([]int, n)
	s := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i
		s[i] = 1
	}
	return &dsu{parent: p, size: s}
}

func (d *dsu) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *dsu) union(a, b int) bool {
	ra := d.find(a)
	rb := d.find(b)
	if ra == rb {
		return false
	}
	if d.size[ra] < d.size[rb] {
		ra, rb = rb, ra
	}
	d.parent[rb] = ra
	d.size[ra] += d.size[rb]
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		solve(reader, writer)
	}
}

func solve(reader *bufio.Reader, writer *bufio.Writer) {
	var n int
	var p int
	fmt.Fscan(reader, &n, &p)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	type pair struct{ val, idx int }
	ord := make([]pair, n)
	for i := 0; i < n; i++ {
		ord[i] = pair{a[i], i}
	}
	sort.Slice(ord, func(i, j int) bool { return ord[i].val < ord[j].val })
	d := newDSU(n)
	used := 0
	var res int64
	for _, pr := range ord {
		if pr.val >= p {
			break
		}
		idx := pr.idx
		val := pr.val
		// extend left
		j := idx - 1
		for j >= 0 {
			if a[j]%val != 0 {
				break
			}
			if !d.union(j, idx) {
				break
			}
			res += int64(val)
			used++
			j--
			if used == n-1 {
				break
			}
		}
		if used == n-1 {
			break
		}
		// extend right
		j = idx + 1
		for j < n {
			if a[j]%val != 0 {
				break
			}
			if !d.union(j, idx) {
				break
			}
			res += int64(val)
			used++
			j++
			if used == n-1 {
				break
			}
		}
		if used == n-1 {
			break
		}
	}
	if used < n-1 {
		res += int64(n-1-used) * int64(p)
	}
	fmt.Fprintln(writer, res)
}
