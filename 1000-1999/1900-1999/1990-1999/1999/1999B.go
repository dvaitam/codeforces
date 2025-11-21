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
		var a1, a2, b1, b2 int
		fmt.Fscan(in, &a1, &a2, &b1, &b2)

		sOrders := [][2]int{{a1, a2}, {a2, a1}}
		bOrders := [][2]int{{b1, b2}, {b2, b1}}
		ans := 0

		for _, s := range sOrders {
			for _, b := range bOrders {
				sWins, bWins := 0, 0
				for i := 0; i < 2; i++ {
					if s[i] > b[i] {
						sWins++
					} else if s[i] < b[i] {
						bWins++
					}
				}
				if sWins > bWins {
					ans++
				}
			}
		}

		fmt.Fprintln(out, ans)
	}
}
