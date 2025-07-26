package main

import (
	"bufio"
	"fmt"
	"os"
)

// Solution for problemB.txt from contest 1848.
// For each color we track the positions of its planks. With sentinels at 0 and
// n+1, we compute all gaps between consecutive positions. By repainting one
// plank we can split the largest gap, reducing it to half (floor). The minimal
// possible maximal distance for this color is the larger of the second largest
// gap and half of the largest gap. We take the minimum of this value over all
// colors.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		pos := make([][]int, k+1)
		for i := 1; i <= n; i++ {
			var c int
			fmt.Fscan(in, &c)
			pos[c] = append(pos[c], i)
		}
		best := n
		for color := 1; color <= k; color++ {
			lst := 0
			g1, g2 := 0, 0
			for _, idx := range pos[color] {
				gap := idx - lst - 1
				if gap >= g1 {
					g2 = g1
					g1 = gap
				} else if gap > g2 {
					g2 = gap
				}
				lst = idx
			}
			gap := n - lst
			if gap >= g1 {
				g2 = g1
				g1 = gap
			} else if gap > g2 {
				g2 = gap
			}
			cand := g2
			if g1/2 > cand {
				cand = g1 / 2
			}
			if cand < best {
				best = cand
			}
		}
		fmt.Fprintln(out, best)
	}
}
