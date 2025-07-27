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
		// allocate deque with head/tail in the middle
		dq := make([]int, 2*n+5)
		l := n
		r := n
		var x int
		fmt.Fscan(reader, &x)
		dq[l] = x
		for i := 1; i < n; i++ {
			fmt.Fscan(reader, &x)
			if x < dq[l] {
				l--
				dq[l] = x
			} else {
				r++
				dq[r] = x
			}
		}
		for i := l; i <= r; i++ {
			if i > l {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, dq[i])
		}
		writer.WriteByte('\n')
	}
}
