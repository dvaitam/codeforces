package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		var m int64
		fmt.Fscan(reader, &n, &m)
		p := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &p[i])
		}
		h := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &h[i])
		}
		adj := make([][]int, n)
		for i := 0; i < n-1; i++ {
			var x, y int
			fmt.Fscan(reader, &x, &y)
			x--
			y--
			adj[x] = append(adj[x], y)
			adj[y] = append(adj[y], x)
		}

		parent := make([]int, n)
		for i := range parent {
			parent[i] = -2
		}
		order := make([]int, 0, n)
		stack := []int{0}
		parent[0] = -1
		for len(stack) > 0 {
			u := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			order = append(order, u)
			for _, v := range adj[u] {
				if v == parent[u] {
					continue
				}
				parent[v] = u
				stack = append(stack, v)
			}
		}

		sumPop := make([]int64, n)
		good := make([]int64, n)
		ok := true
		for i := len(order) - 1; i >= 0 && ok; i-- {
			u := order[i]
			sum := p[u]
			sumGood := int64(0)
			for _, v := range adj[u] {
				if v == parent[u] {
					continue
				}
				sum += sumPop[v]
				sumGood += good[v]
			}
			sumPop[u] = sum
			if (sum+h[u])&1 != 0 {
				ok = false
				break
			}
			g := (sum + h[u]) / 2
			if g < 0 || g > sum {
				ok = false
				break
			}
			if sumGood > g {
				ok = false
				break
			}
			if g-sumGood > p[u] {
				ok = false
				break
			}
			good[u] = g
		}
		if ok {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
