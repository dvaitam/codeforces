package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	p := make([]float64, n+1)
	for i := 0; i <= n; i++ {
		fmt.Fscan(reader, &p[i])
	}
	// d[i] = 2*i - n
	d := make([]float64, n+1)
	for i := 0; i <= n; i++ {
		d[i] = float64(2*i - n)
	}
	best := 0.0
	// consider single i where constraint d[i] <= 0, sum x <= 1
	for i := 0; i <= n; i++ {
		if d[i] <= 0 {
			if p[i] > best {
				best = p[i]
			}
		}
	}
	// consider mixing two points i,j to satisfy d*x = 0 with x_i + x_j = 1
	const eps = 1e-12
	for i := 0; i <= n; i++ {
		for j := 0; j <= n; j++ {
			if d[i] > 0 && d[j] < 0 {
				alpha := (-d[j]) / (d[i] - d[j])
				if alpha < -eps || alpha > 1+eps {
					continue
				}
				val := alpha*p[i] + (1-alpha)*p[j]
				if val > best {
					best = val
				}
			}
		}
	}
	fmt.Printf("%.12f\n", best)
}
