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
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		// arrays sized n+2 to store values per node (1-indexed)
		pre := make([]int, n+2)
		minPre := make([]int, n+2)
		maxPre := make([]int, n+2)
		minSub := make([]int, n+2)
		maxSub := make([]int, n+2)

		// initialize root (node 1 with weight 1)
		pre[1] = 1
		minPre[1] = 0
		maxPre[1] = 1
		minSub[1] = 0
		maxSub[1] = 1

		cur := 1
		for i := 0; i < n; i++ {
			var typ string
			fmt.Fscan(in, &typ)
			if typ == "+" {
				var v, x int
				fmt.Fscan(in, &v, &x)
				cur++
				pre[cur] = pre[v] + x
				if pre[cur] < minPre[v] {
					minPre[cur] = pre[cur]
				} else {
					minPre[cur] = minPre[v]
				}
				if pre[cur] > maxPre[v] {
					maxPre[cur] = pre[cur]
				} else {
					maxPre[cur] = maxPre[v]
				}
				// update maxSub and minSub
				diffMax := pre[cur] - minPre[v]
				if diffMax > maxSub[v] {
					maxSub[cur] = diffMax
				} else {
					maxSub[cur] = maxSub[v]
				}
				diffMin := pre[cur] - maxPre[v]
				if diffMin < minSub[v] {
					minSub[cur] = diffMin
				} else {
					minSub[cur] = minSub[v]
				}
			} else if typ == "?" {
				var u, v, k int
				fmt.Fscan(in, &u, &v, &k)
				if k >= minSub[v] && k <= maxSub[v] {
					fmt.Fprintln(out, "YES")
				} else {
					fmt.Fprintln(out, "NO")
				}
			}
		}
	}
}
