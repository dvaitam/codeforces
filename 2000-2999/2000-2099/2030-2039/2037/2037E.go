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
		var n int
		fmt.Fscan(in, &n)
		var s string
		fmt.Fscan(in, &s)

		possible := false
		for i := 0; i+1 < n; i++ {
			if s[i] == '0' && s[i+1] == '1' {
				possible = true
				break
			}
		}

		if possible {
			fmt.Fprintln(out, s)
		} else {
			fmt.Fprintln(out, "IMPOSSIBLE")
		}
	}
}
