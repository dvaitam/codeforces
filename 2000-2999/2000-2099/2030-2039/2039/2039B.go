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
		var s string
		fmt.Fscan(in, &s)
		n := len(s)

		if n == 1 {
			fmt.Fprintln(out, -1)
			continue
		}

		found := false
		for i := 0; i < n-1; i++ {
			if s[i] == s[i+1] {
				fmt.Fprintln(out, s[i:i+2])
				found = true
				break
			}
		}
		if found {
			continue
		}

		for i := 0; i < n-2; i++ {
			if s[i] != s[i+1] && s[i] != s[i+2] && s[i+1] != s[i+2] {
				fmt.Fprintln(out, s[i:i+3])
				found = true
				break
			}
		}
		if found {
			continue
		}

		fmt.Fprintln(out, -1)
	}
}
