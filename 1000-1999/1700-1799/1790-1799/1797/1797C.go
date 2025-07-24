package main

import (
	"bufio"
	"fmt"
	"os"
)

// Interactive solution for problem C of contest 1797.
// The task is to locate a hidden king on an n x m board using
// at most three distance queries. We first read n and m (if
// provided). Then we follow the standard strategy that works
// for arbitrary board sizes.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	fmt.Println("1 1")
	out.Flush()
	var d1 int
	if _, err := fmt.Fscan(in, &d1); err != nil {
		return
	}

	r1 := 1 + d1
	if r1 > n {
		r1 = n
	}

	fmt.Printf("%d 1\n", r1)
	out.Flush()
	var d2 int
	if _, err := fmt.Fscan(in, &d2); err != nil {
		return
	}

	var row, col int
	if d2 < d1 {
		row = r1
		col = 1 + d2
		if col > m {
			col = m
		}
	} else {
		c1 := 1 + d1
		if c1 > m {
			c1 = m
		}
		fmt.Printf("1 %d\n", c1)
		out.Flush()
		var d3 int
		if _, err := fmt.Fscan(in, &d3); err != nil {
			return
		}
		row = 1 + d3
		if row > n {
			row = n
		}
		col = c1
	}

	fmt.Printf("%d %d\n", row, col)
	out.Flush()
}
