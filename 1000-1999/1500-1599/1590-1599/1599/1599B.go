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

	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var N int64
		var A, B int64
		fmt.Fscan(reader, &N)
		fmt.Fscan(reader, &A, &B)
		var dA, dB string
		fmt.Fscan(reader, &dA, &dB)

		if A != B {
			diff := A - B
			if diff < 0 {
				diff = -diff
			}
			if (A < B && diff%2 == 0) || (A > B && diff%2 == 1) {
				fmt.Fprintln(writer, N-1)
			} else {
				fmt.Fprintln(writer, 0)
			}
		} else {
			if A == 0 {
				fmt.Fprintln(writer, N-1)
			} else if A == N-1 {
				fmt.Fprintln(writer, 0)
			} else {
				if dB == "left" {
					fmt.Fprintln(writer, N-1)
				} else {
					fmt.Fprintln(writer, 0)
				}
			}
		}
	}
}
