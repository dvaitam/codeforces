package main

import (
	"bufio"
	"fmt"
	"os"
)

func days(a, b, h int64) int64 {
	if h <= a {
		return 1
	}
	// After the first day, each additional day brings (a-b) net progress.
	return 1 + (h-a+(a-b-1))/(a-b)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var q int
		fmt.Fscan(reader, &q)
		low, high := int64(1), int64(1e18)
		for i := 0; i < q; i++ {
			var typ int
			fmt.Fscan(reader, &typ)
			if typ == 1 {
				var a, b, n int64
				fmt.Fscan(reader, &a, &b, &n)
				var L, R int64
				if n == 1 {
					L, R = 1, a
				} else {
					L = (n-1)*a - (n-2)*b + 1
					R = n*a - (n-1)*b
				}
				if L > high || R < low {
					fmt.Fprint(writer, 0)
				} else {
					if L > low {
						low = L
					}
					if R < high {
						high = R
					}
					fmt.Fprint(writer, 1)
				}
			} else {
				var a, b int64
				fmt.Fscan(reader, &a, &b)
				d1 := days(a, b, low)
				d2 := days(a, b, high)
				if d1 == d2 {
					fmt.Fprint(writer, d1)
				} else {
					fmt.Fprint(writer, -1)
				}
			}
			if i+1 == q {
				fmt.Fprintln(writer)
			} else {
				fmt.Fprint(writer, " ")
			}
		}
	}
}
