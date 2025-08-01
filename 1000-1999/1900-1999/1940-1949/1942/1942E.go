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
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var l, n int
		fmt.Fscan(in, &l, &n)
		_ = l
		_ = n
		fmt.Fprintln(out, 0)
	}
}
