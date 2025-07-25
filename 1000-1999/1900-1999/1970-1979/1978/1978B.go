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
		var n, a, b int64
		fmt.Fscan(reader, &n, &a, &b)
		maxk := n
		if b < maxk {
			maxk = b
		}
		var result int64
		if b <= a {
			result = n * a
		} else {
			k := b - a
			if k > maxk {
				k = maxk
			}
			best := k*b - k*(k-1)/2 + (n-k)*a
			last := maxk*b - maxk*(maxk-1)/2 + (n-maxk)*a
			if last > best {
				best = last
			}
			result = best
		}
		fmt.Fprintln(writer, result)
	}
}
