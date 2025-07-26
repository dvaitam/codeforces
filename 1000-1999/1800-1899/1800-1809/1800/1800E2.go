package main

import (
	"bufio"
	"fmt"
	"os"
)

type DSU struct {
	parent []int
	size   []int
}

func NewDSU(n int) *DSU {
	p := make([]int, n)
	sz := make([]int, n)
	for i := range p {
		p[i] = i
		sz[i] = 1
	}
	return &DSU{parent: p, size: sz}
}

func (d *DSU) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) union(a, b int) {
	ra := d.find(a)
	rb := d.find(b)
	if ra == rb {
		return
	}
	if d.size[ra] < d.size[rb] {
		ra, rb = rb, ra
	}
	d.parent[rb] = ra
	d.size[ra] += d.size[rb]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		var s, t string
		fmt.Fscan(in, &s)
		fmt.Fscan(in, &t)

		dsu := NewDSU(n)
		for i := 0; i < n; i++ {
			if i+k < n {
				dsu.union(i, i+k)
			}
			if i+k+1 < n {
				dsu.union(i, i+k+1)
			}
		}

		counts := make(map[int][]int)
		ok := true
		for i := 0; i < n; i++ {
			r := dsu.find(i)
			arr, exists := counts[r]
			if !exists {
				arr = make([]int, 26)
				counts[r] = arr
			}
			arr[s[i]-'a']++
			arr[t[i]-'a']--
		}

		for _, arr := range counts {
			for j := 0; j < 26; j++ {
				if arr[j] != 0 {
					ok = false
					break
				}
			}
			if !ok {
				break
			}
		}

		if ok {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
