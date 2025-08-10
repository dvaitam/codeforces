package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	scanner.Scan()
	var t int
	fmt.Sscan(scanner.Text(), &t)

	for i := 0; i < t; i++ {
		scanner.Scan()
		var n int
		fmt.Sscan(scanner.Text(), &n)

		maxIdx := 0
		for j := 1; j < n; j++ {
			cmp := query(writer, scanner, maxIdx, maxIdx, j, j)
			if cmp == "<" {
				maxIdx = j
			}
		}

		winners := make([]int, 0, n-1)
		for j := 0; j < n; j++ {
			if j != maxIdx {
				winners = append(winners, j)
			}
		}

		for len(winners) > 1 {
			newWinners := make([]int, 0, len(winners)/2+1)
			for j := 0; j < len(winners); j += 2 {
				if j+1 >= len(winners) {
					newWinners = append(newWinners, winners[j])
					break
				}
				k := winners[j]
				l := winners[j+1]
				cmpScore := query(writer, scanner, maxIdx, k, maxIdx, l)
				var betterIdx int
				if cmpScore == ">" {
					betterIdx = k
				} else if cmpScore == "<" {
					betterIdx = l
				} else {
					cmpP := query(writer, scanner, k, k, l, l)
					if cmpP == "<" {
						betterIdx = k
					} else if cmpP == ">" {
						betterIdx = l
					} else {
						betterIdx = k
					}
				}
				newWinners = append(newWinners, betterIdx)
			}
			winners = newWinners
		}

		bestJ := winners[0]
		fmt.Fprintf(writer, "! %d %d\n", maxIdx, bestJ)
		writer.Flush()
	}
}

func query(writer *bufio.Writer, scanner *bufio.Scanner, a, b, c, d int) string {
	fmt.Fprintf(writer, "? %d %d %d %d\n", a, b, c, d)
	writer.Flush()
	scanner.Scan()
	return scanner.Text()
}
