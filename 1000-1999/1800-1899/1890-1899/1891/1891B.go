package main

import (
	"bufio"
	"fmt"
	"math/bits"
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
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		xs := make([]int, q)
		for i := 0; i < q; i++ {
			fmt.Fscan(reader, &xs[i])
		}
		add := make([]int64, 31)
		for e := 0; e <= 30; e++ {
			cur := e
			var sum int64
			for _, x := range xs {
				if cur >= x {
					sum += 1 << (x - 1)
					cur = x - 1
				}
			}
			add[e] = sum
		}
		for i := 0; i < n; i++ {
			exp := bits.TrailingZeros64(uint64(arr[i]))
			if exp > 30 {
				exp = 30
			}
			arr[i] += add[exp]
			if i > 0 {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, arr[i])
		}
		writer.WriteByte('\n')
	}
}
