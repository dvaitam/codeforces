package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		v := make([]int64, n)
		t := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &v[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &t[i])
		}
		diff := make([]int64, n)
		var sum int64
		for i := 0; i < n; i++ {
			diff[i] = t[i] - v[i]
			sum += diff[i]
		}
		adj := make([][]int, n)
		for i := 0; i < m; i++ {
			var a, b int
			fmt.Fscan(in, &a, &b)
			a--
			b--
			adj[a] = append(adj[a], b)
			adj[b] = append(adj[b], a)
		}
		color := make([]int, n)
		for i := range color {
			color[i] = -1
		}
		bip := true
		queue := make([]int, 0, n)
		for i := 0; i < n; i++ {
			if color[i] != -1 {
				continue
			}
			color[i] = 0
			queue = append(queue, i)
			for len(queue) > 0 {
				u := queue[0]
				queue = queue[1:]
				for _, w := range adj[u] {
					if color[w] == -1 {
						color[w] = color[u] ^ 1
						queue = append(queue, w)
					} else if color[w] == color[u] {
						bip = false
					}
				}
			}
		}
		if !bip {
			if sum%2 == 0 {
				fmt.Fprintln(out, "YES")
			} else {
				fmt.Fprintln(out, "NO")
			}
			continue
		}
		var diffSum int64
		for i := 0; i < n; i++ {
			if color[i] == 0 {
				diffSum += diff[i]
			} else {
				diffSum -= diff[i]
			}
		}
		if diffSum == 0 {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
