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
		fmt.Fprintln(writer, 2)
		perm := make([]int, 0, n)
		for i := 1; i <= n; i += 2 {
			for x := i; x <= n; x *= 2 {
				perm = append(perm, x)
			}
		}
		for i, v := range perm {
			if i+1 == len(perm) {
				fmt.Fprintln(writer, v)
			} else {
				fmt.Fprint(writer, v, " ")
			}
		}
	}
}
