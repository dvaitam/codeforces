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
		var n, m int64
		fmt.Fscan(in, &n, &m)
		if (n+m)%2 == 0 {
			fmt.Fprintln(out, "Tonya")
		} else {
			fmt.Fprintln(out, "Burenka")
		}
	}
}
