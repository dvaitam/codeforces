package main

import (
	"bufio"
	"fmt"
	"os"
)

func solve(n, rK, cK, rD, cD int64) int64 {
	dx := rK - rD
	dy := cK - cD

	if dx != 0 && dy != 0 {
		var rowTarget int64
		if dx < 0 {
			rowTarget = 0
		} else {
			rowTarget = n
		}
		var colTarget int64
		if dy < 0 {
			colTarget = 0
		} else {
			colTarget = n
		}
		distRow := abs(rowTarget - rD)
		distCol := abs(colTarget - cD)
		if distRow > distCol {
			return distRow
		}
		return distCol
	}

	if dx == 0 {
		if dy > 0 {
			return (n - cK) + dy
		}
		return cD
	}

	if dx > 0 {
		return (n - rK) + dx
	}
	return rD
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, rK, cK, rD, cD int64
		fmt.Fscan(reader, &n, &rK, &cK, &rD, &cD)
		fmt.Fprintln(writer, solve(n, rK, cK, rD, cD))
	}
}
