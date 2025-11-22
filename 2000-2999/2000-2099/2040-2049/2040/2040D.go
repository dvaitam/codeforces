package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxN = 200000
const maxVal = 2 * maxN

func sieve(limit int) []bool {
	isPrime := make([]bool, limit+1)
	for i := 2; i <= limit; i++ {
		isPrime[i] = true
	}
	for i := 2; i*i <= limit; i++ {
		if isPrime[i] {
			for j := i * i; j <= limit; j += i {
				isPrime[j] = false
			}
		}
	}
	return isPrime
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	isPrime := sieve(maxVal)

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)

		adj := make([][]int, n+1)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}

		parent := make([]int, n+1)
		order := make([]int, 0, n)
		q := make([]int, 0, n)
		root := 1
		parent[root] = -1
		q = append(q, root)
		for len(q) > 0 {
			v := q[0]
			q = q[1:]
			order = append(order, v)
			for _, to := range adj[v] {
				if to == parent[v] {
					continue
				}
				parent[to] = v
				q = append(q, to)
			}
		}

		ans := make([]int, n+1)
		cur := 1
		ans[root] = cur
		cur++
		ok := true
		for i := 1; i < len(order); i++ {
			v := order[i]
			p := parent[v]
			for cur <= 2*n && isPrime[cur-ans[p]] {
				cur++
			}
			if cur > 2*n {
				ok = false
				break
			}
			ans[v] = cur
			cur++
		}

		if !ok {
			fmt.Fprintln(out, -1)
			continue
		}

		for i := 1; i <= n; i++ {
			if i > 1 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, ans[i])
		}
		fmt.Fprintln(out)
	}
}
