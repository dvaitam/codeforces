package main

import (
	"bufio"
	"fmt"
	"os"
)

type portal struct {
	l int
	r int
	a int
	b int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		ps := make([]portal, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &ps[i].l, &ps[i].r, &ps[i].a, &ps[i].b)
		}
		var q int
		fmt.Fscan(in, &q)
		xs := make([]int, q)
		for i := 0; i < q; i++ {
			fmt.Fscan(in, &xs[i])
		}
		for _, x := range xs {
			l, r := x, x
			changed := true
			for changed {
				changed = false
				for _, p := range ps {
					if p.r < l || p.l > r {
						continue
					}
					if p.a < l {
						l = p.a
						changed = true
					}
					if p.b > r {
						r = p.b
						changed = true
					}
				}
			}
			fmt.Fprint(out, r, " ")
		}
		fmt.Fprintln(out)
	}
}
