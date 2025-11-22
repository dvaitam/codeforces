package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const INF = 1 << 30

var fib = []int{0, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89}

// computes the minimal possible maximum stack height when cubes with indexes in `bottomMask`
// are used as bottoms of the stacks (each such cube must be present in its own stack),
// and every other cube must be put on one of these stacks (only on a bottom not smaller).
func minMaxHeight(n int, bottomMask int) int {
	bottoms := make([]int, 0)
	others := make([]int, 0)
	for i := 0; i < n; i++ {
		if bottomMask&(1<<i) != 0 {
			bottoms = append(bottoms, i)
		} else {
			others = append(others, i)
		}
	}
	// largest cube must be a bottom, otherwise impossible
	if bottomMask&(1<<(n-1)) == 0 {
		return INF
	}
	// initialise heights with bottom cube sizes
	heights := make([]int, len(bottoms))
	for i, idx := range bottoms {
		heights[i] = fib[idx+1] // fib slice is 1-indexed at position 1
	}
	// sort others descending to improve pruning
	sort.Slice(others, func(i, j int) bool { return fib[others[i]+1] > fib[others[j]+1] })

	best := INF
	var dfs func(pos int)
	dfs = func(pos int) {
		if pos == len(others) {
			// all placed
			maxH := 0
			for _, v := range heights {
				if v > maxH {
					maxH = v
				}
			}
			if maxH < best {
				best = maxH
			}
			return
		}
		// current upper bound for pruning
		curMax := 0
		for _, v := range heights {
			if v > curMax {
				curMax = v
			}
		}
		if curMax >= best {
			return
		}
		sz := fib[others[pos]+1]
		for i, bIdx := range bottoms {
			if fib[bIdx+1] < sz {
				continue
			}
			heights[i] += sz
			maxH := heights[i]
			if maxH < curMax {
				maxH = curMax
			}
			if maxH < best {
				dfs(pos + 1)
			}
			heights[i] -= sz
		}
	}
	dfs(0)
	return best
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}

	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)

		// Precompute best[minSide][height] = minimal sum of bottom sizes needed
		// for some partition where max bottom <= minSide and minimal required height <= height.
		best := make([][]int, 151)
		for i := range best {
			best[i] = make([]int, 151)
			for j := range best[i] {
				best[i][j] = INF
			}
		}

		maxMask := 1 << n
		for mask := 1; mask < maxMask; mask++ {
			if mask&(1<<(n-1)) == 0 { // largest cube must be on bottom
				continue
			}
			maxBottom := 0
			sumBottom := 0
			for i := 0; i < n; i++ {
				if mask&(1<<i) != 0 {
					size := fib[i+1]
					sumBottom += size
					if size > maxBottom {
						maxBottom = size
					}
				}
			}
			if maxBottom > 150 || sumBottom > 150 {
				continue
			}
			minH := minMaxHeight(n, mask)
			if minH > 150 {
				continue
			}
			if sumBottom < best[maxBottom][minH] {
				best[maxBottom][minH] = sumBottom
			}
		}

		// prefix minima to cover larger minSide / height
		for i := 1; i <= 150; i++ {
			for h := 1; h <= 150; h++ {
				v := best[i][h]
				if i > 1 && best[i-1][h] < v {
					v = best[i-1][h]
				}
				if h > 1 && best[i][h-1] < v {
					v = best[i][h-1]
				}
				best[i][h] = v
			}
		}

		for i := 0; i < m; i++ {
			var x, y, z int
			fmt.Fscan(in, &x, &y, &z)
			dims := []int{x, y, z}
			ok := false
			// try all permutations for height choice
			for a := 0; a < 3 && !ok; a++ {
				for b := 0; b < 3 && !ok; b++ {
					if b == a {
						continue
					}
					c := 3 - a - b
					H := dims[a]
					W := dims[b]
					L := dims[c]
					if W < L {
						W, L = L, W
					}
					if H > 150 || L > 150 {
						continue
					}
					if best[L][H] <= W {
						ok = true
					}
				}
			}
			if ok {
				fmt.Fprint(out, "1")
			} else {
				fmt.Fprint(out, "0")
			}
		}
		fmt.Fprintln(out)
	}
}

