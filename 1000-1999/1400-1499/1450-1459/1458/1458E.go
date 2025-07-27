package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	// read shortcut positions
	rowMap := make(map[int][]int) // y -> list of x
	colMap := make(map[int][]int) // x -> list of y
	shortSet := make(map[[2]int]struct{})
	for i := 0; i < n; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		rowMap[y] = append(rowMap[y], x)
		colMap[x] = append(colMap[x], y)
		shortSet[[2]int{x, y}] = struct{}{}
	}
	for y := range rowMap {
		xs := rowMap[y]
		sort.Ints(xs)
		rowMap[y] = xs
	}
	for x := range colMap {
		ys := colMap[x]
		sort.Ints(ys)
		colMap[x] = ys
	}

	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		if _, ok := shortSet[[2]int{a, b}]; ok {
			fmt.Fprintln(out, "LOSE")
			continue
		}
		win := false
		if xs, ok := rowMap[b]; ok {
			// check if there exists shortcut (x,b) with x < a
			idx := sort.SearchInts(xs, a)
			if idx > 0 {
				win = true
			}
		}
		if !win {
			if ys, ok := colMap[a]; ok {
				idx := sort.SearchInts(ys, b)
				if idx > 0 {
					win = true
				}
			}
		}
		if !win {
			if a^b != 0 {
				win = true
			}
		}
		if win {
			fmt.Fprintln(out, "WIN")
		} else {
			fmt.Fprintln(out, "LOSE")
		}
	}
}
