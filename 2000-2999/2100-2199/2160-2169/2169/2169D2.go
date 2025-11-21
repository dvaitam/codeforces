package main

import (
	"bufio"
	"fmt"
	"os"
)

const limit int64 = 1_000_000_000_000

func survivors(n, x, y int64) int64 {
	cur := n
	steps := int64(0)
	for steps < x && cur > 0 {
		q := cur / y
		if q == 0 {
			break
		}
		r := cur - q*y
		t := r/q + 1
		remaining := x - steps
		if t > remaining {
			t = remaining
		}
		cur -= t * q
		steps += t
	}
	return cur
}

func solveCase(x, y, k int64) int64 {
	if y == 1 {
		return -1
	}
	total := survivors(limit, x, y)
	if total < k {
		return -1
	}
	lo, hi := int64(1), limit
	for lo < hi {
		mid := lo + (hi-lo)/2
		if survivors(mid, x, y) >= k {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return lo
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var x, y, k int64
		fmt.Fscan(in, &x, &y, &k)
		fmt.Fprintln(out, solveCase(x, y, k))
	}
}
