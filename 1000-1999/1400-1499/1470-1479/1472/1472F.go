package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		colMask := make(map[int]int)
		for i := 0; i < m; i++ {
			var r, c int
			fmt.Fscan(reader, &r, &c)
			if r == 1 {
				colMask[c] |= 1
			} else {
				colMask[c] |= 2
			}
		}
		cols := make([]int, 0, len(colMask))
		for c := range colMask {
			cols = append(cols, c)
		}
		sort.Ints(cols)

		pendingRow := 0
		pendingCol := 0
		ok := true
		for _, c := range cols {
			mask := colMask[c]
			if mask == 3 {
				if pendingRow != 0 {
					ok = false
					break
				}
				continue
			}
			row := 1
			if mask == 2 {
				row = 2
			}
			if pendingRow == 0 {
				pendingRow = row
				pendingCol = c
			} else {
				diff := c - pendingCol
				same := 0
				if row == pendingRow {
					same = 1
				}
				if diff%2 != same {
					ok = false
					break
				}
				pendingRow = 0
			}
		}
		if pendingRow != 0 {
			ok = false
		}
		if ok {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
