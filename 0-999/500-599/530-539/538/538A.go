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

	var s string
	fmt.Fscan(in, &s)

	target := "CODEFORCES"
	n := len(s)
	m := len(target)
	for i := 0; i <= m; i++ {
		prefix := target[:i]
		suffix := target[i:]
		if len(s) >= len(target) {
			if len(s) >= len(prefix)+len(suffix) {
				if s[:i] == prefix && s[n-(m-i):] == suffix {
					fmt.Fprintln(out, "YES")
					return
				}
			}
		}
	}
	fmt.Fprintln(out, "NO")
}
