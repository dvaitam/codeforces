package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, k int
	if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
		return
	}

	rowHas := make([]bool, n+1)
	colHas := make([]bool, m+1)

	immediateWin := false
	for i := 0; i < k; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		if x == 1 || y == 1 {
			immediateWin = true
		}
		if x >= 1 && x <= n {
			rowHas[x] = true
		}
		if y >= 1 && y <= m {
			colHas[y] = true
		}
	}

	if immediateWin {
		fmt.Fprintln(out, "Chaneka")
		return
	}

	rNo := 0
	for i := 1; i <= n; i++ {
		if !rowHas[i] {
			rNo++
		}
	}
	cNo := 0
	for i := 1; i <= m; i++ {
		if !colHas[i] {
			cNo++
		}
	}

	rPile := rNo - 1
	cPile := cNo - 1
	if rPile^cPile != 0 {
		fmt.Fprintln(out, "Chaneka")
	} else {
		fmt.Fprintln(out, "Bhinneka")
	}
}
