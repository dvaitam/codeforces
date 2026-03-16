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
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		// Odd-indexed (1-based: 1,3,5,...) and even-indexed (2,4,6,...) are independent
		var sumOdd, sumEven int64
		var cntOdd, cntEven int64
		for i := 0; i < n; i++ {
			if i%2 == 0 { // 1-based odd
				sumOdd += a[i]
				cntOdd++
			} else {
				sumEven += a[i]
				cntEven++
			}
		}
		ok := true
		if cntOdd > 0 && sumOdd%cntOdd != 0 {
			ok = false
		}
		if cntEven > 0 && sumEven%cntEven != 0 {
			ok = false
		}
		if ok && cntOdd > 0 && cntEven > 0 && sumOdd/cntOdd != sumEven/cntEven {
			ok = false
		}
		if ok {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
