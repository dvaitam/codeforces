package main

import (
	"bufio"
	"fmt"
	"os"
)

const BIT = 20

var (
	n   int
	a   []int
	g   [][]int
	pre []int
	ans int64
)

func buildPre(v, p int) {
	if p == 0 {
		pre[v] = a[v]
	} else {
		pre[v] = pre[p] ^ a[v]
	}
	for _, u := range g[v] {
		if u == p {
			continue
		}
		buildPre(u, v)
	}
}

func dfs(v, p int) [BIT][2]int64 {
	var cnt [BIT][2]int64
	for b := 0; b < BIT; b++ {
		parity := (pre[v] >> b) & 1
		cnt[b][parity]++
	}
	for _, u := range g[v] {
		if u == p {
			continue
		}
		sub := dfs(u, v)
		for b := 0; b < BIT; b++ {
			av := (a[v] >> b) & 1
			for px := 0; px < 2; px++ {
				for py := 0; py < 2; py++ {
					bit := px ^ py ^ av
					if bit == 1 {
						ans += int64(1<<b) * cnt[b][px] * sub[b][py]
					}
				}
			}
		}
		for b := 0; b < BIT; b++ {
			cnt[b][0] += sub[b][0]
			cnt[b][1] += sub[b][1]
		}
	}
	return cnt
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a = make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	g = make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}

	pre = make([]int, n+1)
	buildPre(1, 0)

	ans = 0
	for i := 1; i <= n; i++ {
		ans += int64(a[i])
	}

	dfs(1, 0)

	fmt.Fprintln(writer, ans)
}
