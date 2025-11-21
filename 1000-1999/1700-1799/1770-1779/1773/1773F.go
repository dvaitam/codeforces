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

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, a, b int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	fmt.Fscan(in, &a)
	fmt.Fscan(in, &b)

	if n == 1 {
		if a == b {
			fmt.Fprintln(out, 1)
		} else {
			fmt.Fprintln(out, 0)
		}
		fmt.Fprintf(out, "%d:%d\n", a, b)
		return
	}

	total := a + b
	draws := 0
	if total < n {
		draws = n - total
	}
	nonDraw := n - draws

	x := make([]int, n)
	y := make([]int, n)

	if nonDraw > 0 {
		winNeed := 0
		if a > 0 {
			winNeed = 1
		}
		loseNeed := 0
		if b > 0 {
			loseNeed = 1
		}

		low := max(winNeed, nonDraw-b)
		high := min(nonDraw-loseNeed, a)
		if low > high {
			high = low
		}
		w := high
		if w < low {
			w = low
		}
		if w < 0 {
			w = 0
		}
		if w > nonDraw {
			w = nonDraw
		}
		l := nonDraw - w
		idx := 0
		for i := 0; i < w; i++ {
			x[idx] = 1
			idx++
		}
		for i := 0; i < l && idx < n; i++ {
			y[idx] = 1
			idx++
		}

		remA := a - w
		remB := b - l
		if w > 0 {
			x[0] += remA
		}
		if l > 0 {
			y[w] += remB
		}
	}

	fmt.Fprintln(out, draws)
	for i := 0; i < n; i++ {
		fmt.Fprintf(out, "%d:%d\n", x[i], y[i])
	}
}
