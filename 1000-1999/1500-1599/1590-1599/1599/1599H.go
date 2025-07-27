package main

import (
	"bufio"
	"fmt"
	"os"
)

func query(w *bufio.Writer, r *bufio.Reader, x, y int) int {
	fmt.Fprintf(w, "? %d %d\n", x, y)
	w.Flush()
	var d int
	fmt.Fscan(r, &d)
	return d
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const mx = 1000000000

	d1 := query(out, in, 1, 1)
	d2 := query(out, in, 1, mx)
	d3 := query(out, in, mx, 1)
	d4 := query(out, in, mx, mx)

	s1 := d1 + 2        // x1 + y1
	s2 := d2 - (mx - 1) // x1 - y2
	s3 := d3 - (mx - 1) // y1 - x2
	s4 := 2*mx - d4     // x2 + y2

	// binary search x1 using queries on y=1
	lo, hi := 1, mx
	x1 := 1
	for lo <= hi {
		mid := (lo + hi) / 2
		d := query(out, in, mid, 1)
		if d == s1-mid-1 {
			x1 = mid
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}

	y1 := s1 - x1
	x2 := s2 + s4 - x1
	y2 := x1 - s2

	fmt.Fprintf(out, "! %d %d %d %d\n", x1, y1, x2, y2)
	out.Flush()
}
