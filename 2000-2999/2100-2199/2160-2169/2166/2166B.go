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
	for i := 0; i < t; i++ {
		var a, b, n int64
		fmt.Fscan(in, &a, &b, &n)

		ans := int64(2)
		if n <= a/b || a == b {
			// Either we never leave the len=b regime, or tabs always end at 'a'.
			ans = 1
		}

		if i > 0 {
			fmt.Fprint(out, "\n")
		}
		fmt.Fprint(out, ans)
	}
}
