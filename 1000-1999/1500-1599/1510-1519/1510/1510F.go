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

	var n int
	var l float64
	if _, err := fmt.Fscan(in, &n, &l); err != nil {
		return
	}

	xs := make([]float64, n)
	ys := make([]float64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &xs[i], &ys[i])
	}

	area := 0.0
	perimeter := 0.0
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		area += xs[i]*ys[j] - xs[j]*ys[i]
		dx := xs[j] - xs[i]
		dy := ys[j] - ys[i]
		perimeter += math.Hypot(dx, dy)
	}
	area = math.Abs(area) * 0.5

	extra := l - perimeter
	r := extra / (2 * math.Pi)
	result := area + r*perimeter + math.Pi*r*r

	fmt.Fprintf(out, "%.15f\n", result)
}
