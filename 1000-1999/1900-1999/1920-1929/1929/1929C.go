package main

import (
	"bufio"
	"fmt"
	"os"
)

// minRequired computes the minimal starting coins needed to
// guarantee at least one coin profit per cycle when the casino
// restricts the player from losing more than x times in a row
// and a winning bet multiplies the stake by k.
func minRequired(k, x int64) int64 {
	var total int64
	for i := int64(0); i <= x; i++ {
		// amount to bet after i losses so far
		bet := (total + 1 + (k - 2)) / (k - 1)
		total += bet
	}
	return total
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var k, x, a int64
		fmt.Fscan(in, &k, &x, &a)
		if a >= minRequired(k, x) {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
