package main

import (
	"bufio"
	"os"
	"strconv"
)

func main() {
	// Set up fast I/O
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	// Helper to read the next integer
	next := func() int {
		scanner.Scan()
		val, _ := strconv.Atoi(scanner.Text())
		return val
	}

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	// Read dimensions and number of operations
	n := next()
	m := next()
	k := next()

	// Arrays to store the time (operation index) and color for rows and columns
	// Using 1-based indexing matching the problem statement
	rTime := make([]int, n+1)
	rColor := make([]int, n+1)
	cTime := make([]int, m+1)
	cColor := make([]int, m+1)

	// Process operations
	for t := 1; t <= k; t++ {
		opType := next()
		idx := next()
		color := next()

		if opType == 1 {
			// Row operation
			rTime[idx] = t
			rColor[idx] = color
		} else {
			// Column operation
			cTime[idx] = t
			cColor[idx] = color
		}
	}

	// Construct and print the resulting table
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if j > 1 {
				writer.WriteByte(' ')
			}

			// Determine the color of cell (i, j)
			// The color is determined by the latest operation that affected this row or column
			var ans int
			if rTime[i] > cTime[j] {
				ans = rColor[i]
			} else {
				ans = cColor[j]
			}
			writer.WriteString(strconv.Itoa(ans))
		}
		writer.WriteByte('\n')
	}
}