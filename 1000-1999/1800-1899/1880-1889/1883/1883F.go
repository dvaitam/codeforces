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

		first := make([]bool, n)
		seen := make(map[int]struct{})
		for i, v := range a {
			if _, ok := seen[v]; !ok {
				first[i] = true
				seen[v] = struct{}{}
			}
		}

		last := make([]bool, n)
		seen = make(map[int]struct{})
		for i := n - 1; i >= 0; i-- {
			v := a[i]
			if _, ok := seen[v]; !ok {
				last[i] = true
				seen[v] = struct{}{}
			}
		}

		suf := make([]int, n+1)
		for i := n - 1; i >= 0; i-- {
			if last[i] {
				suf[i] = suf[i+1] + 1
			} else {
				suf[i] = suf[i+1]
			}
		}

		var ans int64
		for i := 0; i < n; i++ {
			if first[i] {
				ans += int64(suf[i])
			}
		}

		fmt.Fprintln(out, ans)
	}
}
