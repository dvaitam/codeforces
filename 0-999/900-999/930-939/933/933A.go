package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	pre1 := make([]int, n+1)
	for i := 1; i <= n; i++ {
		pre1[i] = pre1[i-1]
		if a[i] == 1 {
			pre1[i]++
		}
	}
	suf2 := make([]int, n+2)
	for i := n; i >= 1; i-- {
		suf2[i] = suf2[i+1]
		if a[i] == 2 {
			suf2[i]++
		}
	}

	ans := 0
	for r := 1; r <= n; r++ {
		dp1, dp2 := 0, 0
		for l := r; l >= 1; l-- {
			x := a[l]
			if x == 1 {
				dp1++
				if dp2 < dp1 {
					dp2 = dp1
				}
			} else {
				if dp2+1 > dp1+1 {
					dp2 = dp2 + 1
				} else {
					dp2 = dp1 + 1
				}
			}
			cand := pre1[l-1] + dp2 + suf2[r+1]
			if cand > ans {
				ans = cand
			}
		}
	}

	fmt.Fprintln(writer, ans)
}
