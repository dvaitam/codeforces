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

	var V int
	if _, err := fmt.Fscan(in, &V); err != nil {
		return
	}
	const n = 250
	data := make([]float64, n)
	for ; V > 0; V-- {
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			data[i] = float64(x)
		}
		mean := 0.0
		for _, v := range data {
			mean += v
		}
		mean /= n
		var sum2 float64
		maxAbs := 0.0
		for _, v := range data {
			d := v - mean
			sum2 += d * d
			if a := math.Abs(v); a > maxAbs {
				maxAbs = a
			}
		}
		varVal := sum2 / n
		if varVal == 0 {
			fmt.Fprintln(out, "uniform")
			continue
		}
		rmax := maxAbs / math.Sqrt(varVal)
		if rmax > 2.2 {
			fmt.Fprintln(out, "poisson")
		} else {
			fmt.Fprintln(out, "uniform")
		}
	}
}
