package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func canPack(pieces []int, n, m int) bool {
	bins := make([]int, m)
	for i := range bins {
		bins[i] = n
	}
	var dfs func(int) bool
	dfs = func(idx int) bool {
		if idx == len(pieces) {
			return true
		}
		piece := pieces[idx]
		for i := 0; i < m; i++ {
			if bins[i] >= piece {
				bins[i] -= piece
				if dfs(idx + 1) {
					return true
				}
				bins[i] += piece
			}
		}
		return false
	}
	return dfs(0)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, a, b int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	fmt.Fscan(in, &a)
	fmt.Fscan(in, &b)

	pieces := []int{a, a, a, a, b, b}
	sort.Slice(pieces, func(i, j int) bool { return pieces[i] > pieces[j] })
	total := 4*a + 2*b
	minBars := (total + n - 1) / n
	for m := minBars; m <= 6; m++ {
		if canPack(pieces, n, m) {
			fmt.Fprintln(out, m)
			return
		}
	}
	fmt.Fprintln(out, 6)
}
