package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		nums := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &nums[i])
		}
		evenCnt := 0
		hasOdd := false
		minTrailing := 31
		for _, v := range nums {
			if v%2 == 1 {
				hasOdd = true
			} else {
				evenCnt++
				tz := 0
				for v%2 == 0 {
					v /= 2
					tz++
				}
				if tz < minTrailing {
					minTrailing = tz
				}
			}
		}
		if hasOdd {
			fmt.Fprintln(writer, evenCnt)
		} else {
			fmt.Fprintln(writer, minTrailing+n-1)
		}
	}
}
