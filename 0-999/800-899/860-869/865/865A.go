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

	var A int
	if _, err := fmt.Fscan(in, &A); err != nil {
		return
	}

	if A == 1 {
		fmt.Fprintln(out, "1 1")
		fmt.Fprintln(out, 1)
		return
	}

	N := 10 * (A - 1)
	fmt.Fprintln(out, N, 2)
	fmt.Fprintln(out, "1 10")
}
