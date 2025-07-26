package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	fmt.Fscan(reader, &n, &m)
	pairCount := m / 4

	totalMin := 0
	totalMax := 0

	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(reader, &s)
		bright := 0
		for j := 0; j < m; j++ {
			if s[j] == '1' {
				bright++
			}
		}
		// maximum number of disjoint bright-bright pairs
		dpBright := make([]int, m+1)
		for j := 2; j <= m; j++ {
			dpBright[j] = dpBright[j-1]
			cand := dpBright[j-2]
			if s[j-2] == '1' && s[j-1] == '1' {
				cand++
			}
			if cand > dpBright[j] {
				dpBright[j] = cand
			}
		}
		maxBrightPairs := dpBright[m]
		if maxBrightPairs > pairCount {
			maxBrightPairs = pairCount
		}
		minOcc := bright - maxBrightPairs

		// maximum number of pairs that include at least one dark window
		dpDark := make([]int, m+1)
		for j := 2; j <= m; j++ {
			dpDark[j] = dpDark[j-1]
			cand := dpDark[j-2]
			if s[j-2] == '0' || s[j-1] == '0' {
				cand++
			}
			if cand > dpDark[j] {
				dpDark[j] = cand
			}
		}
		maxDarkPairs := dpDark[m]
		if maxDarkPairs > pairCount {
			maxDarkPairs = pairCount
		}
		forcedBrightPairs := pairCount - maxDarkPairs
		if forcedBrightPairs < 0 {
			forcedBrightPairs = 0
		}
		maxOcc := bright - forcedBrightPairs

		totalMin += minOcc
		totalMax += maxOcc
	}

	fmt.Fprintln(writer, totalMin, totalMax)
}
