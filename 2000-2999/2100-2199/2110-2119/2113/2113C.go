package main

import (
	"bufio"
	"fmt"
	"os"
)

// For a chosen first explosion center s, all gold strictly inside the initial square
// (Chebyshev distance < k) is lost forever; gold on the boundary is collected immediately,
// and every gold cell outside the square can be collected later. The latter works by
// expanding the empty area step by step and collecting gold in non-decreasing distance
// from s – when a gold cell is collected at distance k from the current center, any cell
// strictly inside that explosion is closer to s and has already been processed.
// Therefore the maximum collectible gold equals totalGold - min_{s in empty cells} lost(s),
// where lost(s) is the number of gold cells with Chebyshev distance < k from s.
// Counting lost(s) is just a square query of radius k-1 around s.

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n, m, k int
		fmt.Fscan(in, &n, &m, &k)
		grid := make([]string, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &grid[i])
		}

		// prefix sums of gold cells
		pref := make([][]int, n+1)
		for i := range pref {
			pref[i] = make([]int, m+1)
		}
		totalGold := 0
		for i := 0; i < n; i++ {
			for j, ch := range grid[i] {
				v := 0
				if ch == 'g' {
					v = 1
					totalGold++
				}
				pref[i+1][j+1] = pref[i][j+1] + pref[i+1][j] - pref[i][j] + v
			}
		}

		if totalGold == 0 {
			fmt.Fprintln(out, 0)
			continue
		}

		// helper to query number of gold in rectangle [x1,x2]x[y1,y2], inclusive
		rectSum := func(x1, y1, x2, y2 int) int {
			return pref[x2+1][y2+1] - pref[x1][y2+1] - pref[x2+1][y1] + pref[x1][y1]
		}

		r := k - 1
		minLost := totalGold + 1
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if grid[i][j] != '.' {
					continue
				}
				x1 := i - r
				if x1 < 0 {
					x1 = 0
				}
				y1 := j - r
				if y1 < 0 {
					y1 = 0
				}
				x2 := i + r
				if x2 >= n {
					x2 = n - 1
				}
				y2 := j + r
				if y2 >= m {
					y2 = m - 1
				}
				lost := rectSum(x1, y1, x2, y2)
				if lost < minLost {
					minLost = lost
				}
			}
		}

		ans := totalGold - minLost
		fmt.Fprintln(out, ans)
	}
}
