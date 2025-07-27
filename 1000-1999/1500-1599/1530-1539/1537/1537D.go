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
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n uint64
		fmt.Fscan(reader, &n)
		if n%2 == 1 {
			fmt.Fprintln(writer, "Bob")
			continue
		}
		if n&(n-1) != 0 { // not a power of two
			fmt.Fprintln(writer, "Alice")
			continue
		}
		// n is power of two
		exp := bits.TrailingZeros64(n)
		if exp%2 == 1 {
			fmt.Fprintln(writer, "Bob")
		} else {
			fmt.Fprintln(writer, "Alice")
		}
	}
}
