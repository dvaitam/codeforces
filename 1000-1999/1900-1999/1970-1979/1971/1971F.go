package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func countInside(R int) int64 {
	if R <= 0 {
		return 0
	}
	r2 := int64(R) * int64(R)
	var cnt int64
	for x := 0; int64(x)*int64(x) < r2; x++ {
		remain := r2 - int64(x*x) - 1
		if remain < 0 {
			remain = 0
		}
		y := int64(math.Sqrt(float64(remain)))
		if x == 0 {
			cnt += 2*y + 1
		} else {
			cnt += 2 * (2*y + 1)
		}
	}
	return cnt
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
		var r int
		fmt.Fscan(in, &r)
		ans := countInside(r+1) - countInside(r)
		fmt.Fprintln(out, ans)
	}
}
