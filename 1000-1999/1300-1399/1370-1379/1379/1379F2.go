package main

import (
	"bufio"
	"fmt"
	"os"
)

// This is a very simplified solution for the problem described in
// problemF2.txt. It only checks if it is possible to place all kings
// using either all white cells with odd coordinates or all white cells
// with even coordinates. This does not cover all valid configurations
// but illustrates the handling of dynamic queries.

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m, q int
	if _, err := fmt.Fscan(reader, &n, &m, &q); err != nil {
		return
	}

	type cell struct{ i, j int }
	blocked := make(map[cell]bool)

	total := n * m
	availA, availB := total, total // A: (odd,odd) cells; B: (even,even) cells

	for ; q > 0; q-- {
		var i, j int
		fmt.Fscan(reader, &i, &j)
		c := cell{i, j}
		odd := i%2 == 1 // true if cell belongs to group A
		if blocked[c] {
			delete(blocked, c)
			if odd {
				availA++
			} else {
				availB++
			}
		} else {
			blocked[c] = true
			if odd {
				availA--
			} else {
				availB--
			}
		}
		if availA == total || availB == total {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
