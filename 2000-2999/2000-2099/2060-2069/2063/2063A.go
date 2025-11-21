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
		var l, r int64
		fmt.Fscan(in, &l, &r)
		ans := int64(0)
		if l == 1 {
			ans++
		}
		start := l
		if start < 2 {
			start = 2
		}
		if r >= start {
			ans += r - start
		}
		fmt.Fprintln(out, ans)
	}
}
