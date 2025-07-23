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

	var n, b int
	if _, err := fmt.Fscan(in, &n, &b); err != nil {
		return
	}

	q := make([]int64, 0)
	ans := make([]int64, n)
	for i := 0; i < n; i++ {
		var t, d int64
		fmt.Fscan(in, &t, &d)
		for len(q) > 0 && q[0] <= t {
			q = q[1:]
		}
		if len(q) > b {
			ans[i] = -1
			continue
		}
		var start int64
		if len(q) == 0 {
			start = t
		} else {
			start = q[len(q)-1]
		}
		finish := start + d
		q = append(q, finish)
		ans[i] = finish
	}

	for i, v := range ans {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, v)
	}
	fmt.Fprintln(out)
}
