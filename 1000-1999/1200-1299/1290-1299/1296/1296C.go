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
		var n int
		var s string
		fmt.Fscan(in, &n)
		fmt.Fscan(in, &s)

		type pair struct{ x, y int }
		pos := pair{0, 0}
		last := map[pair]int{pos: 0}
		bestLen := n + 1
		bestL, bestR := -1, -1

		for i := 1; i <= n; i++ {
			switch s[i-1] {
			case 'L':
				pos.x--
			case 'R':
				pos.x++
			case 'U':
				pos.y++
			case 'D':
				pos.y--
			}
			if prev, ok := last[pos]; ok {
				if i-prev < bestLen {
					bestLen = i - prev
					bestL = prev + 1
					bestR = i
				}
			}
			last[pos] = i
		}

		if bestLen == n+1 {
			fmt.Fprintln(out, -1)
		} else {
			fmt.Fprintln(out, bestL, bestR)
		}
	}
}
