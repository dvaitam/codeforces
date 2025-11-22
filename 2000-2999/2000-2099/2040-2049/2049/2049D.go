package main

import (
	"bufio"
	"fmt"
	"os"
)

const inf int64 = 9e18

// computeRowCosts returns cost[length-1][startColumn] where startColumn is 0-based.
func computeRowCosts(row []int64, m int, k int64) [][]int64 {
	// Duplicate row to handle cyclic segments easily.
	row2 := make([]int64, 2*m)
	copy(row2, row)
	copy(row2[m:], row)

	costs := make([][]int64, m)
	for i := range costs {
		costs[i] = make([]int64, m)
	}

	// For each possible segment length.
	for L := 1; L <= m; L++ {
		val := make([]int64, m)
		// Sliding window to get all cyclic segment sums of length L.
		sum := int64(0)
		for i := 0; i < L; i++ {
			sum += row2[i]
		}
		for start := 0; start < m; start++ {
			val[start] = sum
			if start+L < len(row2) {
				sum += row2[start+L] - row2[start]
			}
		}

		// Build array b[x] = k*x + val[x % m] for x in [0, 2m-2].
		bLen := 2*m - 1
		b := make([]int64, bLen)
		for x := 0; x < bLen; x++ {
			b[x] = k*int64(x) + val[x%m]
		}

		// Sliding window minimum of size m over b.
		deque := make([]int, 0, bLen)
		for x := 0; x < bLen; x++ {
			for len(deque) > 0 && b[deque[len(deque)-1]] >= b[x] {
				deque = deque[:len(deque)-1]
			}
			deque = append(deque, x)
			if x >= m-1 {
				windowStart := x - m + 1
				if deque[0] < windowStart {
					deque = deque[1:]
				}
				minVal := b[deque[0]]
				// cost for entering at column windowStart with length L.
				costs[L-1][windowStart] = minVal - k*int64(windowStart)
			}
		}
	}
	return costs
}

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

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m int
		var k int64
		fmt.Fscan(in, &n, &m, &k)

		grid := make([][]int64, n)
		for i := 0; i < n; i++ {
			row := make([]int64, m)
			for j := 0; j < m; j++ {
				fmt.Fscan(in, &row[j])
			}
			grid[i] = row
		}

		rowCosts := make([][][]int64, n)
		for i := 0; i < n; i++ {
			rowCosts[i] = computeRowCosts(grid[i], m, k)
		}

		dpPrev := make([]int64, m)
		for j := 0; j < m; j++ {
			dpPrev[j] = rowCosts[0][j][0] // length j+1, start column 0
		}

		for i := 1; i < n; i++ {
			dpCurr := make([]int64, m)
			for j := 0; j < m; j++ {
				dpCurr[j] = inf
			}
			for prevCol := 0; prevCol < m; prevCol++ {
				if dpPrev[prevCol] >= inf {
					continue
				}
				for curCol := prevCol; curCol < m; curCol++ {
					length := curCol - prevCol + 1
					rowCost := rowCosts[i][length-1][prevCol]
					if v := dpPrev[prevCol] + rowCost; v < dpCurr[curCol] {
						dpCurr[curCol] = v
					}
				}
			}
			dpPrev = dpCurr
		}

		fmt.Fprintln(out, dpPrev[m-1])
	}
}
