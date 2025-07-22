package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func calc(arr []int, y1, y2, vert int) int {
	best := int(1<<62 - 1)
	idx := sort.SearchInts(arr, y1)
	if idx < len(arr) {
		t := abs(arr[idx]-y1) + abs(arr[idx]-y2) + vert
		if t < best {
			best = t
		}
	}
	if idx > 0 {
		t := abs(arr[idx-1]-y1) + abs(arr[idx-1]-y2) + vert
		if t < best {
			best = t
		}
	}
	idx = sort.SearchInts(arr, y2)
	if idx < len(arr) {
		t := abs(arr[idx]-y1) + abs(arr[idx]-y2) + vert
		if t < best {
			best = t
		}
	}
	if idx > 0 {
		t := abs(arr[idx-1]-y1) + abs(arr[idx-1]-y2) + vert
		if t < best {
			best = t
		}
	}
	return best
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, cl, ce, v int
	if _, err := fmt.Fscan(in, &n, &m, &cl, &ce, &v); err != nil {
		return
	}
	stairs := make([]int, cl)
	for i := 0; i < cl; i++ {
		fmt.Fscan(in, &stairs[i])
	}
	elevators := make([]int, ce)
	for i := 0; i < ce; i++ {
		fmt.Fscan(in, &elevators[i])
	}
	var q int
	fmt.Fscan(in, &q)

	for ; q > 0; q-- {
		var x1, y1, x2, y2 int
		fmt.Fscan(in, &x1, &y1, &x2, &y2)
		if x1 == x2 {
			fmt.Fprintln(out, abs(y1-y2))
			continue
		}
		dist := abs(x1 - x2)
		best := int(1<<62 - 1)
		if cl > 0 {
			t := calc(stairs, y1, y2, dist)
			if t < best {
				best = t
			}
		}
		if ce > 0 {
			vert := (dist + v - 1) / v
			t := calc(elevators, y1, y2, vert)
			if t < best {
				best = t
			}
		}
		fmt.Fprintln(out, best)
	}
}
