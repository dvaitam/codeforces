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
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, q int
		fmt.Fscan(reader, &n, &q)
		var sum int64
		var evenCnt, oddCnt int64
		for i := 0; i < n; i++ {
			var x int64
			fmt.Fscan(reader, &x)
			sum += x
			if x%2 == 0 {
				evenCnt++
			} else {
				oddCnt++
			}
		}
		for ; q > 0; q-- {
			var typ int
			var x int64
			fmt.Fscan(reader, &typ, &x)
			if typ == 0 {
				sum += x * evenCnt
				if x%2 != 0 {
					oddCnt += evenCnt
					evenCnt = 0
				}
			} else {
				sum += x * oddCnt
				if x%2 != 0 {
					evenCnt += oddCnt
					oddCnt = 0
				}
			}
			fmt.Fprintln(writer, sum)
		}
	}
}
