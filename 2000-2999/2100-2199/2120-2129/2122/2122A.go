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
		var n, m int
		fmt.Fscan(in, &n, &m)
		if n == 1 || m == 1 {
			fmt.Fprintln(out, "NO")
		} else {
			fmt.Fprintln(out, "YES")
		}
	}
}
