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
		var s string
		fmt.Fscan(in, &n, &s)
		seen := make(map[byte]bool)
		for i := 0; i < n; i++ {
			seen[s[i]] = true
		}
		ans := n + len(seen)
		fmt.Fprintln(out, ans)
	}
}
