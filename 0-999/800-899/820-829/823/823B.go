package main

import (
	"bufio"
	"fmt"
	"os"
)

type Fenwick struct {
	bit []int
}

func NewFenwick(n int) *Fenwick {
	return &Fenwick{make([]int, n+1)}
}

func (f *Fenwick) Add(i, delta int) {
	for i < len(f.bit) {
		f.bit[i] += delta
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

func (f *Fenwick) Range(l, r int) int {
	if r < l {
		return 0
	}
	return f.Sum(r) - f.Sum(l-1)
}

func idx(c byte) int {
	switch c {
	case 'A':
		return 0
	case 'T':
		return 1
	case 'G':
		return 2
	default: // 'C'
		return 3
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	if _, err := fmt.Fscan(in, &s); err != nil {
		return
	}
	n := len(s)
	data := []byte(s)

	fen := make([][][]*Fenwick, 11)
	for p := 1; p <= 10; p++ {
		fen[p] = make([][]*Fenwick, 4)
		for c := 0; c < 4; c++ {
			fen[p][c] = make([]*Fenwick, p)
			for o := 0; o < p; o++ {
				fen[p][c][o] = NewFenwick(n)
			}
		}
	}

	for i := 1; i <= n; i++ {
		ch := data[i-1]
		c := idx(ch)
		for p := 1; p <= 10; p++ {
			o := (i - 1) % p
			fen[p][c][o].Add(i, 1)
		}
	}

	var q int
	fmt.Fscan(in, &q)
	for ; q > 0; q-- {
		var t int
		fmt.Fscan(in, &t)
		if t == 1 {
			var x int
			var cs string
			fmt.Fscan(in, &x, &cs)
			nc := cs[0]
			if data[x-1] == nc {
				continue
			}
			oc := data[x-1]
			data[x-1] = nc
			oi := idx(oc)
			ni := idx(nc)
			for p := 1; p <= 10; p++ {
				o := (x - 1) % p
				fen[p][oi][o].Add(x, -1)
				fen[p][ni][o].Add(x, 1)
			}
		} else {
			var l, r int
			var e string
			fmt.Fscan(in, &l, &r, &e)
			m := len(e)
			ans := 0
			for i := 0; i < m; i++ {
				c := idx(e[i])
				o := (l + i - 1) % m
				ans += fen[m][c][o].Range(l, r)
			}
			fmt.Fprintln(out, ans)
		}
	}
}
