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
		fmt.Fscan(reader, &n)
		c := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &c[i])
		}
		var k int64
		fmt.Fscan(reader, &k)

		a := make([]int64, n)
		const INF int64 = 1 << 60
		mn := INF
		prefix := int64(0)
		for i := n - 1; i >= 0; i-- {
			if c[i] < mn {
				mn = c[i]
			}
			times := k / mn
			prefix += times
			k -= times * mn
			a[i] = prefix
		}
		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(writer, " ")
			}
			fmt.Fprint(writer, a[i])
		}
		fmt.Fprintln(writer)
	}
}
