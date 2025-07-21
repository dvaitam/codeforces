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
		if bits.OnesCount64(n) == 1 {
			fmt.Fprintln(writer, 1)
			fmt.Fprintln(writer, n)
		} else {
			cnt := bits.OnesCount64(n)
			fmt.Fprintln(writer, cnt+1)
			for i := uint(0); i < 64; i++ {
				if n&(1<<i) != 0 {
					fmt.Fprint(writer, n^(1<<i), " ")
				}
			}
			fmt.Fprintln(writer, n)
		}
	}
}
