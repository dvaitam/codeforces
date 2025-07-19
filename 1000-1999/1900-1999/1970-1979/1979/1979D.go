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
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, k int
		var s string
		fmt.Fscan(reader, &n, &k)
		fmt.Fscan(reader, &s)

		// prefix sums of '1's
		pre := make([]int, n+1)
		for i := 0; i < n; i++ {
			pre[i+1] = pre[i] + int(s[i]-'0')
		}
		// f[i] indicates from i to end is valid segment partition
		f := make([]bool, n+1)
		f[n] = true
		for i := n - 1; i >= 0; i-- {
			j := i + k
			if j > n {
				j = n
			}
			ones := pre[j] - pre[i]
			allSame := ones == 0 || ones == j-i
			okNext := f[j]
			okChange := (j == n) || (s[i] != s[j])
			f[i] = okNext && allSame && okChange
		}

		b0 := int(s[0] - '0')
		lastParity := ((n - 1) / k) % 2
		expectedLast := byte('0' + byte(b0^lastParity))
		ans := -1
		for p := 1; p <= n; p++ {
			parity := ((p - 1) / k) % 2
			expected := byte('0' + byte(b0^parity))
			if s[p-1] != expected {
				break
			}
			if f[p] && (p == n || s[p] == expectedLast) {
				ans = p
				break
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
