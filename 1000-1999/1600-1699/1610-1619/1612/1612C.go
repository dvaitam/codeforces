package main

import (
	"bufio"
	"fmt"
	"os"
)

func sumFirst(n, k int64) int64 {
	if n <= k {
		return n * (n + 1) / 2
	}
	m := n - k
	p := k - m - 1
	return k*k - p*(p+1)/2
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var k, x int64
		fmt.Fscan(reader, &k, &x)
		if x >= k*k {
			fmt.Fprintln(writer, 2*k-1)
			continue
		}
		low, high := int64(1), 2*k-1
		for low < high {
			mid := (low + high) / 2
			if sumFirst(mid, k) >= x {
				high = mid
			} else {
				low = mid + 1
			}
		}
		fmt.Fprintln(writer, low)
	}
}
