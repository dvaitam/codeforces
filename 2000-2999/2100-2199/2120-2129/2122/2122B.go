package main

import (
	"bufio"
	"fmt"
	"os"
)

const inf int64 = 1 << 60

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		type pile struct {
			a, b, c, d int64
		}
		piles := make([]pile, n)
		totalLenMismatch := false
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &piles[i].a, &piles[i].b, &piles[i].c, &piles[i].d)
			if piles[i].a+piles[i].b != piles[i].c+piles[i].d {
				totalLenMismatch = true
			}
		}
		if totalLenMismatch {
			// not expected per statement, but safe
			fmt.Fprintln(out, -1)
			continue
		}

		flag := make([]bool, n)
		for i := 0; i < n; i++ {
			if piles[i].a != piles[i].c {
				flag[i] = true
			}
		}
		// Build arr of piles where we need row adjustments
		arrRow := make([]struct{ a, target int64 }, 0, n)
		for i := 0; i < n; i++ {
			if !flag[i] {
				continue
			}
			arrRow = append(arrRow, struct{ a, target int64 }{piles[i].a, piles[i].c})
		}
		dp := make(map[int64]int64)
		dp[0] = 0
		for _, item := range arrRow {
			next := make(map[int64]int64)
			diff := item.target - item.a
			for cost, val := range dp {
				if v, ok := next[cost]; !ok || v < val {
					next[cost] = val
				}
				if diff >= 0 {
					c := cost + diff
					if v, ok := next[c]; !ok || v < val+diff {
						next[c] = val + diff
					}
				}
			}
			dp = next
		}
		ansRow := int64(inf)
		for _, v := range dp {
			if v < ansRow {
				ansRow = v
			}
		}
		if len(dp) == 0 {
			ansRow = 0
		}

		// Similarly for column adjustments
		arrCol := make([]struct{ a, target int64 }, 0, n)
		for i := 0; i < n; i++ {
			if piles[i].b != piles[i].d {
				arrCol = append(arrCol, struct{ a, target int64 }{piles[i].b, piles[i].d})
			}
		}
		dp = make(map[int64]int64)
		dp[0] = 0
		for _, item := range arrCol {
			next := make(map[int64]int64)
			diff := item.target - item.a
			for cost, val := range dp {
				if v, ok := next[cost]; !ok || v < val {
					next[cost] = val
				}
				if diff >= 0 {
					c := cost + diff
					if v, ok := next[c]; !ok || v < val+diff {
						next[c] = val + diff
					}
				}
			}
			dp = next
		}
		ansCol := int64(inf)
		for _, v := range dp {
			if v < ansCol {
				ansCol = v
			}
		}
		if len(dp) == 0 {
			ansCol = 0
		}

		if ansRow >= inf || ansCol >= inf {
			fmt.Fprintln(out, -1)
		} else {
			fmt.Fprintln(out, ansRow+ansCol)
		}
	}
}
