package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	val := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &val[i])
	}
	g := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}

	parent := make([]int, n+1)
	parity := make([]int, n+1)
	order := make([]int, 0, n)
	stack := []int{1}
	parent[1] = 0
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, u)
		for _, v := range g[u] {
			if v == parent[u] {
				continue
			}
			parent[v] = u
			parity[v] = parity[u] ^ 1
			stack = append(stack, v)
		}
	}

	cnt := make([][2]int64, n+1)
	size := make([]int, n+1)
	for i := len(order) - 1; i >= 0; i-- {
		u := order[i]
		size[u] = 1
		if parity[u] == 0 {
			cnt[u][0] = 1
		} else {
			cnt[u][1] = 1
		}
		for _, v := range g[u] {
			if v == parent[u] {
				continue
			}
			size[u] += size[v]
			cnt[u][0] += cnt[v][0]
			cnt[u][1] += cnt[v][1]
		}
	}

	totalEven := cnt[1][0]
	totalOdd := cnt[1][1]
	n64 := int64(n)
	F := make([]int64, n+1)
	for x := 1; x <= n; x++ {
		res := n64
		px := parity[x]
		for _, v := range g[x] {
			if v == parent[x] {
				up0 := totalEven - cnt[x][0]
				up1 := totalOdd - cnt[x][1]
				sizeC := n64 - int64(size[x])
				sign := up0 - up1
				if px == 1 {
					sign = up1 - up0
				}
				res += sign * (n64 - sizeC)
			} else {
				sizeC := int64(size[v])
				sign := cnt[v][0] - cnt[v][1]
				if px == 1 {
					sign = cnt[v][1] - cnt[v][0]
				}
				res += sign * (n64 - sizeC)
			}
		}
		F[x] = res
	}

	ans := int64(0)
	for i := 1; i <= n; i++ {
		f := F[i] % mod
		if f < 0 {
			f += mod
		}
		v := val[i] % mod
		if v < 0 {
			v += mod
		}
		ans = (ans + f*v) % mod
	}
	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(out, ans%mod)
	out.Flush()
}
