package main

import (
	"bufio"
	"fmt"
	"os"
)

type op struct {
	t int
	r int
	c int
	x int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m, q int
	if _, err := fmt.Fscan(reader, &n, &m, &q); err != nil {
		return
	}

	ops := make([]op, q)
	for i := 0; i < q; i++ {
		var t int
		fmt.Fscan(reader, &t)
		if t == 1 {
			fmt.Fscan(reader, &ops[i].r)
		} else if t == 2 {
			fmt.Fscan(reader, &ops[i].c)
		} else {
			fmt.Fscan(reader, &ops[i].r, &ops[i].c, &ops[i].x)
		}
		ops[i].t = t
	}

	// initialize matrix with zeros
	mat := make([][]int, n)
	for i := 0; i < n; i++ {
		mat[i] = make([]int, m)
	}

	// process operations in reverse
	for i := q - 1; i >= 0; i-- {
		op := ops[i]
		switch op.t {
		case 1:
			r := op.r - 1
			tmp := mat[r][m-1]
			for j := m - 1; j > 0; j-- {
				mat[r][j] = mat[r][j-1]
			}
			mat[r][0] = tmp
		case 2:
			c := op.c - 1
			tmp := mat[n-1][c]
			for j := n - 1; j > 0; j-- {
				mat[j][c] = mat[j-1][c]
			}
			mat[0][c] = tmp
		case 3:
			mat[op.r-1][op.c-1] = op.x
		}
	}

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, mat[i][j])
		}
		writer.WriteByte('\n')
	}
}
