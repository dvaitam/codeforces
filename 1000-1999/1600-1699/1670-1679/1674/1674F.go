package main

import (
	"bufio"
	"fmt"
	"os"
)

type Fenwick struct{ tree []int }

func NewFenwick(n int) *Fenwick {
	return &Fenwick{make([]int, n+2)}
}

func (f *Fenwick) Add(i, delta int) {
	for i < len(f.tree) {
		f.tree[i] += delta
		i += i & -i
	}
}

func (f *Fenwick) Sum(i int) int {
	s := 0
	for i > 0 {
		s += f.tree[i]
		i -= i & -i
	}
	return s
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, q int
	fmt.Fscan(in, &n, &m, &q)

	totalCells := n * m
	fw := NewFenwick(totalCells)
	state := make([]byte, totalCells+1) // 1-based
	totalStars := 0
	for i := 1; i <= n; i++ {
		var row string
		fmt.Fscan(in, &row)
		for j := 1; j <= m; j++ {
			if row[j-1] == '*' {
				pos := (j-1)*n + i
				state[pos] = 1
				fw.Add(pos, 1)
				totalStars++
			}
		}
	}

	for ; q > 0; q-- {
		var x, y int
		fmt.Fscan(in, &x, &y)
		pos := (y-1)*n + x
		if state[pos] == 1 {
			state[pos] = 0
			fw.Add(pos, -1)
			totalStars--
		} else {
			state[pos] = 1
			fw.Add(pos, 1)
			totalStars++
		}
		good := fw.Sum(totalStars)
		moves := totalStars - good
		fmt.Fprintln(out, moves)
	}
}
