package main

import (
	"bufio"
	"fmt"
	"os"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, a, b int
	if _, err := fmt.Fscan(in, &n, &a, &b); err != nil {
		return
	}
	x := make([]int, n)
	y := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &x[i], &y[i])
	}
	best := 0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			for r1 := 0; r1 < 2; r1++ {
				w1, h1 := x[i], y[i]
				if r1 == 1 {
					w1, h1 = y[i], x[i]
				}
				for r2 := 0; r2 < 2; r2++ {
					w2, h2 := x[j], y[j]
					if r2 == 1 {
						w2, h2 = y[j], x[j]
					}
					if w1+w2 <= a && max(h1, h2) <= b {
						area := w1*h1 + w2*h2
						if area > best {
							best = area
						}
					}
					if max(w1, w2) <= a && h1+h2 <= b {
						area := w1*h1 + w2*h2
						if area > best {
							best = area
						}
					}
				}
			}
		}
	}
	fmt.Println(best)
}
