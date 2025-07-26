package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 998244353

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
	ra, rb := d.Find(a), d.Find(b)
	if ra != rb {
		d.parent[rb] = ra
	}
}

func pow2(n int) []int {
	res := make([]int, n+1)
	res[0] = 1
	for i := 1; i <= n; i++ {
		res[i] = res[i-1] * 2 % MOD
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var r, c int
	fmt.Fscan(in, &r, &c)
	grid := make([]string, r)
	for i := 0; i < r; i++ {
		fmt.Fscan(in, &grid[i])
	}

	k := 0
	rowFixed := make([]int, r)
	colFixed := make([]int, c)
	rowQ := make([]int, r)
	colQ := make([]int, c)
	var edges [][2]int
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			ch := grid[i][j]
			if ch == '?' {
				k++
				rowQ[i]++
				colQ[j]++
				edges = append(edges, [2]int{i, j})
			} else if ch == '1' {
				rowFixed[i] ^= 1
				colFixed[j] ^= 1
			}
		}
	}

	pow := pow2(k)

	if r%2 == 0 && c%2 == 0 {
		fmt.Fprintln(out, pow[k])
		return
	}

	if r%2 == 0 && c%2 == 1 {
		ans := 0
		for p := 0; p < 2; p++ {
			cur := 1
			for i := 0; i < r; i++ {
				if rowQ[i] == 0 {
					if rowFixed[i] != p {
						cur = 0
						break
					}
				} else {
					cur = cur * pow[rowQ[i]-1] % MOD
				}
			}
			ans = (ans + cur) % MOD
		}
		fmt.Fprintln(out, ans)
		return
	}

	if r%2 == 1 && c%2 == 0 {
		ans := 0
		for p := 0; p < 2; p++ {
			cur := 1
			for j := 0; j < c; j++ {
				if colQ[j] == 0 {
					if colFixed[j] != p {
						cur = 0
						break
					}
				} else {
					cur = cur * pow[colQ[j]-1] % MOD
				}
			}
			ans = (ans + cur) % MOD
		}
		fmt.Fprintln(out, ans)
		return
	}

	// r and c are both odd
	dsu := NewDSU(r + c)
	deg := make([]int, r+c)
	for _, e := range edges {
		i, j := e[0], e[1]
		dsu.Union(i, r+j)
		deg[i]++
		deg[r+j]++
	}
	compVerts := make(map[int][]int)
	for v := 0; v < r+c; v++ {
		root := dsu.Find(v)
		compVerts[root] = append(compVerts[root], v)
	}
	compEdgeCnt := make(map[int]int)
	for _, e := range edges {
		root := dsu.Find(e[0])
		compEdgeCnt[root]++
	}
	ans := 0
	for p := 0; p < 2; p++ {
		rank := 0
		ok := true
		for root, verts := range compVerts {
			xor := 0
			for _, v := range verts {
				if v < r {
					xor ^= p ^ rowFixed[v]
				} else {
					xor ^= p ^ colFixed[v-r]
				}
			}
			if compEdgeCnt[root] == 0 {
				if xor%2 == 1 {
					ok = false
					break
				}
				rank++
			} else {
				if xor%2 == 1 {
					ok = false
					break
				}
				rank += len(verts) - 1
			}
		}
		if ok {
			free := k - rank
			if free >= 0 {
				ans = (ans + pow[free]) % MOD
			}
		}
	}
	fmt.Fprintln(out, ans)
}
