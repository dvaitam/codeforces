package main

import (
	"bufio"
	"fmt"
	"os"
)

func countDistinct(l, r int64) int64 {
	if r-l <= 0 {
		return 0
	}
	diff := r - l
	res := int64(0)
	// gcd values >= l
	start := l
	end := r / 2
	if end > diff {
		end = diff
	}
	if end >= start {
		res += end - start + 1
	}
	// gcd values < l
	maxG := l - 1
	if maxG > diff {
		maxG = diff
	}
	g := int64(1)
	for g <= maxG {
		k := (l - 1) / g
		var nextG int64
		if k == 0 {
			nextG = maxG
		} else {
			nextG = (l - 1) / k
			if nextG > maxG {
				nextG = maxG
			}
		}
		upper := r / (k + 2)
		if upper >= g {
			if upper > nextG {
				upper = nextG
			}
			res += upper - g + 1
		}
		g = nextG + 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var l, r int64
		fmt.Fscan(in, &l, &r)
		ans := countDistinct(l, r)
		fmt.Fprintln(out, ans)
	}
}
