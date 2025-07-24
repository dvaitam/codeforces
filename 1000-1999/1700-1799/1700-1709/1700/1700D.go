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

	var n int
	fmt.Fscan(in, &n)
	pref := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		var x int64
		fmt.Fscan(in, &x)
		pref[i] = pref[i-1] + x
	}

	maxTime := int64(0)
	for i := 1; i <= n; i++ {
		t := (pref[i] + int64(i) - 1) / int64(i)
		if t > maxTime {
			maxTime = t
		}
	}

	var q int
	fmt.Fscan(in, &q)
	for ; q > 0; q-- {
		var t int64
		fmt.Fscan(in, &t)
		if t < maxTime {
			fmt.Fprintln(out, -1)
		} else {
			k := (pref[n] + t - 1) / t
			fmt.Fprintln(out, k)
		}
	}
}
