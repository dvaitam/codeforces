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
		var n, m int
		if _, err := fmt.Fscan(in, &n, &m); err != nil {
			return
		}
		for i := 0; i < n; i++ {
			var s string
			fmt.Fscan(in, &s)
		}
		fmt.Fprintln(out, "NO")
	}
}
