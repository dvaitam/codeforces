package main

import (
	"bufio"
	"fmt"
	"os"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	p := make([][2]int, k+1)
	for i := 1; i <= k; i++ {
		p[i][0] = 0
		p[i][1] = -1
	}
	xc := (k + 1) / 2
	for i := 0; i < n; i++ {
		var m int
		fmt.Fscan(reader, &m)
		if m > k {
			fmt.Fprintln(writer, -1)
			continue
		}
		dm := -1
		xm, lm, rm := 0, 0, 0
		for x := 1; x <= k; x++ {
			if p[x][0] > p[x][1] {
				d := abs(x-xc)*m + (m/2)*((m+1)/2)
				if dm < 0 || d < dm {
					dm = d
					xm = x
					lm = (k-m)/2 + 1
					rm = lm + m - 1
				}
			} else {
				if p[x][0] > m {
					d := abs(x-xc)*m + (xc-p[x][0])*m + m*(m+1)/2
					if dm < 0 || d < dm {
						dm = d
						xm = x
						lm = p[x][0] - m
						rm = p[x][0] - 1
					}
				}
				if p[x][1] <= k-m {
					d := abs(x-xc)*m + (p[x][1]-xc)*m + m*(m+1)/2
					if dm < 0 || d < dm {
						dm = d
						xm = x
						lm = p[x][1] + 1
						rm = p[x][1] + m
					}
				}
			}
		}
		if dm < 0 {
			fmt.Fprintln(writer, -1)
		} else {
			fmt.Fprintf(writer, "%d %d %d\n", xm, lm, rm)
			if p[xm][0] > p[xm][1] {
				p[xm][0] = lm
				p[xm][1] = rm
			} else if p[xm][0] == rm+1 {
				p[xm][0] = lm
			} else if p[xm][1] == lm-1 {
				p[xm][1] = rm
			}
		}
	}
}
