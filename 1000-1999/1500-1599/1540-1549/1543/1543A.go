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
		var a, b int64
		fmt.Fscan(in, &a, &b)
		if a == b {
			fmt.Fprintln(out, 0, 0)
			continue
		}
		if a < b {
			a, b = b, a
		}
		diff := a - b
		rem := a % diff
		moves := rem
		if diff-rem < moves {
			moves = diff - rem
		}
		fmt.Fprintln(out, diff, moves)
	}
}
