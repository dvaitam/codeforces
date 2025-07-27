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
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		var c int64
		fmt.Fscan(in, &n, &c)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		b := make([]int64, n-1)
		for i := 0; i < n-1; i++ {
			fmt.Fscan(in, &b[i])
		}

		curDay := int64(0)
		curMoney := int64(0)
		ans := int64(1 << 60)
		for i := 0; i < n; i++ {
			if curMoney >= c {
				if curDay < ans {
					ans = curDay
				}
			} else {
				need := c - curMoney
				days := curDay + (need+a[i]-1)/a[i]
				if days < ans {
					ans = days
				}
			}
			if i == n-1 {
				break
			}
			if curMoney < b[i] {
				need := b[i] - curMoney
				d := (need + a[i] - 1) / a[i]
				curDay += d
				curMoney += d * a[i]
			}
			curMoney -= b[i]
			curDay++
		}
		fmt.Fprintln(out, ans)
	}
}
