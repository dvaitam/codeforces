package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const eps = 1e-9

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		var px, py, qx, qy float64
		fmt.Fscan(in, &px, &py, &qx, &qy)
		sum := 0.0
		maxVal := 0.0
		for i := 0; i < n; i++ {
			var a float64
			fmt.Fscan(in, &a)
			sum += a
			if a > maxVal {
				maxVal = a
			}
		}
		dist := math.Hypot(px-qx, py-qy)
		low := 2*maxVal - sum
		if low < 0 {
			low = 0
		}
		if dist+eps >= low && dist <= sum+eps {
			fmt.Fprintln(out, "Yes")
		} else {
			fmt.Fprintln(out, "No")
		}
	}
}
