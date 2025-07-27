package main

import (
	"bufio"
	"fmt"
	"os"
)

func cost(weights []int64) int64 {
	n := len(weights)
	if n == 0 {
		return 0
	}
	preW := make([]int64, n)
	preP := make([]int64, n)
	var sum int64
	for i := 0; i < n; i++ {
		sum += weights[i]
		preW[i] = sum
		if i == 0 {
			preP[i] = int64(i) * weights[i]
		} else {
			preP[i] = preP[i-1] + int64(i)*weights[i]
		}
	}
	total := preW[n-1]
	target := (total + 1) / 2
	r := 0
	for r < n && preW[r] < target {
		r++
	}
	leftW := preW[r]
	leftP := preP[r]
	rightW := total - leftW
	rightP := preP[n-1] - leftP
	return int64(r)*leftW - leftP + rightP - int64(r)*rightW
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, q int
	if _, err := fmt.Fscan(in, &n, &m, &q); err != nil {
		return
	}
	grid := make([][]int64, n)
	row := make([]int64, n)
	col := make([]int64, m)
	for i := 0; i < n; i++ {
		grid[i] = make([]int64, m)
		for j := 0; j < m; j++ {
			var val int64
			fmt.Fscan(in, &val)
			grid[i][j] = val
			row[i] += val
			col[j] += val
		}
	}

	answers := make([]int64, 0, q+1)
	answers = append(answers, cost(row)+cost(col))
	for k := 0; k < q; k++ {
		var x, y int
		var z int64
		fmt.Fscan(in, &x, &y, &z)
		x--
		y--
		diff := z - grid[x][y]
		if diff != 0 {
			grid[x][y] = z
			row[x] += diff
			col[y] += diff
		}
		answers = append(answers, cost(row)+cost(col))
	}

	for i, v := range answers {
		if i > 0 {
			out.WriteByte(' ')
		}
		fmt.Fprint(out, v)
	}
	out.WriteByte('\n')
}
