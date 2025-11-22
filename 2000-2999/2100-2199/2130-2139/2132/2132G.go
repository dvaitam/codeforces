package main

import (
	"bufio"
	"fmt"
	"os"
)

// We transpose the grid if needed so that the number of rows (n) does not exceed the number of columns (m).
// For every possible vertical offset dr (ensuring overlap), we compute the set of horizontal offsets dc that
// keep the overlapping cells equal after 180-degree rotation. For a fixed pair of rows (r1, r2) and an offset dc,
// the condition reduces to matching a suffix/prefix of one row with a prefix/suffix of the reversed other row.
// Using rolling hashes we can find the maximum possible overlap lengths L1 (suffix of row1 vs prefix of rev(row2))
// and L2 (prefix of row1 vs suffix of rev(row2)). These yield simple intervals of valid dc values:
//   dc >= 0 valid in [m-L1, m-1]; dc <= 0 valid in [1-m, L2-m].
// Intersecting these intervals across all paired rows for the chosen dr gives all feasible dc values.
// The cost (number of added cells) for offsets (dr, dc) is m*|dr| + n*|dc| + |dr|*|dc| (area difference).
// We search for minimal cost over all dr with non-empty feasible dc.

const base uint64 = 911382323

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		grid := make([][]byte, n)
		for i := 0; i < n; i++ {
			var s string
			fmt.Fscan(in, &s)
			grid[i] = []byte(s)
		}

		// Transpose if rows more than cols to keep n small
		if n > m {
			newGrid := make([][]byte, m)
			for j := 0; j < m; j++ {
				row := make([]byte, n)
				for i := 0; i < n; i++ {
					row[i] = grid[i][j]
				}
				newGrid[j] = row
			}
			grid = newGrid
			n, m = m, n
		}

		// precompute powers
		pow := make([]uint64, m+1)
		pow[0] = 1
		for i := 1; i <= m; i++ {
			pow[i] = pow[i-1] * base
		}

		// prefix hashes for each row and its reverse
		h := make([][]uint64, n)
		hr := make([][]uint64, n)
		for i := 0; i < n; i++ {
			h[i] = make([]uint64, m+1)
			hr[i] = make([]uint64, m+1)
			for j := 0; j < m; j++ {
				h[i][j+1] = h[i][j]*base + uint64(grid[i][j])
				hr[i][j+1] = hr[i][j]*base + uint64(grid[i][m-1-j])
			}
		}

		hashSub := func(pref []uint64, l, r int) uint64 {
			return pref[r] - pref[l]*pow[r-l]
		}

		// function to compute longest L where suffix of row a length L == prefix of reversed row b length L
		l1Func := func(a, b int) int {
			lo, hi := 0, m
			for lo < hi {
				mid := (lo + hi + 1) >> 1
				if hashSub(h[a], m-mid, m) == hashSub(hr[b], 0, mid) {
					lo = mid
				} else {
					hi = mid - 1
				}
			}
			return lo
		}
		// longest L where prefix of row a length L == suffix of reversed row b length L
		l2Func := func(a, b int) int {
			lo, hi := 0, m
			for lo < hi {
				mid := (lo + hi + 1) >> 1
				if hashSub(h[a], 0, mid) == hashSub(hr[b], m-mid, m) {
					lo = mid
				} else {
					hi = mid - 1
				}
			}
			return lo
		}

		best := int64(n*m + 1) // upper bound sentinel

		for dr := -(n - 1); dr <= n-1; dr++ {
			posL, posR := 0, m-1
			negL, negR := 1-m, 0
			valid := true
			xLow := dr
			if xLow < 0 {
				xLow = 0
			}
			xHigh := n + dr - 1
			if xHigh > n-1 {
				xHigh = n - 1
			}
			if xLow > xHigh {
				continue
			}
			for x := xLow; x <= xHigh; x++ {
				y := n + dr - 1 - x
				L1 := l1Func(x, y)
				L2 := l2Func(x, y)
				// positive interval
				if L1 == 0 {
					posL = 1
					posR = 0
				} else {
					if m-L1 > posL {
						posL = m - L1
					}
					if m-1 < posR {
						posR = m - 1
					}
				}
				// negative interval
				if L2 == 0 {
					negL = 1
					negR = 0
				} else {
					if 1-m > negL {
						negL = 1 - m
					}
					if L2-m < negR {
						negR = L2 - m
					}
				}
				if posL > posR && negL > negR {
					valid = false
					break
				}
			}
			if !valid {
				continue
			}
			minDC := int64(1 << 60)
			if posL <= posR {
				if int64(posL) < minDC {
					minDC = int64(posL)
				}
			}
			if negL <= negR {
				cand := int64(-negR) // negR <= 0, smallest abs
				if cand < minDC {
					minDC = cand
				}
			}
			if minDC == int64(1<<60) {
				continue
			}
			ad := absInt(int64(dr))
			cost := int64(m)*ad + int64(n)*minDC + ad*minDC
			if cost < best {
				best = cost
			}
		}

		fmt.Fprintln(out, best)
	}
}

func absInt(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
