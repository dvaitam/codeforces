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

	var q int
	fmt.Fscan(in, &q)
	for ; q > 0; q-- {
		var n int
		fmt.Fscan(in, &n)
		p := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &p[i])
			p[i]--
		}

		ans := make([]int, n)
		visited := make([]bool, n)
		for i := 0; i < n; i++ {
			if visited[i] {
				continue
			}
			curr := i
			cycle := []int{}
			for !visited[curr] {
				visited[curr] = true
				cycle = append(cycle, curr)
				curr = p[curr]
			}
			length := len(cycle)
			for _, idx := range cycle {
				ans[idx] = length
			}
		}

		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, ans[i])
		}
		fmt.Fprintln(out)
	}
}
