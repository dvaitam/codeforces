package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD1 = 1000000007
const MOD2 = 666666667
const SZ = 131072

var q [][]int
var e [][]pair
var p []int
var g []int
var st []bool
var n, m int

type pair struct {
	to, w int
}

func dfs(x int) {
	st[x] = true
	for _, edge := range e[x] {
		y, w := edge.to, edge.w
		if p[y] == 0 {
			p[y] = x
			g[y] = w
			dfs(y)
		} else if st[y] && y != p[x] {
			t := []int{w}
			curr := x
			for curr != y {
				t = append(t, g[curr])
				curr = p[curr]
			}
			q = append(q, t)
		}
	}
	st[x] = false
}

func wht(a []int, mod int) {
	for i := 1; i < SZ; i <<= 1 {
		for j := 0; j < SZ; j++ {
			if (j & i) != 0 {
				u := a[j-i]
				v := a[j]
				
				val1 := u + v
				if val1 >= mod {
					val1 -= mod
				}
				
				val2 := u - v
				if val2 < 0 {
					val2 += mod
				}
				
				a[j-i] = val1
				a[j] = val2
			}
		}
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	e = make([][]pair, n+1)
	p = make([]int, n+1)
	g = make([]int, n+1)
	st = make([]bool, n+1)

	sve := 0
	for i := 0; i < m; i++ {
		var u, v, w int
		fmt.Fscan(reader, &u, &v, &w)
		e[u] = append(e[u], pair{v, w})
		e[v] = append(e[v], pair{u, w})
		sve ^= w
	}

	p[1] = 1
	dfs(1)

	sol1 := make([]int, SZ)
	sol2 := make([]int, SZ)
	sol1[sve] = 1
	sol2[sve] = 1

	wht(sol1, MOD1)
	wht(sol2, MOD2)

	for _, vec := range q {
		t1 := make([]int, SZ)
		t2 := make([]int, SZ)
		for _, x := range vec {
			t1[x]++
			t2[x]++
		}
		wht(t1, MOD1)
		wht(t2, MOD2)
		for i := 0; i < SZ; i++ {
			sol1[i] = int(int64(sol1[i]) * int64(t1[i]) % MOD1)
			sol2[i] = int(int64(sol2[i]) * int64(t2[i]) % MOD2)
		}
	}

	wht(sol1, MOD1)
	wht(sol2, MOD2)

	invSZ1 := 742744451 

	for i := 0; i < SZ; i++ {
		if sol1[i] != 0 || sol2[i] != 0 {
			ans := int64(sol1[i]) * int64(invSZ1) % MOD1
			fmt.Fprintln(writer, i, ans)
			return
		}
	}
}