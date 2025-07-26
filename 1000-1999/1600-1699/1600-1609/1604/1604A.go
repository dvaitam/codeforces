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
		var n int
		fmt.Fscan(in, &n)
		ans := 0
		for i := 1; i <= n; i++ {
			var a int
			fmt.Fscan(in, &a)
			if a > i+ans {
				ans += a - (i + ans)
			}
		}
		fmt.Fprintln(out, ans)
	}
}
