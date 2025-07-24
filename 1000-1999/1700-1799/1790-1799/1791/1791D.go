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
		var s string
		fmt.Fscan(in, &s)
		prefix := make([]int, n)
		seen := [26]bool{}
		cnt := 0
		for i := 0; i < n; i++ {
			c := s[i] - 'a'
			if !seen[c] {
				seen[c] = true
				cnt++
			}
			prefix[i] = cnt
		}
		seen = [26]bool{}
		suffix := make([]int, n)
		cnt = 0
		for i := n - 1; i >= 0; i-- {
			c := s[i] - 'a'
			if !seen[c] {
				seen[c] = true
				cnt++
			}
			suffix[i] = cnt
		}
		ans := 0
		for i := 0; i < n-1; i++ {
			val := prefix[i] + suffix[i+1]
			if val > ans {
				ans = val
			}
		}
		fmt.Fprintln(out, ans)
	}
}
