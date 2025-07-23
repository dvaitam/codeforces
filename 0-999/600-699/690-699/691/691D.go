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
	d := &dsu{parent: make([]int, n+1), size: make([]int, n+1)}
	for i := 1; i <= n; i++ {
		d.parent[i] = i
		d.size[i] = 1
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
	if d.size[a] < d.size[b] {
		a, b = b, a
	}
	d.parent[b] = a
	d.size[a] += d.size[b]
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	p := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &p[i])
	}
	d := newDSU(n)
	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(reader, &a, &b)
		d.union(a, b)
	}

	compPos := make(map[int][]int)
	compVal := make(map[int][]int)
	for i := 1; i <= n; i++ {
		r := d.find(i)
		compPos[r] = append(compPos[r], i)
		compVal[r] = append(compVal[r], p[i])
	}

	ans := make([]int, n+1)
	for root := range compPos {
		pos := compPos[root]
		val := compVal[root]
		sort.Ints(pos)
		sort.Slice(val, func(i, j int) bool { return val[i] > val[j] })
		for i := 0; i < len(pos); i++ {
			ans[pos[i]] = val[i]
		}
	}

	for i := 1; i <= n; i++ {
		if i > 1 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, ans[i])
	}
	writer.WriteByte('\n')
}
