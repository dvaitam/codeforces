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

func (f *Fenwick) Add(i, v int) {
	for i <= f.n {
		f.bit[i] += v
		i += i & -i
	}
}

func (f *Fenwick) Sum(i int) int {
	res := 0
	for i > 0 {
		res += f.bit[i]
		i -= i & -i
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	ps := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &ps[i])
	}
	removed := make([]bool, n+2)
	maxRem := n
	fw := NewFenwick(n)
	res := make([]int, 0, n+1)
	for step := 0; step <= n; step++ {
		if maxRem <= 1 {
			res = append(res, 1)
		} else {
			res = append(res, 1+fw.Sum(maxRem-1))
		}
		if step == n {
			break
		}
		p := ps[step]
		removed[p] = true
		fw.Add(p, 1)
		if p == maxRem {
			for maxRem > 0 && removed[maxRem] {
				maxRem--
			}
		}
	}
	out := bufio.NewWriter(os.Stdout)
	for i, v := range res {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, v)
	}
	out.WriteByte('\n')
	out.Flush()
}
