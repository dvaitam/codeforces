package main

import (
	"bufio"
	"fmt"
	"os"
)

// DSU represents a disjoint set union structure.
type DSU struct {
	parent []int
}

// NewDSU initializes a DSU of size n.
func NewDSU(n int) *DSU {
	p := make([]int, n)
	for i := range p {
		p[i] = i
	}
	return &DSU{parent: p}
}

// Find returns the representative for x with path compression.
func (d *DSU) Find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.Find(d.parent[x])
	}
	return d.parent[x]
}

// Union merges the sets containing a and b.
func (d *DSU) Union(a, b int) {
	ra := d.Find(a)
	rb := d.Find(b)
	if ra != rb {
		d.parent[ra] = rb
	}
}

// For k=3, we connect positions i with i+3 and i+4. Any letters within the
// same connected component can be rearranged freely by performing allowed
// swaps along a path. Thus, s can be transformed to t iff for every component
// the multiset of letters is the same in s and t.
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
		for i := 0; i+k < n; i++ {
			dsu.Union(i, i+k)
		}
		for i := 0; i+k+1 < n; i++ {
			dsu.Union(i, i+k+1)
		}

		countsS := make(map[int][26]int)
		countsT := make(map[int][26]int)
		for i := 0; i < n; i++ {
			r := dsu.Find(i)
			cs := countsS[r]
			cs[int(s[i]-'a')]++
			countsS[r] = cs
			ct := countsT[r]
			ct[int(t[i]-'a')]++
			countsT[r] = ct
		}

		ok := true
		for r, cs := range countsS {
			ct := countsT[r]
			for j := 0; j < 26; j++ {
				if cs[j] != ct[j] {
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
