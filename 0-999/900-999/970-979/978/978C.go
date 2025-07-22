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
	fmt.Fscan(in, &n, &m)
	pref := make([]int64, n)
	for i := 0; i < n; i++ {
		var a int64
		fmt.Fscan(in, &a)
		if i == 0 {
			pref[i] = a
		} else {
			pref[i] = pref[i-1] + a
		}
	}

	idx := 0
	var prev int64
	for ; m > 0; m-- {
		var b int64
		fmt.Fscan(in, &b)
		for idx < n && b > pref[idx] {
			idx++
		}
		if idx == 0 {
			prev = 0
		} else {
			prev = pref[idx-1]
		}
		fmt.Fprintln(out, idx+1, b-prev)
	}
}
