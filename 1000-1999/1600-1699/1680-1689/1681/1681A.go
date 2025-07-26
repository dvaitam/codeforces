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
		a := make([]int, n)
		maxA := 0
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
			if a[i] > maxA {
				maxA = a[i]
			}
		}
		var m int
		fmt.Fscan(reader, &m)
		b := make([]int, m)
		maxB := 0
		for i := 0; i < m; i++ {
			fmt.Fscan(reader, &b[i])
			if b[i] > maxB {
				maxB = b[i]
			}
		}
		if maxA >= maxB {
			fmt.Fprintln(writer, "Alice")
		} else {
			fmt.Fprintln(writer, "Bob")
		}
		if maxB >= maxA {
			fmt.Fprintln(writer, "Bob")
		} else {
			fmt.Fprintln(writer, "Alice")
		}
	}
}
