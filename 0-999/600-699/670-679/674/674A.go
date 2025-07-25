package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}
	ans := make([]int, n+1)
	for l := 0; l < n; l++ {
		cnt := make([]int, n+1)
		best := 0
		for r := l; r < n; r++ {
			c := arr[r]
			cnt[c]++
			if cnt[c] > cnt[best] || (cnt[c] == cnt[best] && c < best) {
				best = c
			}
			ans[best]++
		}
	}
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for i := 1; i <= n; i++ {
		if i > 1 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, ans[i])
	}
	fmt.Fprintln(out)
}
