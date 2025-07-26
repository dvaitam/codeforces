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
		var k int
		fmt.Fscan(reader, &n, &k)
		var s string
		fmt.Fscan(reader, &s)

		counts := make([]int, n)
		origK := k
		parity := origK % 2
		res := []byte(s)

		for i := 0; i < n-1; i++ {
			bit := int(res[i] - '0')
			bit ^= parity
			if bit == 0 && k > 0 {
				counts[i] = 1
				k--
				bit ^= 1
			}
			res[i] = byte(bit + '0')
		}

		counts[n-1] = k
		last := int(res[n-1] - '0')
		last ^= parity
		if k%2 == 1 {
			last ^= 1
		}
		res[n-1] = byte(last + '0')

		fmt.Fprintln(writer, string(res))
		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(writer, " ")
			}
			fmt.Fprint(writer, counts[i])
		}
		fmt.Fprintln(writer)
	}
}
