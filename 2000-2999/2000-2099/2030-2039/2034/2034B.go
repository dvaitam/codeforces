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
		var n64, m, k int64
		fmt.Fscan(in, &n64, &m, &k)
		n := int(n64)
		var s string
		fmt.Fscan(in, &s)

		var ans int64
		for i := 0; i < n; {
			if s[i] == '1' {
				i++
				continue
			}
			j := i
			for j < n && s[j] == '0' {
				j++
			}
			L := int64(j - i)
			if L >= m {
				num := L - (m - 1)
				den := k + m - 1
				ans += (num + den - 1) / den
			}
			i = j
		}
		fmt.Fprintln(out, ans)
	}
}
