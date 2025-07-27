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

	var n, q int
	var k int64
	if _, err := fmt.Fscan(reader, &n, &q, &k); err != nil {
		return
	}

	a := make([]int64, n+1) // 1-indexed
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	pref := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		pref[i] = pref[i-1]
		if i >= 2 && i <= n-1 {
			pref[i] += a[i+1] - a[i-1] - 2
		}
	}

	for ; q > 0; q-- {
		var l, r int
		fmt.Fscan(reader, &l, &r)
		if l == r {
			fmt.Fprintln(writer, k-1)
		} else {
			ans := a[l+1] - 2 + k - a[r-1] - 1 + pref[r-1] - pref[l]
			fmt.Fprintln(writer, ans)
		}
	}
}
