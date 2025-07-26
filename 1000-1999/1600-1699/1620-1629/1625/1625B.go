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

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		last := make(map[int]int)
		ans := -1
		for i, v := range a {
			if p, ok := last[v]; ok {
				d := i - p
				cand := n - d
				if cand > ans {
					ans = cand
				}
			}
			last[v] = i
		}
		fmt.Fprintln(out, ans)
	}
}
