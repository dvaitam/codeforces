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
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for tc := 0; tc < t; tc++ {
		var n int
		fmt.Fscan(in, &n)
		// initialize permutation 1..n
		ans := make([]int, n)
		for i := 0; i < n; i++ {
			ans[i] = i + 1
		}
		// swap adjacent pairs
		for i := 1; i < n; i += 2 {
			ans[i], ans[i-1] = ans[i-1], ans[i]
		}
		// if odd, swap last two once more
		if n%2 == 1 && n > 1 {
			ans[n-1], ans[n-2] = ans[n-2], ans[n-1]
		}
		// output
		for i, v := range ans {
			if i > 0 {
				out.WriteByte(' ')
			}
			fmt.Fprintf(out, "%d", v)
		}
		out.WriteByte('\n')
	}
}
