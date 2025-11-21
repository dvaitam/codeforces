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
		var n, m int
		fmt.Fscan(in, &n, &m)
		a := make([]int, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &a[i])
		}

		resets := 0
		firstResetIdx := -1
		for i := 1; i < m; i++ {
			if a[i] == 1 {
				resets++
				if firstResetIdx == -1 {
					firstResetIdx = i
				}
			}
		}

		ans := 0
		if resets == 0 {
			last := a[m-1]
			if last <= n {
				ans = n - last + 1
			}
		} else {
			startSeg := a[firstResetIdx-1]
			if startSeg <= n && startSeg+resets <= n {
				ans = 1
			}
		}

		fmt.Fprintln(out, ans)
	}
}
