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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	var s string
	fmt.Fscan(reader, &s)

	ones := 0
	zeros := 0
	for _, ch := range s {
		if ch == '1' {
			ones++
		} else {
			zeros++
		}
	}

	N := 1 << n
	L := 1 << ones
	R := N - (1 << zeros) + 1

	for i := L; i <= R; i++ {
		if i > L {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, i)
	}
	fmt.Fprintln(writer)
}
