package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353
const maxN = 200000 + 5

var inv [maxN]int64

func init() {
	inv[1] = 1
	for i := 2; i < maxN; i++ {
		inv[i] = (mod - (mod/int64(i))*inv[int(mod%int64(i))]%mod) % mod
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)

		adj := make([][]int, n+1)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(reader, &u, &v)
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}

		parent := make([]int, n+1)
		children := make([][]int, n+1)
		order := make([]int, 0, n)

		stack := []int{1}
		parent[1] = 0
		for len(stack) > 0 {
			v := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			order = append(order, v)
			for _, to := range adj[v] {
				if to == parent[v] {
					continue
				}
				parent[to] = v
				children[v] = append(children[v], to)
				stack = append(stack, to)
			}
		}

		distLeaf := make([]int, n+1)
		for i := len(order) - 1; i >= 0; i-- {
			v := order[i]
			if len(children[v]) == 0 {
				if v == 1 {
					distLeaf[v] = 0
				} else {
					distLeaf[v] = 0
				}
				continue
			}
			minChild := distLeaf[children[v][0]]
			for _, to := range children[v][1:] {
				if distLeaf[to] < minChild {
					minChild = distLeaf[to]
				}
			}
			distLeaf[v] = minChild + 1
		}

		prob := make([]int64, n+1)
		prob[1] = 1

		for _, v := range order {
			for _, to := range children[v] {
				k := distLeaf[to]
				if k == 0 {
					prob[to] = 0
				} else {
					ratio := (int64(k) * inv[k+1]) % mod
					prob[to] = prob[v] * ratio % mod
				}
			}
		}

		for i := 1; i <= n; i++ {
			if i > 1 {
				fmt.Fprint(writer, " ")
			}
			fmt.Fprint(writer, prob[i]%mod)
		}
		fmt.Fprintln(writer)
	}
}
