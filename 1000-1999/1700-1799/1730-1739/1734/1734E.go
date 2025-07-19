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
	fmt.Fscan(reader, &n)

	b := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &b[i])
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			x := (((i-j)*j)%n + b[i] + n) % n
			fmt.Fprint(writer, x, " ")
		}
		fmt.Fprintln(writer)
	}
}
