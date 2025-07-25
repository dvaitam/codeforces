package main

import (
	"bufio"
	"fmt"
	"os"
)

type Fenwick struct {
	n   int
	bit []int
}

func NewFenwick(n int) *Fenwick {
	return &Fenwick{n: n, bit: make([]int, n+2)}
}

func (f *Fenwick) Add(i, val int) {
	for i <= f.n {
		f.bit[i] += val
		i += i & -i
	}
}

func (f *Fenwick) Sum(i int) int {
	s := 0
	for i > 0 {
		s += f.bit[i]
		i -= i & -i
	}
	return s
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		p := make([]int, n)
		pos := make([]int, n+1)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &p[i])
			pos[p[i]] = i + 1
		}

		// inversion count for values <= k
		fw := NewFenwick(n)
		invA := make([]int, n+1)
		for val := 1; val <= n; val++ {
			idx := pos[val]
			greater := (val - 1) - fw.Sum(idx)
			invA[val] = invA[val-1] + greater
			fw.Add(idx, 1)
		}
		// inversion count for values > k
		fw = NewFenwick(n)
		invB := make([]int, n+1)
		for val := n; val >= 1; val-- {
			idx := pos[val]
			smaller := fw.Sum(idx - 1)
			invB[val-1] = invB[val] + smaller
			fw.Add(idx, 1)
		}

		for k := 0; k <= n; k++ {
			if k > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, invA[k]+invB[k])
		}
		fmt.Fprintln(out)
	}
}
