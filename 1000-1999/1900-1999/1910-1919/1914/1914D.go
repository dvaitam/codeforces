package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type pair struct {
	val int64
	idx int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)

		a := make([]pair, n)
		for i := 0; i < n; i++ {
			var x int64
			fmt.Fscan(reader, &x)
			a[i] = pair{val: x, idx: i}
		}
		b := make([]pair, n)
		for i := 0; i < n; i++ {
			var x int64
			fmt.Fscan(reader, &x)
			b[i] = pair{val: x, idx: i}
		}
		c := make([]pair, n)
		for i := 0; i < n; i++ {
			var x int64
			fmt.Fscan(reader, &x)
			c[i] = pair{val: x, idx: i}
		}

		sort.Slice(a, func(i, j int) bool { return a[i].val > a[j].val })
		sort.Slice(b, func(i, j int) bool { return b[i].val > b[j].val })
		sort.Slice(c, func(i, j int) bool { return c[i].val > c[j].val })

		top := 5
		if n < top {
			top = n
		}

		var ans int64
		for i := 0; i < top; i++ {
			for j := 0; j < top; j++ {
				if a[i].idx == b[j].idx {
					continue
				}
				for k := 0; k < top; k++ {
					if a[i].idx == c[k].idx || b[j].idx == c[k].idx {
						continue
					}
					sum := a[i].val + b[j].val + c[k].val
					if sum > ans {
						ans = sum
					}
				}
			}
		}

		fmt.Fprintln(writer, ans)
	}
}
