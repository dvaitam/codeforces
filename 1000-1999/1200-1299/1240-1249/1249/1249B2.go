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
	if _, err := fmt.Fscan(in, &q); err != nil {
		return
	}
	for ; q > 0; q-- {
		var n int
		fmt.Fscan(in, &n)
		p := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &p[i])
			p[i]--
		}
		ans := make([]int, n)
		vis := make([]bool, n)
		for i := 0; i < n; i++ {
			if !vis[i] {
				cycle := []int{}
				j := i
				for !vis[j] {
					vis[j] = true
					cycle = append(cycle, j)
					j = p[j]
				}
				l := len(cycle)
				for _, v := range cycle {
					ans[v] = l
				}
			}
		}
		for i := 0; i < n; i++ {
			if i > 0 {
				out.WriteByte(' ')
			}
			fmt.Fprint(out, ans[i])
		}
		out.WriteByte('\n')
	}
}
