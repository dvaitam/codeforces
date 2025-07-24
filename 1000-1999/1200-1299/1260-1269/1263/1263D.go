package main

import (
	"bufio"
	"fmt"
	"os"
)

type DSU struct {
	parent []int
}

func NewDSU(n int) *DSU {
	d := &DSU{parent: make([]int, n)}
	for i := range d.parent {
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
	pa := d.Find(a)
	pb := d.Find(b)
	if pa != pb {
		d.parent[pa] = pb
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	dsu := NewDSU(26)
	used := make([]bool, 26)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(reader, &s)
		first := -1
		for _, ch := range s {
			idx := int(ch - 'a')
			used[idx] = true
			if first == -1 {
				first = idx
			} else {
				dsu.Union(first, idx)
			}
		}
	}

	comps := make(map[int]bool)
	for i := 0; i < 26; i++ {
		if used[i] {
			root := dsu.Find(i)
			comps[root] = true
		}
	}
	fmt.Fprintln(writer, len(comps))
}
