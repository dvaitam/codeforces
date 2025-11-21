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

	const maxN = 1000000
	logCeil := make([]int, maxN+2)
	power := 1
	exp := 0
	for i := 1; i <= maxN+1; i++ {
		if i > power {
			power <<= 1
			exp++
		}
		logCeil[i] = exp
	}

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		children := make([][]int, n)
		for v := 1; v < n; v++ {
			var p int
			fmt.Fscan(in, &p)
			p--
			children[p] = append(children[p], v)
		}

		depth := make([]int, n)
		stack := []int{0}
		order := make([]int, 0, n)
		for len(stack) > 0 {
			v := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			order = append(order, v)
			for _, to := range children[v] {
				depth[to] = depth[v] + 1
				stack = append(stack, to)
			}
		}

		leaves := make([]int, n)
		for i := len(order) - 1; i >= 0; i-- {
			v := order[i]
			if len(children[v]) == 0 {
				leaves[v] = 1
			} else {
				sum := 0
				for _, to := range children[v] {
					sum += leaves[to]
				}
				leaves[v] = sum
			}
		}

		ans := 0
		for v := 0; v < n; v++ {
			k := logCeil[leaves[v]]
			if val := depth[v] + k; val > ans {
				ans = val
			}
		}
		fmt.Fprintln(out, ans)
	}
}
