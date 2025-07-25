package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type pair struct {
	delta int
	depth int
	dist  int
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
		adj := make([][]int, n)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			u--
			v--
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}

		dist := make([][]int, n)
		for i := 0; i < n; i++ {
			dist[i] = make([]int, n)
			for j := 0; j < n; j++ {
				dist[i][j] = -1
			}
			queue := make([]int, 0, n)
			dist[i][i] = 0
			queue = append(queue, i)
			for head := 0; head < len(queue); head++ {
				cur := queue[head]
				for _, nb := range adj[cur] {
					if dist[i][nb] == -1 {
						dist[i][nb] = dist[i][cur] + 1
						queue = append(queue, nb)
					}
				}
			}
		}

		depth := dist[0]
		e1 := 0
		for _, d := range depth {
			if d > e1 {
				e1 = d
			}
		}

		ans := make([]int, n+1)
		for x := 1; x <= n; x++ {
			ans[x] = e1
		}

		arr := make([]pair, n)
		prefix := make([]int, n+1)
		suffix := make([]int, n+1)

		for v := 1; v < n; v++ {
			for u := 0; u < n; u++ {
				arr[u] = pair{
					delta: depth[u] - dist[v][u],
					depth: depth[u],
					dist:  dist[v][u],
				}
			}
			sort.Slice(arr, func(i, j int) bool { return arr[i].delta < arr[j].delta })

			prefix[0] = 0
			for i := 0; i < n; i++ {
				if arr[i].depth > prefix[i] {
					prefix[i+1] = arr[i].depth
				} else {
					prefix[i+1] = prefix[i]
				}
			}
			suffix[n] = 0
			for i := n - 1; i >= 0; i-- {
				if arr[i].dist > suffix[i+1] {
					suffix[i] = arr[i].dist
				} else {
					suffix[i] = suffix[i+1]
				}
			}

			for x := 1; x <= n; x++ {
				idx := sort.Search(n, func(i int) bool { return arr[i].delta > x })
				val := prefix[idx]
				tmp := x + suffix[idx]
				if tmp > val {
					val = tmp
				}
				if val < ans[x] {
					ans[x] = val
				}
			}
		}

		for x := 1; x <= n; x++ {
			if x > 1 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, ans[x])
		}
		fmt.Fprintln(out)
	}
}
