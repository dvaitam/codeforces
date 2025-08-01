package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func solveCase(grid [][]byte, n, m int) int32 {
	// Initialize prefix sum array
	pref := make([][]int32, n+1)
	for i := range pref {
		pref[i] = make([]int32, m)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			pref[i+1][j] = pref[i][j] + int32(grid[i][j])
		}
	}

	best := int32(1<<31 - 1) // Equivalent to std::i32::MAX
	vertZero := make([]int32, m)
	sumCol := make([]int32, m)
	ps := make([]int32, m+1)

	for r1 := 0; r1 < n; r1++ {
		for r2 := r1 + 4; r2 < n; r2++ {
			h := int32(r2 - r1 + 1)
			for j := 0; j < m; j++ {
				interiorOnes := pref[r2][j] - pref[r1+1][j]
				vertZero[j] = (h - 2) - interiorOnes
				topBottomMissing := 2 - (int32(grid[r1][j]) + int32(grid[r2][j]))
				sumCol[j] = interiorOnes + topBottomMissing
			}
			ps[0] = 0
			for j := 0; j < m; j++ {
				ps[j+1] = ps[j] + sumCol[j]
			}
			minVal := int32(1<<31 - 1)
			for c2 := 3; c2 < m; c2++ {
				c1 := c2 - 3
				candidate := vertZero[c1] - ps[c1+1]
				if candidate < minVal {
					minVal = candidate
				}
				cost := vertZero[c2] + ps[c2] + minVal
				if cost < best {
					best = cost
				}
			}
		}
	}
	return best
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	var t int
	fmt.Sscanf(scanner.Text(), "%d", &t)

	for ; t > 0; t-- {
		scanner.Scan()
		var n, m int
		fmt.Sscanf(scanner.Text(), "%d %d", &n, &m)
		grid := make([][]byte, n)
		for i := 0; i < n; i++ {
			scanner.Scan()
			grid[i] = []byte(scanner.Text())
			for j := range grid[i] {
				grid[i][j] -= '0' // Convert '0' or '1' to 0 or 1
			}
		}
		ans := solveCase(grid, n, m)
		fmt.Println(ans)
	}
}
