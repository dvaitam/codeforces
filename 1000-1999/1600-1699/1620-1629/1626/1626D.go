package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

func nextPow2(x int) int {
	if x <= 1 {
		return 1
	}
	return 1 << bits.Len(uint(x-1))
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		freq := make([]int, n+2)
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			if x >= 1 && x <= n {
				freq[x]++
			}
		}
		pre := make([]int, n+1)
		for i := 1; i <= n; i++ {
			pre[i] = pre[i-1] + freq[i]
		}
		powers := []int{}
		for p := 1; p <= 2*n; p <<= 1 {
			powers = append(powers, p)
		}

		ans := int(1e9)
		for i := 0; i <= n; i++ {
			c1 := pre[i]
			cost1 := nextPow2(c1) - c1
			for _, limit := range powers {
				target := pre[i] + limit
				j := sort.Search(len(pre), func(k int) bool { return pre[k] > target }) - 1
				if j < i+1 {
					continue
				}
				c2 := pre[j] - pre[i]
				c3 := n - pre[j]
				cost2 := nextPow2(c2) - c2
				cost3 := nextPow2(c3) - c3
				total := cost1 + cost2 + cost3
				if total < ans {
					ans = total
				}
			}
		}
		fmt.Fprintln(out, ans)
	}
}
