package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	n, a, b int
	h       [12]int
	ans     = 1 << 30
	best    [12]int
	cur     [12]int
)

func dfs(d, sum int) {
	if sum >= ans {
		return
	}
	if d == n {
		if h[n-1] < 0 && h[n] < 0 {
			ans = sum
			best = cur
		}
		return
	}
	needed := 0
	if h[d-1] >= 0 {
		needed = h[d-1]/b + 1
	}
	oldPrev, oldCur, oldNext := h[d-1], h[d], h[d+1]
	for i := needed; ; i++ {
		h[d-1] = oldPrev - i*b
		h[d] = oldCur - i*a
		h[d+1] = oldNext - i*b
		cur[d] = i
		dfs(d+1, sum+i)
		if h[d-1] < 0 && h[d] < 0 && h[d+1] < 0 {
			break
		}
	}
	h[d-1] = oldPrev
	h[d] = oldCur
	h[d+1] = oldNext
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fscan(reader, &n, &a, &b)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &h[i])
	}
	dfs(2, 0)
	fmt.Println(ans)
	first := true
	for d := 2; d < n; d++ {
		for j := 0; j < best[d]; j++ {
			if !first {
				fmt.Print(" ")
			}
			fmt.Print(d)
			first = false
		}
	}
	fmt.Println()
}
