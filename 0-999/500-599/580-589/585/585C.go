package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var x, y int64
	if _, err := fmt.Fscan(in, &x, &y); err != nil {
		return
	}
	if gcd(x, y) != 1 {
		fmt.Fprintln(out, "Impossible")
		return
	}
	type step struct {
		k int64
		c byte
	}
	var res []step
	for x > 1 || y > 1 {
		if x > y {
			k := (x - 1) / y
			res = append(res, step{k, 'A'})
			x -= k * y
		} else {
			k := (y - 1) / x
			res = append(res, step{k, 'B'})
			y -= k * x
		}
	}
	for _, st := range res {
		fmt.Fprintf(out, "%d%c", st.k, st.c)
	}
	fmt.Fprintln(out)
}
