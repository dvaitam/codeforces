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
		var n, r int
		fmt.Fscan(in, &n, &r)
		sumPairs := 0
		leftovers := 0
		total := 0
		for i := 0; i < n; i++ {
			var a int
			fmt.Fscan(in, &a)
			sumPairs += a / 2
			leftovers += a % 2
			total += a
		}
		emptySeats := 2*r - total
		if emptySeats < 0 {
			emptySeats = 0
		}
		singles := leftovers
		if singles > emptySeats {
			singles = emptySeats
		}
		ans := 2*sumPairs + singles
		fmt.Fprintln(out, ans)
	}
}
