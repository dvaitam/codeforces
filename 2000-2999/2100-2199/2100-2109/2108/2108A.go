package main

import (
	"bufio"
	"fmt"
	"os"
)

func maxAbs(n int) int {
	return n * (n - 1) / 2 * 2
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		res := maxAbs(n) + 1
		fmt.Fprintln(out, res)
	}
}
