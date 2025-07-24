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
		var a, b, c int
		fmt.Fscan(in, &a, &b, &c)
		if a+b == c || a+c == b || b+c == a {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
