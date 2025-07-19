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
	for T > 0 {
		T--
		var N, K int
		fmt.Fscan(reader, &N, &K)
		if N%2 == 1 {
			x := (N - 1) / 2
			fmt.Fprintln(writer, 1, x, x)
		} else {
			half := N / 2
			if half%2 == 0 {
				fmt.Fprintln(writer, half, half/2, half/2)
			} else {
				fmt.Fprintln(writer, half-1, half-1, 2)
			}
		}
	}
}
