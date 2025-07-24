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
		var n int
		var s string
		fmt.Fscan(reader, &n, &s)

		var total0, total1 int64
		var cur0, cur1 int64
		var max0, max1 int64
		for i := 0; i < n; i++ {
			if s[i] == '0' {
				total0++
				cur0++
				if cur0 > max0 {
					max0 = cur0
				}
				cur1 = 0
			} else {
				total1++
				cur1++
				if cur1 > max1 {
					max1 = cur1
				}
				cur0 = 0
			}
		}
		prod := total0 * total1
		sq0 := max0 * max0
		sq1 := max1 * max1
		ans := prod
		if sq0 > ans {
			ans = sq0
		}
		if sq1 > ans {
			ans = sq1
		}
		fmt.Fprintln(writer, ans)
	}
}
