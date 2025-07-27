package main

import (
	"bufio"
	"fmt"
	"os"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
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
		c := make([]int, n)
		a := make([]int, n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &c[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}
		cur := 0
		ans := 0
		for i := 1; i < n; i++ {
			diff := abs(a[i] - b[i])
			if diff == 0 {
				cur = 0
			} else {
				cur = max(diff, cur+c[i-1]-diff)
			}
			ans = max(ans, cur+c[i])
		}
		fmt.Fprintln(out, ans+1)
	}
}
