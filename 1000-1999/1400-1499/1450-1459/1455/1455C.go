package main

import (
	"bufio"
	"fmt"
	"os"
)

// For each test case, Alice and Bob play optimally in a simplified ping-pong.
// The optimal result is that Alice wins x-1 plays and Bob wins y plays,
// where x and y are their initial stamina (x>=1, y>=1).
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for i := 0; i < t; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		// Alice's wins: x-1, Bob's wins: y
		alice := x - 1
		bob := y
		fmt.Fprintln(out, alice, bob)
	}
}
