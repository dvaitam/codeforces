package main

import (
	"bufio"
	"fmt"
	"os"
)

type DSU struct{ parent [26]int }

func newDSU() *DSU {
	d := &DSU{}
	for i := 0; i < 26; i++ {
		d.parent[i] = i
	}
	return d
}

func (d *DSU) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) union(x, y int) {
	fx := d.find(x)
	fy := d.find(y)
	if fx != fy {
		d.parent[fx] = fy
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
		var t string
		fmt.Fscan(in, &t)

		dsu := newDSU()
		pre := make([]int, 26)
		used := make([]bool, 26)
		for i := range pre {
			pre[i] = -1
		}

		res := make([]byte, n)
		for i := 0; i < n; i++ {
			x := int(t[i] - 'a')
			if pre[x] != -1 {
				res[i] = byte(pre[x] + 'a')
				continue
			}
			assigned := false
			for y := 0; y < 26; y++ {
				if y == x || used[y] {
					continue
				}
				if dsu.find(x) != dsu.find(y) {
					pre[x] = y
					used[y] = true
					dsu.union(x, y)
					assigned = true
					break
				}
			}
			if !assigned {
				for y := 0; y < 26; y++ {
					if y == x || used[y] {
						continue
					}
					pre[x] = y
					used[y] = true
					dsu.union(x, y)
					break
				}
			}
			res[i] = byte(pre[x] + 'a')
		}
		fmt.Fprintln(out, string(res))
	}
}
