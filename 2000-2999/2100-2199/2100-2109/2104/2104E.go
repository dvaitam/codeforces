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

	var n, k int
	fmt.Fscan(in, &n, &k)
	var s string
	fmt.Fscan(in, &s)

	occ := make([]int64, k)
	for _, ch := range s {
		occ[ch-'a']++
	}

	var q int
	fmt.Fscan(in, &q)
	for ; q > 0; q-- {
		var t string
		fmt.Fscan(in, &t)
		ans := int64(0)
		count := make([]int64, k)
		for _, ch := range t {
			count[ch-'a']++
		}
		for i := 0; i < k; i++ {
			if count[i]+ans > occ[i] {
				ans = count[i] - occ[i]
			}
		}
		fmt.Fprintln(out, ans)
	}
}
