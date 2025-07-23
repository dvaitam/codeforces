package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	var x int64
	if _, err := fmt.Fscan(in, &n, &x); err != nil {
		return
	}
	a, b := x, x
	var cost int64
	for i := 0; i < n; i++ {
		var l, r int64
		fmt.Fscan(in, &l, &r)
		if b < l {
			cost += l - b
			a, b = b, l
		} else if r < a {
			cost += a - r
			a, b = r, a
		} else {
			if a < l {
				a = l
			}
			if b > r {
				b = r
			}
		}
	}
	fmt.Println(cost)
}
