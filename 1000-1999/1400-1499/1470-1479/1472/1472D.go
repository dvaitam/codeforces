package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		sort.Slice(a, func(i, j int) bool { return a[i] > a[j] })
		alice, bob := 0, 0
		for i, v := range a {
			if i%2 == 0 { // Alice's move
				if v%2 == 0 {
					alice += v
				}
			} else { // Bob's move
				if v%2 == 1 {
					bob += v
				}
			}
		}
		if alice > bob {
			fmt.Fprintln(out, "Alice")
		} else if bob > alice {
			fmt.Fprintln(out, "Bob")
		} else {
			fmt.Fprintln(out, "Tie")
		}
	}
}
