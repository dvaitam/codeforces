package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type pair struct {
	pos int64
	idx int
}

func floorDiv(a, b int64) int64 {
	if a >= 0 {
		return a / b
	}
	return -((-a + b - 1) / b)
}

func mod(a, b int64) int64 {
	r := a % b
	if r < 0 {
		r += b
	}
	return r
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	var m, t int64
	if _, err := fmt.Fscan(reader, &n, &m, &t); err != nil {
		return
	}

	pos0 := make([]int64, n)
	dir := make([]int, n)
	initOrder := make([]pair, n)
	for i := 0; i < n; i++ {
		var s int64
		var d string
		fmt.Fscan(reader, &s, &d)
		pos0[i] = s - 1
		if d == "R" {
			dir[i] = 1
		} else {
			dir[i] = -1
		}
		initOrder[i] = pair{pos: pos0[i], idx: i}
	}

	sort.Slice(initOrder, func(i, j int) bool { return initOrder[i].pos < initOrder[j].pos })

	final := make([]int64, n)
	var shiftSum int64
	for i := 0; i < n; i++ {
		move := pos0[i] + int64(dir[i])*t
		final[i] = mod(move, m)
		shiftSum += floorDiv(move, m)
	}

	sort.Slice(final, func(i, j int) bool { return final[i] < final[j] })

	shift := mod(shiftSum, int64(n))

	res := make([]int64, n)
	for k := 0; k < n; k++ {
		idx := initOrder[k].idx
		res[idx] = final[(int64(k)+shift)%int64(n)] + 1
	}

	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, res[i])
	}
}
