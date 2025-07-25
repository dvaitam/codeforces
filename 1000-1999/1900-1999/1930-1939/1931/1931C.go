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

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		if n == 1 {
			fmt.Fprintln(out, 0)
			continue
		}
		startVal := a[0]
		prefix := 1
		for prefix < n && a[prefix] == startVal {
			prefix++
		}
		sufStart := 0
		for i := n - 1; i >= 0 && a[i] == startVal; i-- {
			sufStart++
		}
		cost1 := n - prefix - sufStart
		if cost1 < 0 {
			cost1 = 0
		}

		endVal := a[n-1]
		suffix := 1
		for i := n - 2; i >= 0 && a[i] == endVal; i-- {
			suffix++
		}
		preEnd := 0
		for i := 0; i < n && a[i] == endVal; i++ {
			preEnd++
		}
		cost2 := n - preEnd - suffix
		if cost2 < 0 {
			cost2 = 0
		}
		if cost2 < cost1 {
			cost1 = cost2
		}
		fmt.Fprintln(out, cost1)
	}
}
