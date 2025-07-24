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
		seen := make(map[string]int)
		found := false
		for i := 0; i < n-1 && !found; i++ {
			pair := s[i : i+2]
			if first, ok := seen[pair]; ok {
				if i-first >= 2 {
					found = true
				}
			} else {
				seen[pair] = i
			}
		}
		if found {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
