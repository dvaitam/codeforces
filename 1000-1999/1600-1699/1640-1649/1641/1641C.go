package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemC.txt.
// We maintain intervals that must contain at least one sick
// person. For queries marking people healthy we use a disjoint
// set union to skip processed indices and efficiently update
// the affected intervals. Each interval is represented only by
// the index of its current leftmost possible sick person. When
// this index becomes healthy we move the interval to the next
// candidate. For questions about a person we check if among the
// intervals currently pointing to this position there is one
// whose right boundary lies before the next unchecked position;
// if so, this person must be sick.

type Fenwick struct {
	n int
	t []int
}

func NewFenwick(n int) *Fenwick {
	return &Fenwick{n: n, t: make([]int, n+2)}
}

func (f *Fenwick) Add(i, v int) {
	for i <= f.n {
		f.t[i] += v
		i += i & -i
	}
}

func (f *Fenwick) Sum(i int) int {
	s := 0
	for i > 0 {
		s += f.t[i]
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

const inf = int(1e9)

type Interval struct {
	l, r int
	pos  int
	act  bool
}

var (
	n      int
	status []int // 0 unknown, -1 healthy, 1 sick
	parent []int
	fenw   *Fenwick
	ivs    []*Interval
	byPos  [][]int
	minR   []int
)

func find(x int) int {
	if x > n+1 {
		return n + 1
	}
	if parent[x] != x {
		parent[x] = find(parent[x])
	}
	return parent[x]
}

func setSick(i int) {
	if i < 1 || i > n || status[i] == 1 {
		return
	}
	status[i] = 1
	fenw.Add(i, 1)
	for _, id := range byPos[i] {
		iv := ivs[id]
		if iv.act {
			iv.act = false
		}
	}
	byPos[i] = nil
	minR[i] = inf
	parent[i] = find(i + 1)
}

func processHealthyRange(l, r int) {
	i := find(l)
	var ids []int
	for i <= r {
		if status[i] == 0 {
			ids = append(ids, byPos[i]...)
			byPos[i] = nil
			minR[i] = inf
			status[i] = -1
		}
		nxt := find(i + 1)
		parent[i] = nxt
		i = nxt
	}
	for _, id := range ids {
		iv := ivs[id]
		if !iv.act {
			continue
		}
		if fenw.Range(iv.l, iv.r) > 0 {
			iv.act = false
			continue
		}
		pos := find(iv.pos)
		if pos > iv.r {
			iv.act = false
			continue
		}
		iv.pos = pos
		byPos[pos] = append(byPos[pos], id)
		if iv.r < minR[pos] {
			minR[pos] = iv.r
		}
	}
}

func addInterval(l, r int) {
	if fenw.Range(l, r) > 0 {
		return
	}
	pos := find(l)
	if pos > r {
		return
	}
	id := len(ivs)
	ivs = append(ivs, &Interval{l: l, r: r, pos: pos, act: true})
	byPos[pos] = append(byPos[pos], id)
	if r < minR[pos] {
		minR[pos] = r
	}
}

func queryPerson(j int) string {
	if status[j] == 1 {
		return "YES"
	}
	if status[j] == -1 {
		return "NO"
	}
	r := minR[j]
	nxt := find(j + 1)
	if r < nxt {
		setSick(j)
		return "YES"
	}
	return "N/A"
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var q int
	fmt.Fscan(reader, &n, &q)
	status = make([]int, n+2)
	parent = make([]int, n+2)
	byPos = make([][]int, n+2)
	minR = make([]int, n+2)
	for i := 1; i <= n+1; i++ {
		parent[i] = i
		minR[i] = inf
	}
	fenw = NewFenwick(n)
	for ; q > 0; q-- {
		var t int
		fmt.Fscan(reader, &t)
		if t == 0 {
			var l, r, x int
			fmt.Fscan(reader, &l, &r, &x)
			if x == 0 {
				processHealthyRange(l, r)
			} else {
				addInterval(l, r)
			}
		} else {
			var j int
			fmt.Fscan(reader, &j)
			fmt.Fprintln(writer, queryPerson(j))
		}
	}
}
