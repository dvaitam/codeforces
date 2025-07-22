package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	x := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &x[i])
	}
	y := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &y[i])
	}

	i, j := 0, 0
	sumX, sumY := 0, 0
	count := 0
	for i < n || j < m {
		if sumX == sumY {
			if sumX != 0 {
				count++
				sumX, sumY = 0, 0
				continue
			}
			if i < n {
				sumX += x[i]
				i++
			}
			if j < m {
				sumY += y[j]
				j++
			}
		} else if sumX < sumY {
			if i < n {
				sumX += x[i]
				i++
			}
		} else {
			if j < m {
				sumY += y[j]
				j++
			}
		}
	}
	if sumX == sumY && sumX != 0 {
		count++
	}
	fmt.Fprintln(out, count)
}
