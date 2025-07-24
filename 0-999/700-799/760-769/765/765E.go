package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
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
	dist := make([]int, n)
	for i := range dist {
		dist[i] = -1
	}
	q := []int{0}
	dist[0] = 0
	for head := 0; head < len(q); head++ {
		v := q[head]
		for _, to := range g[v] {
			if dist[to] == -1 {
				dist[to] = dist[v] + 1
				q = append(q, to)
			}
		}
	}
	d := 0
	for i := 0; i < n; i++ {
		if len(g[i]) == 1 {
			if d == 0 {
				d = dist[i]
			} else {
				d = gcd(d, dist[i])
			}
		}
	}
	if d == 0 {
		fmt.Println(0)
		return
	}
	for d%2 == 0 {
		d /= 2
	}
	fmt.Println(d)
}
