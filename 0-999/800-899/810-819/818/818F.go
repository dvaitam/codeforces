package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func edges(n, x int64) int64 {
	if x < 0 {
		return 0
	}
	b := n - x
	if b < 0 {
		return 0
	}
	cap := x * (x - 1) / 2
	if cap > b {
		cap = b
	}
	return b + cap
}
func solve(n int64) int64 {
	if n <= 1 {
		return 0
	}
	t := int64(math.Sqrt(float64(1 + 8*n)))
	for (t+1)*(t+1) <= 1+8*n {
		t++
	}
	for t*t > 1+8*n {
		t--
	}
	x0 := (t - 1) / 2
	if x0*(x0+1) < 2*n {
		x0++
	}
	e0 := edges(n, x0)
	e1 := edges(n, x0-1)
	if e1 > e0 {
		e0 = e1
	}
	return e0
}
func main() {
	r := bufio.NewReader(os.Stdin)
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	var q int
	if _, err := fmt.Fscan(r, &q); err != nil {
		return
	}
	for ; q > 0; q-- {
		var n int64
		fmt.Fscan(r, &n)
		fmt.Fprintln(w, solve(n))
	}
}
