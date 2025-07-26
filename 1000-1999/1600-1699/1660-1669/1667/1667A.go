package main

import (
	"bufio"
	"fmt"
	"os"
)

// floorDiv returns floor(x/y) for y>0
func floorDiv(x, y int64) int64 {
	if x >= 0 {
		return x / y
	}
	return -((-x + y - 1) / y)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	const inf int64 = 1<<63 - 1
	best := inf
	for i := 0; i < n; i++ {
		var moves int64
		var prev int64
		// to the right of i
		prev = 0
		for j := i + 1; j < n; j++ {
			k := prev/a[j] + 1
			prev = k * a[j]
			moves += k
		}
		// to the left of i
		prev = 0
		for j := i - 1; j >= 0; j-- {
			k := floorDiv(prev-1, a[j])
			prev = k * a[j]
			if k < 0 {
				moves += -k
			} else {
				moves += k
			}
		}
		if moves < best {
			best = moves
		}
	}

	fmt.Fprintln(out, best)
}
