package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func hIndex(vals []int) int {
	sort.Slice(vals, func(i, j int) bool { return vals[i] > vals[j] })
	h := 0
	for i, v := range vals {
		if v >= i+1 {
			h = i + 1
		} else {
			break
		}
	}
	return h
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	g := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}

	parent := make([]int, n+1)
	order := make([]int, 0, n)
	stack := []int{1}
	parent[1] = -1
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, u)
		for _, v := range g[u] {
			if v == parent[u] {
				continue
			}
			parent[v] = u
			stack = append(stack, v)
		}
	}

	f := make([][]int, n+1)
	var total int64 = int64(n) * int64(n)

	for i := len(order) - 1; i >= 0; i-- {
		u := order[i]
		children := make([]int, 0, len(g[u]))
		for _, v := range g[u] {
			if parent[v] == u {
				children = append(children, v)
			}
		}
		if len(children) == 0 {
			f[u] = []int{0}
			continue
		}
		F2 := len(children)
		arrF := []int{F2}
		total += int64(F2)

		depth := 3
		for {
			vals := make([]int, len(children))
			for idx, v := range children {
				if depth-3 < len(f[v]) {
					vals[idx] = f[v][depth-3]
				} else {
					vals[idx] = 0
				}
			}
			h := hIndex(vals)
			if h == 0 {
				break
			}
			arrF = append(arrF, h)
			total += int64(h)
			depth++
		}
		f[u] = arrF
	}

	fmt.Println(total)
}
