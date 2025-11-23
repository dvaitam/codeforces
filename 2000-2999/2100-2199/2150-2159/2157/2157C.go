package main

import (
	"bufio"
	"fmt"
	"os"
)

type Query struct {
	c, l, r int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, k, q int
		fmt.Fscan(in, &n, &k, &q)

		qs := make([]Query, q)
		for i := 0; i < q; i++ {
			fmt.Fscan(in, &qs[i].c, &qs[i].l, &qs[i].r)
		}

		minLower := make([]bool, n)
		inC2 := make([]bool, n)

		for _, qu := range qs {
			if qu.c == 1 {
				for i := qu.l - 1; i < qu.r; i++ {
					minLower[i] = true
				}
			}
		}
		for _, qu := range qs {
			if qu.c == 2 {
				for i := qu.l - 1; i < qu.r; i++ {
					inC2[i] = true
				}
			}
		}

		a := make([]int, n)
		if k > 0 {
			var U []int
			for i := 0; i < n; i++ {
				if !minLower[i] {
					U = append(U, i)
				}
			}
			for idx, pos := range U {
				a[pos] = idx % k
			}
		}

		for i := 0; i < n; i++ {
			if minLower[i] {
				a[i] = k + 1
			}
		}

		allowedK := make([]bool, n)
		for i := 0; i < n; i++ {
			if minLower[i] && !inC2[i] {
				allowedK[i] = true
			}
		}

		for _, qu := range qs {
			if qu.c == 1 {
				pos := -1
				for i := qu.l - 1; i < qu.r; i++ {
					if allowedK[i] {
						pos = i
						break
					}
				}
				if pos == -1 {
					pos = qu.l - 1
				}
				a[pos] = k
			}
		}

		for i := 0; i < n; i++ {
			fmt.Fprint(out, a[i])
			if i+1 < n {
				fmt.Fprint(out, " ")
			}
		}
		fmt.Fprintln(out)
	}
}

