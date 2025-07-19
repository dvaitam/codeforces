package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	var s string
	if _, err := fmt.Fscan(reader, &n, &s); err != nil {
		return
	}
	totalGold := 0
	for i := 0; i < n; i++ {
		if s[i] == 'G' {
			totalGold++
		}
	}
	if totalGold == 0 {
		fmt.Println(0)
		return
	}

	type run struct{ start, end, length int }
	runs := make([]run, 0, 16)
	for i := 0; i < n; {
		if s[i] == 'G' {
			j := i
			for j < n && s[j] == 'G' {
				j++
			}
			runs = append(runs, run{start: i, end: j, length: j - i})
			i = j
		} else {
			i++
		}
	}

	ans := 0
	// single runs: extend by swapping if possible
	for _, r := range runs {
		ans = max(ans, r.length)
		// can swap one G from elsewhere
		ans = max(ans, min(totalGold, r.length+1))
	}
	// merge adjacent runs separated by one S
	for k := 0; k+1 < len(runs); k++ {
		gap := runs[k+1].start - runs[k].end
		if gap == 1 {
			merged := runs[k].length + runs[k+1].length
			ans = max(ans, min(totalGold, merged+1))
		}
	}
	fmt.Println(ans)
}
