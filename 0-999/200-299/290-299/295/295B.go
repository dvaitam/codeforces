package main

import (
	"bufio"
	"os"
	"strconv"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	nextInt := func() int {
		scanner.Scan()
		val, _ := strconv.Atoi(scanner.Text())
		return val
	}

	if !scanner.Scan() {
		return
	}
	n, _ := strconv.Atoi(scanner.Text())

	dist := make([][]int64, n)
	for i := 0; i < n; i++ {
		dist[i] = make([]int64, n)
		for j := 0; j < n; j++ {
			dist[i][j] = int64(nextInt())
		}
	}

	removeOrder := make([]int, n)
	for i := 0; i < n; i++ {
		removeOrder[i] = nextInt() - 1
	}

	results := make([]int64, n)
	activeVertices := make([]int, 0, n)

	for k := n - 1; k >= 0; k-- {
		pivot := removeOrder[k]
		activeVertices = append(activeVertices, pivot)

	
pivotRow := dist[pivot]
		for i := 0; i < n; i++ {
			rowI := dist[i]
			distIP := rowI[pivot]
			for j := 0; j < n; j++ {
				sum := distIP + pivotRow[j]
				if sum < rowI[j] {
					rowI[j] = sum
				}
			}
		}

		var currentSum int64
		for _, i := range activeVertices {
			rowI := dist[i]
			for _, j := range activeVertices {
				currentSum += rowI[j]
			}
		}
		results[k] = currentSum
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	for i := 0; i < n; i++ {
		if i > 0 {
			out.WriteByte(' ')
		}
		out.WriteString(strconv.FormatInt(results[i], 10))
	}
	out.WriteByte('\n')
}