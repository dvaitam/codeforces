package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var V int
	if _, err := fmt.Fscan(in, &V); err != nil {
		return
	}
	n := 250
	for ; V > 0; V-- {
		sum := 0.0
		sumSq := 0.0
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			fx := float64(x)
			sum += fx
			sumSq += fx * fx
		}
		mean := sum / float64(n)
		variance := sumSq/float64(n) - mean*mean
		ratio := 0.0
		if mean > 0 {
			ratio = variance / mean
		}
		var p float64
		if ratio > 2.0 {
			// likely uniform: variance ≈ (P^2 + P) / 3
			p = (math.Sqrt(1+12*variance) - 1) / 2
		} else {
			// likely Poisson: mean ≈ P
			p = mean
		}
		fmt.Println(int(math.Round(p)))
	}
}
