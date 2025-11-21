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
		fmt.Fscan(in, &n)
		var s string
		fmt.Fscan(in, &s)
		if n%2 == 0 {
			oddCounts := make([]int, 26)
			evenCounts := make([]int, 26)
			for i := 0; i < n; i++ {
				idx := int(s[i] - 'a')
				if (i+1)%2 == 1 {
					oddCounts[idx]++
				} else {
					evenCounts[idx]++
				}
			}
			half := n / 2
			maxOdd := 0
			maxEven := 0
			for c := 0; c < 26; c++ {
				if oddCounts[c] > maxOdd {
					maxOdd = oddCounts[c]
				}
				if evenCounts[c] > maxEven {
					maxEven = evenCounts[c]
				}
			}
			ans := (half - maxOdd) + (half - maxEven)
			fmt.Fprintln(out, ans)
			continue
		}

		prefOdd := make([]int, 26)
		prefEven := make([]int, 26)
		suffOdd := make([]int, 26)
		suffEven := make([]int, 26)

		prefOddCount, prefEvenCount := 0, 0
		suffOddCount, suffEvenCount := 0, 0

		for i := 0; i < n; i++ {
			idx := int(s[i] - 'a')
			if (i+1)%2 == 1 {
				suffOdd[idx]++
				suffOddCount++
			} else {
				suffEven[idx]++
				suffEvenCount++
			}
		}

		ans := int(1e9)
		for i := 0; i < n; i++ {
			idx := int(s[i] - 'a')
			pos := i + 1
			if pos%2 == 1 {
				suffOdd[idx]--
				suffOddCount--
			} else {
				suffEven[idx]--
				suffEvenCount--
			}

			maxX, maxY := 0, 0
			for c := 0; c < 26; c++ {
				if val := prefOdd[c] + suffEven[c]; val > maxX {
					maxX = val
				}
				if val := prefEven[c] + suffOdd[c]; val > maxY {
					maxY = val
				}
			}

			oddPrefEvenSuff := prefOddCount + suffEvenCount
			evenPrefOddSuff := prefEvenCount + suffOddCount
			cost := oddPrefEvenSuff - maxX + evenPrefOddSuff - maxY + 1
			if cost < ans {
				ans = cost
			}

			if pos%2 == 1 {
				prefOdd[idx]++
				prefOddCount++
			} else {
				prefEven[idx]++
				prefEvenCount++
			}
		}

		fmt.Fprintln(out, ans)
	}
}
