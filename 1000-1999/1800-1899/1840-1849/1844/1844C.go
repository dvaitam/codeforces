package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		sumOdd, sumEven := int64(0), int64(0)
		maxVal := int64(-1 << 63)
		for i := 1; i <= n; i++ {
			var x int64
			fmt.Fscan(in, &x)
			if x > maxVal {
				maxVal = x
			}
			if i%2 == 1 {
				if x > 0 {
					sumOdd += x
				}
			} else {
				if x > 0 {
					sumEven += x
				}
			}
		}
		if sumOdd == 0 && sumEven == 0 {
			fmt.Fprintln(out, maxVal)
		} else if sumOdd > sumEven {
			fmt.Fprintln(out, sumOdd)
		} else {
			fmt.Fprintln(out, sumEven)
		}
	}
}
