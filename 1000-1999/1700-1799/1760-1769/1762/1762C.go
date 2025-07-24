package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n int
		var s string
		fmt.Fscan(reader, &n)
		fmt.Fscan(reader, &s)
		var same, large int64 = 1, 0
		var ans int64 = 1
		for i := 1; i < n; i++ {
			if s[i] == s[i-1] {
				newSame := same
				newLarge := (same + 2*large) % mod
				same = newSame % mod
				large = newLarge
			} else {
				same = same % mod
				large = 0
			}
			ans = (ans + same + large) % mod
		}
		fmt.Fprintln(writer, ans%mod)
	}
}
