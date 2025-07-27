package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 998244353
const maxN = 300000

func modPow2(pre []int) {
	pre[0] = 1
	for i := 1; i <= maxN; i++ {
		pre[i] = (pre[i-1] << 1) % mod
	}
}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1<<20)
	writer := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer writer.Flush()

	pow2 := make([]int, maxN+1)
	modPow2(pow2)

	var t int
	fmt.Fscan(reader, &t)

	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)

		adj := make([][]int, n)
		for i := 0; i < m; i++ {
			var u, v int
			fmt.Fscan(reader, &u, &v)
			u--
			v--
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}

		color := make([]int8, n)
		for i := range color {
			color[i] = -1
		}

		ans := 1
		bad := false
		queue := make([]int, 0)

		for i := 0; i < n && !bad; i++ {
			if color[i] != -1 {
				continue
			}
			// start BFS
			cnt := [2]int{0, 0}
			color[i] = 0
			cnt[0]++
			queue = queue[:0]
			queue = append(queue, i)

			for len(queue) > 0 && !bad {
				v := queue[0]
				queue = queue[1:]
				for _, to := range adj[v] {
					if color[to] == -1 {
						color[to] = 1 - color[v]
						cnt[color[to]]++
						queue = append(queue, to)
					} else if color[to] == color[v] {
						bad = true
						break
					}
				}
			}

			if bad {
				break
			}

			ways := (pow2[cnt[0]] + pow2[cnt[1]]) % mod
			ans = (ans * ways) % mod
		}

		if bad {
			fmt.Fprintln(writer, 0)
		} else {
			fmt.Fprintln(writer, ans)
		}
	}
}