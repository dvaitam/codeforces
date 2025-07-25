package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: This implementation is incomplete. It simply assumes the best good string
// is the original string s itself, which is always a valid good string. The
// optimal answer could be smaller, so this solution is not correct for all cases.

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		var s, t string
		fmt.Fscan(in, &s)
		fmt.Fscan(in, &t)
		dist := 0
		for i := 0; i < n; i++ {
			if s[i] != t[i] {
				dist++
			}
		}
		fmt.Fprintln(out, dist)
	}
}
