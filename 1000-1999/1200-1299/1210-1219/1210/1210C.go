package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007

type pair struct {
	g     int64
	count int64
}

type state struct {
	node   int
	parent int
	list   []pair
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	x := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &x[i])
	}
	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var a, b int
		fmt.Fscan(reader, &a, &b)
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}

	var ans int64
	stack := []state{{node: 1, parent: 0, list: []pair{{g: x[1], count: 1}}}}

	for len(stack) > 0 {
		st := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, p := range st.list {
			ans = (ans + (p.g%mod)*(p.count%mod)) % mod
		}
		for _, to := range adj[st.node] {
			if to == st.parent {
				continue
			}
			nl := make([]pair, 1)
			nl[0] = pair{g: x[to], count: 1}
			for _, p := range st.list {
				g := gcd(p.g, x[to])
				if nl[len(nl)-1].g == g {
					nl[len(nl)-1].count += p.count
				} else {
					nl = append(nl, pair{g: g, count: p.count})
				}
			}
			stack = append(stack, state{node: to, parent: st.node, list: nl})
		}
	}

	fmt.Fprintln(writer, ans%mod)
}
