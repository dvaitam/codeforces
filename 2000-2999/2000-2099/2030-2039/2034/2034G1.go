package main

import (
	"bufio"
	"fmt"
	"os"
)

type Interval struct {
	l, r, id int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		intervals := make([]Interval, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &intervals[i].l, &intervals[i].r)
			intervals[i].id = i
		}

		// Try k colors greedily
		for k := 1; ; k++ {
			colors := make([]int, n)
			if assignColors(intervals, colors, k) {
				fmt.Fprintln(out, k)
				for i := 0; i < n; i++ {
					if i > 0 {
						fmt.Fprint(out, " ")
					}
					fmt.Fprint(out, colors[i])
				}
				fmt.Fprintln(out)
				break
			}
		}
	}
}

func assignColors(intervals []Interval, colors []int, k int) bool {
	n := len(intervals)
	if k == 1 {
		colors[0] = 1
		for i := 1; i < n; i++ {
			if intervals[i].l <= intervals[i-1].r {
				return false
			}
			colors[i] = 1
		}
		return true
	}
	// greedy: assign colors based on parity
	for i := range colors {
		colors[i] = i%k + 1
	}
	return true
}
