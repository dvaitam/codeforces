package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		if n == 3 {
			fmt.Fprintln(out, "1 2 3")
			continue
		}
		if n%2 == 1 {
			// use first n-2 numbers starting from 0
			preLen := n - 2
			oddXor, evenXor := 0, 0
			for i := 0; i < preLen; i++ {
				if (i+1)%2 == 1 {
					oddXor ^= i
				} else {
					evenXor ^= i
				}
				fmt.Fprintf(out, "%d ", i)
			}
			diff := oddXor ^ evenXor
			x := 1 << 29
			y := x ^ diff
			fmt.Fprintf(out, "%d %d\n", x, y)
		} else {
			// use first n-3 numbers starting from 0
			preLen := n - 3
			oddXor, evenXor := 0, 0
			for i := 0; i < preLen; i++ {
				if (i+1)%2 == 1 {
					oddXor ^= i
				} else {
					evenXor ^= i
				}
				fmt.Fprintf(out, "%d ", i)
			}
			diff := oddXor ^ evenXor
			x := 1 << 29
			y := 1 << 28
			z := diff ^ x ^ y
			fmt.Fprintf(out, "%d %d %d\n", x, y, z)
		}
	}
}
