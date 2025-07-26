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
		var x int64
		fmt.Fscan(reader, &n, &x)
		a := make([]int64, n)
		var maxA int64
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
			if a[i] > maxA {
				maxA = a[i]
			}
		}
		lo := int64(1)
		hi := maxA + x + 1 // exclusive upper bound
		for lo+1 < hi {
			mid := (lo + hi) / 2
			var need int64
			for _, v := range a {
				if mid > v {
					need += mid - v
					if need > x {
						break
					}
				}
			}
			if need <= x {
				lo = mid
			} else {
				hi = mid
			}
		}
		fmt.Fprintln(writer, lo)
	}
}
