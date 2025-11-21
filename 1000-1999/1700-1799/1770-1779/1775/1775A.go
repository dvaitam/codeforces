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

		allSame := true
		for i := 1; i < n; i++ {
			if s[i] != s[0] {
				allSame = false
				break
			}
		}
		if allSame {
			fmt.Fprintf(out, "%c %s %c\n", s[0], s[1:n-1], s[n-1])
			continue
		}

		found := false
		for i := 1; i < n-1; i++ {
			if s[i] == 'a' {
				fmt.Fprintf(out, "%s %c %s\n", s[:i], 'a', s[i+1:])
				found = true
				break
			}
		}
		if !found {
			fmt.Fprintf(out, "%c %s %c\n", s[0], s[1:n-1], s[n-1])
		}
	}
}
