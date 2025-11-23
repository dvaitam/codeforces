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
		var n, m int
		fmt.Fscan(in, &n, &m)

		limit := n << 2
		cnt := make([]int, limit)

		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			cnt[x-1]++
		}

		ans := 0
		lst := -1

		for i := 0; i < limit; i++ {
			if cnt[i] > m {
				if i-lst > ans {
					ans = i - lst
				}
				cnt[i+1] += cnt[i] - 1
			} else {
				lst = i
			}
		}

		fmt.Fprintln(out, ans)
	}
}

