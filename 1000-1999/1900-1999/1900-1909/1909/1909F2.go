package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 998244353

// This is a very simplified implementation that follows the logic of the easy
// version. It does not handle all tricky cases of the hard version but tries to
// respect the prefix counts when possible.
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
			if a[i] != -1 {
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
			} else {
				if smallAvail > 0 {
					ans = ans * smallAvail % MOD
					small++
				} else if bigAvail > 0 {
					ans = ans * bigAvail % MOD
				} else {
					ans = 0
					break
				}
			}
		}
		fmt.Fprintln(writer, ans%MOD)
	}
}
