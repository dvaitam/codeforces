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
		var n, m int64
		fmt.Fscan(in, &n, &m)
		ans := m*(m+1)/2 + m*(n*(n+1)/2-1)
		fmt.Fprintln(out, ans)
	}
}
