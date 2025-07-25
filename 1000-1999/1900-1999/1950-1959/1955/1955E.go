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

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		var s string
		fmt.Fscan(reader, &s)
		diff := make([]int, n+1)
		diff[0] = int(s[0] - '0')
		for i := 1; i < n; i++ {
			diff[i] = int(s[i]-'0') ^ int(s[i-1]-'0')
		}
		diff[n] = int(s[n-1] - '0')

		ans := 1
		for k := n; k >= 1; k-- {
			parity := make([]int, k)
			for i := 0; i <= n; i++ {
				parity[i%k] ^= diff[i]
			}
			expect := make([]int, k)
			expect[0] ^= 1
			expect[n%k] ^= 1
			ok := true
			for r := 0; r < k; r++ {
				if parity[r] != expect[r] {
					ok = false
					break
				}
			}
			if ok {
				ans = k
				break
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
