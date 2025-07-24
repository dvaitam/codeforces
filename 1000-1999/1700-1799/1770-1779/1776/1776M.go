package main

import (
	"bufio"
	"fmt"
	"os"
)

type pair struct {
	v int
	p int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	g := make([][]int, n+1)
	deg := make([]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
		deg[u]++
		deg[v]++
	}
	even := make([]bool, n+1)
	odd := make([]bool, n+1)
	q := make([]pair, 0)
	for i := 1; i <= n; i++ {
		if deg[i] <= 1 {
			q = append(q, pair{i, 0})
			even[i] = true
		}
	}
	for h := 0; h < len(q); h++ {
		cur := q[h]
		for _, nb := range g[cur.v] {
			np := cur.p ^ 1
			if np == 0 {
				if !even[nb] {
					even[nb] = true
					q = append(q, pair{nb, np})
				}
			} else {
				if !odd[nb] {
					odd[nb] = true
					q = append(q, pair{nb, np})
				}
			}
		}
	}
	ans := 1
	for i := n; i >= 1; i-- {
		if even[i] || deg[i]%2 == 1 {
			ans = i
			break
		}
	}
	fmt.Println(ans)
}
