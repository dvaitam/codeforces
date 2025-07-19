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
	for i := 0; i < t; i++ {
		var n int
		fmt.Fscan(reader, &n)
		k0, k1 := 0, 0
		for j := 0; j < n; j++ {
			var a int
			fmt.Fscan(reader, &a)
			if a == 1 {
				k1++
			} else {
				k0++
			}
		}
		if k1 <= n/2 {
			// output zeros, ensure not all zeros to match conditions
			length := n - k1
			if k0 == n {
				length--
				k0--
			}
			fmt.Fprintln(writer, length)
			for j := 0; j < k0; j++ {
				fmt.Fprint(writer, "0 ")
			}
			fmt.Fprintln(writer)
		} else {
			// output ones with even count
			length := k1
			if length%2 != 0 {
				length--
				k1--
			}
			fmt.Fprintln(writer, length)
			for j := 0; j < k1; j++ {
				fmt.Fprint(writer, "1 ")
			}
			fmt.Fprintln(writer)
		}
	}
}
