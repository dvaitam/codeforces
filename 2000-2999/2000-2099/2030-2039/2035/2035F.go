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

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, root int
		fmt.Fscan(in, &n, &root)
		root--

		a := make([]int64, n)
		var sumA int64
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			sumA += a[i]
		}

		adj := make([][]int, n)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			u--
			v--
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}

		children := make([][]int, n)
		parent := make([]int, n)
		for i := range parent {
			parent[i] = -1
		}
		stack := []int{root}
		parent[root] = root
		for len(stack) > 0 {
			u := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			for _, v := range adj[u] {
				if v == parent[u] {
					continue
				}
				parent[v] = u
				children[u] = append(children[u], v)
				stack = append(stack, v)
			}
		}

		capacity := make([]int64, n)
		var dfs func(int) (int64, int64, int)
		dfs = func(u int) (int64, int64, int) {
			low, high := int64(0), int64(0)
			parity := 0
			for _, v := range children[u] {
				cl, ch, cp := dfs(v)
				low += cl
				high += ch
				parity ^= cp
			}
			low += a[u]
			high += a[u]
			if a[u]&1 == 1 {
				parity ^= 1
			}
			cap := capacity[u]
			low -= cap
			high += cap
			if cap&1 == 1 {
				parity ^= 1
			}
			return low, high, parity
		}

		check := func(T int64) bool {
			if T < 0 {
				return false
			}
			for i := 0; i < n; i++ {
				idx := int64(i + 1)
				if T < idx {
					capacity[i] = 0
				} else {
					capacity[i] = (T-idx)/int64(n) + 1
				}
			}
			low, high, parity := dfs(root)
			return low <= 0 && 0 <= high && parity == 0
		}

		if check(sumA) {
			fmt.Fprintln(out, sumA)
			continue
		}

		lo := sumA + 1
		if lo < 0 {
			lo = 0
		}
		hi := lo
		if hi == 0 {
			hi = 1
		}
		limit := int64(1) << 60
		for hi < limit && !check(hi) {
			hi <<= 1
		}
		if !check(hi) {
			fmt.Fprintln(out, -1)
			continue
		}

		for lo < hi {
			mid := (lo + hi) / 2
			if check(mid) {
				hi = mid
			} else {
				lo = mid + 1
			}
		}
		fmt.Fprintln(out, lo)
	}
}
