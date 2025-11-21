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
		var n int
		var s string
		fmt.Fscan(in, &n)
		fmt.Fscan(in, &s)

		cnt1 := 0
		cnt0 := 0
		first1 := -1
		last1 := -1
		for i := 0; i < n; i++ {
			if s[i] == '1' {
				cnt1++
				if first1 == -1 {
					first1 = i
				}
				last1 = i
			} else {
				cnt0++
			}
		}

		if cnt1 == 0 || cnt0 == 0 {
			fmt.Fprintln(out, 0)
			continue
		}

		blockLen := last1 - first1 + 1
		ans := blockLen - cnt1
		if ans < 0 {
			ans = 0
		}
		fmt.Fprintln(out, ans)
	}
}

