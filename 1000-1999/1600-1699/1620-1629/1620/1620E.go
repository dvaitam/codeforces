package main

import (
	"bufio"
	"fmt"
	"os"
)

type Query struct {
	t int
	x int
	y int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var q int
	if _, err := fmt.Fscan(in, &q); err != nil {
		return
	}
	queries := make([]Query, q)
	maxVal := 500000
	for i := 0; i < q; i++ {
		var t int
		fmt.Fscan(in, &t)
		if t == 1 {
			var x int
			fmt.Fscan(in, &x)
			queries[i] = Query{t: 1, x: x}
			if x > maxVal {
				maxVal = x
			}
		} else {
			var x, y int
			fmt.Fscan(in, &x, &y)
			queries[i] = Query{t: 2, x: x, y: y}
			if x > maxVal {
				maxVal = x
			}
			if y > maxVal {
				maxVal = y
			}
		}
	}
	// mapping for replacement, initialize identity
	m := make([]int, maxVal+1)
	for i := 0; i <= maxVal; i++ {
		m[i] = i
	}

	var res []int
	// process in reverse
	for i := q - 1; i >= 0; i-- {
		qu := queries[i]
		if qu.t == 1 {
			val := m[qu.x]
			res = append(res, val)
		} else {
			m[qu.x] = m[qu.y]
		}
	}

	// output in correct order
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for i := len(res) - 1; i >= 0; i-- {
		if i != len(res)-1 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, res[i])
	}
	if len(res) > 0 {
		fmt.Fprintln(out)
	}
}
