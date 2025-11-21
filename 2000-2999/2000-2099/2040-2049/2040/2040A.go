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
		var n, k int
		fmt.Fscan(in, &n, &k)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		ans := -1
		for i := 0; i < n; i++ {
			winnable := true
			for j := 0; j < n; j++ {
				if i == j {
					continue
				}
				if abs(a[i]-a[j])%k == 0 {
					winnable = false
					break
				}
			}
			if winnable {
				ans = i + 1
				break
			}
		}

		if ans == -1 {
			fmt.Fprintln(out, "NO")
		} else {
			fmt.Fprintln(out, "YES")
			fmt.Fprintln(out, ans)
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
