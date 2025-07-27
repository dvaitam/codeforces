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

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var h, c, t int64
		fmt.Fscan(in, &h, &c, &t)
		if t >= h {
			fmt.Fprintln(out, 1)
			continue
		}
		if 2*t <= h+c {
			fmt.Fprintln(out, 2)
			continue
		}
		k := (h - t) / (2*t - h - c)
		n1 := 2*k + 1
		diff1 := math.Abs(float64((k+1)*h+k*c)/float64(n1) - float64(t))
		k++
		n2 := 2*k + 1
		diff2 := math.Abs(float64((k+1)*h+k*c)/float64(n2) - float64(t))
		if diff1 <= diff2 {
			fmt.Fprintln(out, n1)
		} else {
			fmt.Fprintln(out, n2)
		}
	}
}
