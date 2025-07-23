package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	a := make([]int, n)
	maxA := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
		if a[i] > maxA {
			maxA = a[i]
		}
	}
	g := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}

	check := func(th int) bool {
		vis := make([]bool, n)
		stack := make([]int, 0, n)
		for i := 0; i < n; i++ {
			if !vis[i] && a[i] >= th {
				size := 0
				stack = append(stack[:0], i)
				vis[i] = true
				for len(stack) > 0 {
					v := stack[len(stack)-1]
					stack = stack[:len(stack)-1]
					size++
					for _, to := range g[v] {
						if !vis[to] && a[to] >= th {
							vis[to] = true
							stack = append(stack, to)
						}
					}
				}
				if size >= k {
					return true
				}
			}
		}
		return false
	}

	low, high := 1, maxA
	for low < high {
		mid := (low + high + 1) / 2
		if check(mid) {
			low = mid
		} else {
			high = mid - 1
		}
	}
	fmt.Println(low)
}
