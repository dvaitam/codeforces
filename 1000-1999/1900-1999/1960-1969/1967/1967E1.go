package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m, b0 int
		if _, err := fmt.Fscan(in, &n, &m, &b0); err != nil {
			return
		}
		// TODO: implement actual algorithm
		// Placeholder prints 0 as result.
		fmt.Fprintln(out, 0)
	}
}
