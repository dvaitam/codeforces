package main

import (
	"bufio"
	"fmt"
	"os"
)

// Alice wins iff there is a block of ones touching an end of the string.
// A ones block strictly inside the string needs both adjacent edges chosen as OR.
// Bob, moving after every Alice move, can always respond by marking the other edge
// of any internal block with AND, so Alice can never secure both edges there.
// If a ones block touches an end, only one adjacent edge (or none when all ones)
// is needed; Alice can take it on her first move, guaranteeing a true segment.
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
		fmt.Fscan(in, &n)
		var s string
		fmt.Fscan(in, &s)
		if s[0] == '1' || s[n-1] == '1' {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
