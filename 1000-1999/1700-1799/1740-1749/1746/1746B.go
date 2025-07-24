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
		a := make([]int, n)
		zeros := 0
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			if a[i] == 0 {
				zeros++
			}
		}
		ans := n // upper bound
		onesPrefix := 0
		zeroSuffix := zeros
		for i := 0; i <= n; i++ {
			if onesPrefix+zeroSuffix < ans {
				ans = onesPrefix + zeroSuffix
			}
			if i == n {
				break
			}
			if a[i] == 1 {
				onesPrefix++
			} else {
				zeroSuffix--
			}
		}
		fmt.Fprintln(out, ans)
	}
}
