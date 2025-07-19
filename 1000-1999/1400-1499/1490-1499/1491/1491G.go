package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	fmt.Fscan(reader, &n)
	c := make([]int, n+1)
	pos := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &c[i])
		pos[c[i]] = i
	}
	// find cycles
	vis := make([]bool, n+1)
	var G [][]int
	for i := 1; i <= n; i++ {
		if !vis[i] {
			var cycle []int
			for u := i; !vis[u]; u = c[u] {
				cycle = append(cycle, u)
				vis[u] = true
			}
			if len(cycle) > 1 {
				G = append(G, cycle)
			}
		}
	}
	// reset visited for two_cycles
	for i := 1; i <= n; i++ {
		vis[i] = false
	}
	// operations
	var ops [][2]int
	// swap function
	var Swap = func(x, y int) {
		ops = append(ops, [2]int{pos[x], pos[y]})
		// swap in c
		cx, cy := pos[x], pos[y]
		c[cx], c[cy] = c[cy], c[cx]
		// swap positions
		pos[x], pos[y] = pos[y], pos[x]
	}
	// two_cycles function
	var two_cycles func(int, int)
	two_cycles = func(p1, p2 int) {
		for !vis[c[p1]] && c[p1] != p2 {
			np1 := c[p1]
			vis[p1] = true
			Swap(p1, np1)
			p1 = np1
		}
		for !vis[c[p2]] && c[p2] != p1 {
			np2 := c[p2]
			vis[p2] = true
			Swap(p2, np2)
			p2 = np2
		}
		Swap(p1, p2)
		vis[p1] = true
		vis[p2] = true
	}
	// process pairs of cycles
	for len(G) >= 2 {
		g1 := G[len(G)-1]
		g2 := G[len(G)-2]
		u := g1[len(g1)-1]
		v := g2[len(g2)-1]
		Swap(u, v)
		two_cycles(u, v)
		G = G[:len(G)-2]
	}
	// one cycle left
	if len(G) == 1 {
		g1 := G[0]
		if len(g1) == n {
			a1, a2, a3 := g1[0], g1[1], g1[2]
			Swap(a1, a2)
			Swap(a1, a3)
			two_cycles(a2, a3)
		} else {
			// find a fixed point
			node := 0
			for i := 1; i <= n; i++ {
				if c[i] == i {
					node = i
					break
				}
			}
			u := g1[len(g1)-1]
			Swap(u, node)
			two_cycles(u, node)
		}
	}
	// output
	fmt.Fprintln(writer, len(ops))
	for _, p := range ops {
		fmt.Fprintln(writer, p[0], p[1])
	}
}
