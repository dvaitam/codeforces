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

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		var s string
		fmt.Fscan(in, &n)
		fmt.Fscan(in, &s)
		r, b := 0, 0
		for _, ch := range s {
			if ch == 'R' {
				r++
			} else if ch == 'B' {
				b++
			}
		}
		if r > b {
			fmt.Fprintln(out, "Alice")
			continue
		}
		if b > r {
			fmt.Fprintln(out, "Bob")
			continue
		}
		// r == b
		isPal := true
		for i := 0; i < n/2; i++ {
			if s[i] != s[n-1-i] {
				isPal = false
				break
			}
		}
		if isPal {
			fmt.Fprintln(out, "Bob")
		} else {
			fmt.Fprintln(out, "Alice")
		}
	}
}
