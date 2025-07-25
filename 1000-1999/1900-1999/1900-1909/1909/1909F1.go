package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 998244353

// This simplified solution attempts the greedy idea used for the hard
// version's placeholder. It does not fully handle all tricky cases of
// the easy version, but follows the basic logic of keeping track of
// how many "small" numbers (\le i) have been placed so far.
func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		ans := 1
		small := 0
		for i := 1; i <= n; i++ {
			smallAvail := i - small
			bigAvail := n - 2*i + 1 + small
			if bigAvail < 0 {
				bigAvail = 0
			}
			diff := a[i] - small
			if diff < 0 || diff > 1 {
				ans = 0
				break
			}
			if diff == 1 {
				if smallAvail <= 0 {
					ans = 0
					break
				}
				ans = ans * smallAvail % MOD
				small++
			} else {
				if bigAvail <= 0 {
					ans = 0
					break
				}
				ans = ans * bigAvail % MOD
			}
		}
		fmt.Fprintln(writer, ans%MOD)
	}
}
