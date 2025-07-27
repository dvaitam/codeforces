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
		var s string
		fmt.Fscan(in, &n, &s)
		seen := make(map[byte]bool)
		prev := byte(0)
		ok := true
		for i := 0; i < n; i++ {
			c := s[i]
			if c != prev {
				if seen[c] {
					ok = false
					break
				}
				seen[c] = true
				prev = c
			}
		}
		if ok {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
