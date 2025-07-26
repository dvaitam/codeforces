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

	var n, k, A int
	if _, err := fmt.Fscan(in, &n, &k, &A); err != nil {
		return
	}

	total := 0
	for i := 0; i < n; i++ {
		var x, y, c int
		fmt.Fscan(in, &x, &y, &c)
		triCost := A * (k - x - y)
		if triCost < c {
			total += triCost
		} else {
			total += c
		}
	}

	if A*k < total {
		fmt.Fprintln(out, A*k)
	} else {
		fmt.Fprintln(out, total)
	}
}
