package main

import (
	"bufio"
	"fmt"
	"os"
)

func smallestDivisor(n int64) int64 {
	if n%2 == 0 {
		return 2
	}
	for i := int64(3); i*i <= n; i += 2 {
		if n%i == 0 {
			return i
		}
	}
	return n
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
		var n, k int64
		fmt.Fscan(reader, &n, &k)
		if n%2 == 0 {
			n += 2 * k
		} else {
			d := smallestDivisor(n)
			n += d
			k--
			n += 2 * k
		}
		fmt.Fprintln(writer, n)
	}
}
