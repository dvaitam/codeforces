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

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m int
		var l, r int
		fmt.Fscan(in, &n, &m, &l, &r)
		a := -l
		b := r
		left := m - b
		if left < 0 {
			left = 0
		}
		if left > a {
			left = a
		}
		lPrime := -left
		rPrime := m - left
		fmt.Fprintf(out, "%d %d\n", lPrime, rPrime)
	}
}
