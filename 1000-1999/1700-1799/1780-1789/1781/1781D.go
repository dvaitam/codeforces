package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func isSquare(x int64) bool {
	if x < 0 {
		return false
	}
	r := int64(math.Sqrt(float64(x)))
	if r*r == x || (r+1)*(r+1) == x || (r-1)*(r-1) == x {
		return true
	}
	return false
}

func solve(in *bufio.Reader, out *bufio.Writer) {
	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			var val int64
			fmt.Fscan(in, &val)
			a[i] = val
		}
		ans := 1
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				diff := a[j] - a[i]
				if diff == 0 {
					continue
				}
				for d := int64(1); d*d <= diff; d++ {
					if diff%d != 0 {
						continue
					}
					d1 := d
					d2 := diff / d
					if (d1+d2)%2 != 0 {
						// t and s would not be integers
						continue
					}
					tVal := (d1 + d2) / 2
					x := tVal*tVal - a[j]
					if x < 0 || x > 1e18 {
						continue
					}
					cnt := 0
					for k := 0; k < n; k++ {
						if isSquare(a[k] + x) {
							cnt++
						}
					}
					if cnt > ans {
						ans = cnt
					}
				}
			}
		}
		fmt.Fprintln(out, ans)
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	solve(in, out)
}
