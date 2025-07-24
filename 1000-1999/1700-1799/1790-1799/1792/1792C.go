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
		fmt.Fscan(in, &n)
		p := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &p[i])
		}

		pre := 0
		for pre < n && p[pre] == pre+1 {
			pre++
		}
		suf := 0
		for suf < n-pre && p[n-1-suf] == n-suf {
			suf++
		}
		rem := n - pre - suf
		if rem < 0 {
			rem = 0
		}
		ans := (rem + 1) / 2
		fmt.Fprintln(out, ans)
	}
}
