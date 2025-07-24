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

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	prefixRight := make([]int64, n)
	for i := 1; i < n; i++ {
		diff := a[i-1] - a[i]
		if diff > 0 {
			prefixRight[i] = prefixRight[i-1] + diff
		} else {
			prefixRight[i] = prefixRight[i-1]
		}
	}

	prefixLeft := make([]int64, n)
	for i := n - 2; i >= 0; i-- {
		diff := a[i+1] - a[i]
		if diff > 0 {
			prefixLeft[i] = prefixLeft[i+1] + diff
		} else {
			prefixLeft[i] = prefixLeft[i+1]
		}
	}

	for ; m > 0; m-- {
		var s, t int
		fmt.Fscan(in, &s, &t)
		if s < t {
			fmt.Fprintln(out, prefixRight[t-1]-prefixRight[s-1])
		} else {
			fmt.Fprintln(out, prefixLeft[t-1]-prefixLeft[s-1])
		}
	}
}
