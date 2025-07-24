package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 2520

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	k := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &k[i])
	}

	m := make([]int, n)
	edges := make([][]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &m[i])
		edges[i] = make([]int, m[i])
		for j := 0; j < m[i]; j++ {
			fmt.Fscan(reader, &edges[i][j])
			edges[i][j]--
		}
	}

	total := n * mod
	next := make([]int, total)
	for i := 0; i < n; i++ {
		for r := 0; r < mod; r++ {
			r2 := (r + k[i]) % mod
			if r2 < 0 {
				r2 += mod
			}
			v := edges[i][r2%m[i]]
			next[i*mod+r] = v*mod + r2
		}
	}

	ans := make([]int, total)
	vis := make([]byte, total)

	for s := 0; s < total; s++ {
		if vis[s] != 0 {
			continue
		}
		stack := []int{}
		cur := s
		for vis[cur] == 0 {
			vis[cur] = 1
			stack = append(stack, cur)
			cur = next[cur]
		}
		if vis[cur] == 1 {
			cycle := []int{cur}
			for v := next[cur]; v != cur; v = next[v] {
				cycle = append(cycle, v)
			}
			uniq := make(map[int]struct{})
			for _, v := range cycle {
				uniq[v/mod] = struct{}{}
			}
			val := len(uniq)
			for _, v := range cycle {
				ans[v] = val
				vis[v] = 2
			}
		}
		for i := len(stack) - 1; i >= 0; i-- {
			u := stack[i]
			if vis[u] == 2 {
				continue
			}
			ans[u] = ans[next[u]]
			vis[u] = 2
		}
	}

	var q int
	fmt.Fscan(reader, &q)
	for i := 0; i < q; i++ {
		var x, y int
		fmt.Fscan(reader, &x, &y)
		x--
		r := y % mod
		if r < 0 {
			r += mod
		}
		fmt.Fprintln(writer, ans[x*mod+r])
	}
}
