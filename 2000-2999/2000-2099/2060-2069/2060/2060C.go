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
		var n, k int
		fmt.Fscan(in, &n, &k)
		freq := make([]int, n+2)
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			if x <= n {
				freq[x]++
			}
		}
		ans := 0
		for x := 1; x <= n; x++ {
			y := k - x
			if y < 1 || y > n {
				continue
			}
			if x > y {
				continue
			}
			if x == y {
				ans += freq[x] / 2
			} else {
				if freq[x] < freq[y] {
					ans += freq[x]
				} else {
					ans += freq[y]
				}
			}
		}
		fmt.Fprintln(out, ans)
	}
}
