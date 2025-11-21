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
		var n, a, b, c int64
		fmt.Fscan(in, &n, &a, &b, &c)
		total := a + b + c
		full := n / total
		rem := n % total
		day := full * 3
		if rem == 0 {
			fmt.Fprintln(out, day)
			continue
		}
		if rem <= a {
			fmt.Fprintln(out, day+1)
		} else if rem <= a+b {
			fmt.Fprintln(out, day+2)
		} else {
			fmt.Fprintln(out, day+3)
		}
	}
}
