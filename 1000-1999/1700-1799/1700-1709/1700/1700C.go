package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		var pos, neg int64
		for i := 1; i < n; i++ {
			diff := a[i] - a[i-1]
			if diff > 0 {
				pos += diff
			} else {
				neg += -diff
			}
		}
		ans := pos + neg + int64(math.Abs(float64(a[0]-neg)))
		fmt.Fprintln(out, ans)
	}
}
