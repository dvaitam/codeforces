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
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func outcome(h, w, xa, ya, xb, yb int) string {
	if xa >= xb {
		return "Draw"
	}
	dx := xb - xa
	if dx%2 == 1 {
		s := (dx + 1) / 2
		if s > h-xa || s > xb-1 {
			return "Draw"
		}
		la := max(1, ya-(s-1))
		ra := min(w, ya+(s-1))
		lb := max(1, yb-(s-1))
		rb := min(w, yb+(s-1))
		left := max(1, la-1)
		right := min(w, ra+1)
		if lb >= left && rb <= right {
			return "Alice"
		}
		return "Draw"
	}
	s := dx / 2
	if s > h-xa || s > xb-1 {
		return "Draw"
	}
	la := max(1, ya-(s-1))
	ra := min(w, ya+(s-1))
	lb := max(1, yb-(s-1))
	rb := min(w, yb+(s-1))
	aLeft := max(1, la-1)
	aRight := min(w, ra+1)
	bLeft := max(1, lb-1)
	bRight := min(w, rb+1)
	if aLeft >= bLeft && aRight <= bRight {
		return "Bob"
	}
	return "Draw"
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var h, w int
		var xa, ya, xb, yb int
		fmt.Fscan(reader, &h, &w, &xa, &ya, &xb, &yb)
		fmt.Fprintln(writer, outcome(h, w, xa, ya, xb, yb))
	}
}
