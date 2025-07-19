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

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for i := 0; i < t; i++ {
		var n, a, b int64
		fmt.Fscan(in, &n, &a, &b)
		if 2*a < b {
			fmt.Fprintln(out, n*a)
		} else {
			fmt.Fprintln(out, (n/2)*b+(n%2)*a)
		}
	}
}
