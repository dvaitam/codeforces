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
		var n, m, r, c int64
		fmt.Fscan(in, &n, &m, &r, &c)
		remaining := (n-r+1)*m - c
		boundary := n - r
		ans := remaining + boundary*(m-1)
		fmt.Fprintln(out, ans)
	}
}
