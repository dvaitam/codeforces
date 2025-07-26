package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

// This program solves the problem described in problemA.txt for contest 1657A.
// We can move the chip from (0,0) to (x,y) in one operation if the Euclidean
// distance \sqrt{x^2 + y^2} is an integer. Otherwise two operations suffice by
// visiting (x,0) or (0,y) first. The answer is therefore:
//
//	0 if (x,y) = (0,0)
//	1 if x^2 + y^2 is a perfect square
//	2 otherwise.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var x, y int
		fmt.Fscan(in, &x, &y)
		if x == 0 && y == 0 {
			fmt.Fprintln(out, 0)
			continue
		}
		d2 := x*x + y*y
		r := int(math.Sqrt(float64(d2) + 0.5))
		if r*r == d2 {
			fmt.Fprintln(out, 1)
		} else {
			fmt.Fprintln(out, 2)
		}
	}
}
