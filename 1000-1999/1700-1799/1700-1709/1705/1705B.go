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
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		last := n - 1
		for last >= 0 && a[last] == 0 {
			last--
		}
		var ans int64
		for i := 0; i < last; i++ {
			ans += a[i]
			if a[i] == 0 {
				ans++
			}
		}
		fmt.Fprintln(out, ans)
	}
}
