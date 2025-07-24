package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	g := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		a--
		b--
		g[a] = append(g[a], b)
		g[b] = append(g[b], a)
	}
	if n%2 == 1 {
		fmt.Println(0)
		return
	}
	color := make([]int, n)
	parent := make([]int, n)
	order := []int{0}
	parent[0] = -1
	for idx := 0; idx < len(order); idx++ {
		v := order[idx]
		for _, to := range g[v] {
			if to == parent[v] {
				continue
			}
			color[to] = color[v] ^ 1
			parent[to] = v
			order = append(order, to)
		}
	}
	cnt0 := make([]int, n)
	cnt1 := make([]int, n)
	for i := n - 1; i >= 0; i-- {
		v := order[i]
		if color[v] == 0 {
			cnt0[v]++
		} else {
			cnt1[v]++
		}
		if parent[v] != -1 {
			cnt0[parent[v]] += cnt0[v]
			cnt1[parent[v]] += cnt1[v]
		}
	}
	total0 := cnt0[0]
	total1 := cnt1[0]
	half := n / 2
	var ans int64
	for _, v := range order[1:] {
		a0 := cnt0[v]
		a1 := cnt1[v]
		b0 := total0 - a0
		b1 := total1 - a1
		if a0+b0 == half {
			ans += int64(a0)*int64(b1) + int64(a1)*int64(b0)
		}
		if a0+b1 == half {
			ans += int64(a0)*int64(b0) + int64(a1)*int64(b1)
		}
	}
	fmt.Println(ans)
}
