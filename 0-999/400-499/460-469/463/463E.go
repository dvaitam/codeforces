package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	input, _ := io.ReadAll(os.Stdin)
	var offset int
	readInt := func() int {
		for offset < len(input) && (input[offset] < '0' || input[offset] > '9') {
			offset++
		}
		if offset >= len(input) {
			return 0
		}
		res := 0
		for offset < len(input) && input[offset] >= '0' && input[offset] <= '9' {
			res = res*10 + int(input[offset]-'0')
			offset++
		}
		return res
	}

	n := readInt()
	q := readInt()
	if n == 0 {
		return
	}

	val := make([]int, n+1)
	for i := 1; i <= n; i++ {
		val[i] = readInt()
	}

	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		u := readInt()
		v := readInt()
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	const MAXA = 2000005
	spf := make([]int32, MAXA)
	for i := 2; i < MAXA; i++ {
		spf[i] = int32(i)
	}
	for i := 2; i*i < MAXA; i++ {
		if spf[i] == int32(i) {
			for j := i * i; j < MAXA; j += i {
				if spf[j] == int32(j) {
					spf[j] = int32(i)
				}
			}
		}
	}

	depth := make([]int, n+1)
	depth[0] = -1

	var initDFS func(u, p, d int)
	initDFS = func(u, p, d int) {
		depth[u] = d
		for _, v := range adj[u] {
			if v != p {
				initDFS(v, u, d+1)
			}
		}
	}
	initDFS(1, 0, 1)

	ans := make([]int, n+1)
	deepest_node := make([]int, MAXA)

	var dfs func(u, p int)
	dfs = func(u, p int) {
		var prs [8]int32
		var old_nodes [8]int
		count := 0

		best_node := 0
		x := val[u]
		for x > 1 {
			pr := spf[x]
			prs[count] = pr

			node := deepest_node[pr]
			if depth[node] > depth[best_node] {
				best_node = node
			}

			old_nodes[count] = node
			deepest_node[pr] = u
			count++

			for x%int(pr) == 0 {
				x /= int(pr)
			}
		}

		ans[u] = best_node

		for _, v := range adj[u] {
			if v != p {
				dfs(v, u)
			}
		}

		for i := 0; i < count; i++ {
			deepest_node[prs[i]] = old_nodes[i]
		}
	}

	dfs(1, 0)

	out := make([]byte, 0, 1024*1024)
	writeInt := func(x int) {
		if x == -1 {
			out = append(out, '-', '1', '\n')
			return
		}
		if x == 0 {
			out = append(out, '0', '\n')
			return
		}
		var buf [20]byte
		pos := 20
		for x > 0 {
			pos--
			buf[pos] = byte(x%10 + '0')
			x /= 10
		}
		out = append(out, buf[pos:]...)
		out = append(out, '\n')
	}

	for i := 0; i < q; i++ {
		typ := readInt()
		if typ == 1 {
			v := readInt()
			if ans[v] == 0 {
				writeInt(-1)
			} else {
				writeInt(ans[v])
			}
		} else if typ == 2 {
			v := readInt()
			w := readInt()
			val[v] = w
			dfs(1, 0)
		}
	}
	fmt.Print(string(out))
}
