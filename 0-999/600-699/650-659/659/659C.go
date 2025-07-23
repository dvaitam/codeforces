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
	var m int64
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	have := make(map[int64]bool, n)
	for i := 0; i < n; i++ {
		var x int64
		fmt.Fscan(in, &x)
		have[x] = true
	}
	var res []int64
	var cost int64 = 0
	for t := int64(1); cost+t <= m; t++ {
		if !have[t] {
			res = append(res, t)
			cost += t
		}
	}
	fmt.Fprintln(out, len(res))
	for i, v := range res {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, v)
	}
	if len(res) > 0 {
		fmt.Fprintln(out)
	}
}
