package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		solve(in, out)
	}
}

func solve(in *bufio.Reader, out *bufio.Writer) {
	var n int
	fmt.Fscan(in, &n)
	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &b[i])
	}

	edges := make([][]int, n+1)
	for i := 0; i < n; i++ {
		x := a[i]
		y := b[i]
		edges[x] = append(edges[x], y)
		edges[y] = append(edges[y], x)
	}

	visited := make([]bool, n+1)
	var ans int64 = 1
	for v := 1; v <= n; v++ {
		if !visited[v] {
			stack := []int{v}
			visited[v] = true
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				for _, nb := range edges[top] {
					if !visited[nb] {
						visited[nb] = true
						stack = append(stack, nb)
					}
				}
			}
			ans = ans * 2 % mod
		}
	}
	fmt.Fprintln(out, ans)
}
