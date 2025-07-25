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
		var x, y int
		fmt.Fscan(reader, &x, &y)
		d := x ^ y
		ans := 1 << bits.TrailingZeros(uint(d))
		fmt.Fprintln(writer, ans)
	}
}
