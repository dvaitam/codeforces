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
		var n, k int64
		fmt.Fscan(reader, &n, &k)
		var pos int64
		if n%2 == 0 {
			pos = (k-1)%n + 1
		} else {
			m := (n - 1) / 2
			extra := (k - 1) / m
			pos = ((k - 1 + extra) % n) + 1
		}
		fmt.Fprintln(writer, pos)
	}
}
