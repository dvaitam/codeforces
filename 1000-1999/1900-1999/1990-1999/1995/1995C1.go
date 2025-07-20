package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

const eps = 1e-12 // tolerance for “almost an integer’’

// delta = ceil( log2( log(x) / log(y) ) )  computed robustly
func delta(x, y int64) int64 {
	if x == y { // ratio == 1  →  log2 1 == 0
		return 0
	}
	if x == 1 { // log 1 == 0  →  −Inf after log2; clamp far negative
		return -60
	}

	val := math.Log(float64(x)) / math.Log(float64(y)) // log(x)/log(y)
	val = math.Log2(val)                               // log2(…)

	nearest := math.Round(val)
	if math.Abs(val-nearest) <= eps {
		return int64(nearest) // treat as the exact integer
	}
	return int64(math.Ceil(val)) // proper ceiling for the rest
}

func solve(r *bufio.Reader, w *bufio.Writer) {
	var n, x int64
	if _, err := fmt.Fscan(r, &n, &x); err != nil {
		return
	}
	n-- // we already consumed x

	var ans, temp int64
	for n > 0 {
		var y int64
		fmt.Fscan(r, &y)

		// invalid edge: once seen, whole test‑case is −1
		if y == 1 && x != 1 {
			ans = -1
		}

		if ans != -1 {
			d := delta(x, y)

			// temp = max(0, temp + d)
			if nt := temp + d; nt > 0 {
				temp = nt
			} else {
				temp = 0
			}
			ans += temp
		}

		x, n = y, n-1
	}
	fmt.Fprintln(w, ans)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		solve(in, out)
	}
}
