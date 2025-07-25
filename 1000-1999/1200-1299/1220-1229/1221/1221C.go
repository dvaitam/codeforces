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

	var q int
	fmt.Fscan(reader, &q)
	for ; q > 0; q-- {
		var c, m, x int64
		fmt.Fscan(reader, &c, &m, &x)
		teams := (c + m + x) / 3
		if c < teams {
			teams = c
		}
		if m < teams {
			teams = m
		}
		fmt.Fprintln(writer, teams)
	}
}
