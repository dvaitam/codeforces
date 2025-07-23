package main

import (
	"bufio"
	"fmt"
	"os"
)

type pair struct{ a, b int }

func cost(piece []string, n int) (int, int) {
	cost1 := 0 // top-left should be 1
	cost0 := 0 // top-left should be 0
	for i := 0; i < n; i++ {
		row := piece[i]
		for j := 0; j < n; j++ {
			val := int(row[j] - '0')
			if (i+j)%2 == 0 {
				if val != 1 {
					cost1++
				}
				if val != 0 {
					cost0++
				}
			} else {
				if val != 0 {
					cost1++
				}
				if val != 1 {
					cost0++
				}
			}
		}
	}
	return cost1, cost0
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	costs := make([]pair, 4)
	for k := 0; k < 4; k++ {
		piece := make([]string, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &piece[i])
		}
		a, b := cost(piece, n)
		costs[k] = pair{a, b}
	}

	best := 1<<31 - 1
	for i := 0; i < 4; i++ {
		for j := i + 1; j < 4; j++ {
			total := costs[i].a + costs[j].a
			for k := 0; k < 4; k++ {
				if k != i && k != j {
					total += costs[k].b
				}
			}
			if total < best {
				best = total
			}
		}
	}
	fmt.Fprintln(writer, best)
}
