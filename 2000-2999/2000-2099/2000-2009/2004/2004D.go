package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// colorIndex maps portal character to internal id.
func colorIndex(c byte) int {
	switch c {
	case 'B':
		return 0
	case 'G':
		return 1
	case 'R':
		return 2
	case 'Y':
		return 3
	}
	return -1
}

// pairID returns an index in [0,5] for an unordered pair of distinct colors.
// The order of colors does not matter. If colors are the same it returns -1.
func pairID(a, b int) int {
	if a == b {
		return -1
	}
	if a > b {
		a, b = b, a
	}
	switch {
	case a == 0 && b == 1: // B,G
		return 0
	case a == 0 && b == 2: // B,R
		return 1
	case a == 0 && b == 3: // B,Y
		return 2
	case a == 1 && b == 2: // G,R
		return 3
	case a == 1 && b == 3: // G,Y
		return 4
	case a == 2 && b == 3: // R,Y
		return 5
	}
	return -1
}

// minCostThrough returns the minimal value of |x-pos| + |y-pos| over positions in vec.
// vec must be sorted.
func minCostThrough(x, y int, vec []int) int {
	if len(vec) == 0 {
		return int(1e18)
	}
	if x > y {
		x, y = y, x
	}
	best := int(1e18)
	idx := sort.SearchInts(vec, x)
	if idx < len(vec) {
		p := vec[idx]
		v := abs(p-x) + abs(p-y)
		if v < best {
			best = v
		}
	}
	if idx > 0 {
		p := vec[idx-1]
		v := abs(p-x) + abs(p-y)
		if v < best {
			best = v
		}
	}
	return best
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	const inf = int(1e18)
	for ; T > 0; T-- {
		var n, q int
		fmt.Fscan(in, &n, &q)

		colorsA := make([]int, n+1)
		colorsB := make([]int, n+1)
		lists := make([][]int, 6) // positions for each color pair

		for i := 1; i <= n; i++ {
			var s string
			fmt.Fscan(in, &s)
			a := colorIndex(s[0])
			b := colorIndex(s[1])
			colorsA[i] = a
			colorsB[i] = b
			id := pairID(a, b)
			if id >= 0 {
				lists[id] = append(lists[id], i)
			}
		}

		for ; q > 0; q-- {
			var x, y int
			fmt.Fscan(in, &x, &y)
			if x == y {
				fmt.Fprintln(out, 0)
				continue
			}
			ax, bx := colorsA[x], colorsB[x]
			ay, by := colorsA[y], colorsB[y]
			// If they share a color, direct move is always optimal.
			if ax == ay || ax == by || bx == ay || bx == by {
				fmt.Fprintln(out, abs(x-y))
				continue
			}

			best := inf
			// Four combinations of colors between cities.
			comb := [][2]int{{ax, ay}, {ax, by}, {bx, ay}, {bx, by}}
			for _, c := range comb {
				id := pairID(c[0], c[1])
				if id == -1 {
					continue
				}
				cand := minCostThrough(x, y, lists[id])
				if cand < best {
					best = cand
				}
			}
			if best == inf {
				fmt.Fprintln(out, -1)
			} else {
				fmt.Fprintln(out, best)
			}
		}
	}
}
