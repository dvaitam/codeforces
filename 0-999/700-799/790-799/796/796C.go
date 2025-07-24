package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	g := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}

	type pair struct {
		val int64
		idx int
	}
	ps := make([]pair, n)
	for i := 0; i < n; i++ {
		ps[i] = pair{a[i], i}
	}
	sort.Slice(ps, func(i, j int) bool { return ps[i].val > ps[j].val })

	banned := make([]bool, n)
	ans := int64(1<<63 - 1)
	for r := 0; r < n; r++ {
		banned[r] = true
		for _, v := range g[r] {
			banned[v] = true
		}
		other := int64(-1 << 60)
		for _, p := range ps {
			if !banned[p.idx] {
				other = p.val
				break
			}
		}
		neighMax := int64(-1 << 60)
		for _, v := range g[r] {
			if a[v] > neighMax {
				neighMax = a[v]
			}
		}
		cur := a[r]
		if neighMax != -1<<60 && neighMax+1 > cur {
			cur = neighMax + 1
		}
		if other != -1<<60 && other+2 > cur {
			cur = other + 2
		}
		if cur < ans {
			ans = cur
		}
		banned[r] = false
		for _, v := range g[r] {
			banned[v] = false
		}
	}
	fmt.Println(ans)
}
