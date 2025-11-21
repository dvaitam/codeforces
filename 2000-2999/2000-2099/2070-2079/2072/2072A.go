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
	for ; t > 0; t-- {
		var n, k, p int64
		fmt.Fscan(in, &n, &k, &p)

		if k == 0 {
			fmt.Fprintln(out, 0)
			continue
		}

		absK := k
		if absK < 0 {
			absK = -absK
		}
		ops := (absK + p - 1) / p
		if ops <= n {
			fmt.Fprintln(out, ops)
		} else {
			fmt.Fprintln(out, -1)
		}
	}
}
