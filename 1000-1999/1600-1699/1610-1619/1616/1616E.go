package main

import (
	"bufio"
	"fmt"
	"os"
)

type fenwick struct {
	n   int
	bit []int64
}

func newFenwick(n int) *fenwick {
	return &fenwick{n: n, bit: make([]int64, n+2)}
}

func (f *fenwick) add(i int, v int64) {
	for i <= f.n {
		f.bit[i] += v
		i += i & -i
	}
}

func (f *fenwick) sum(i int) int64 {
	res := int64(0)
	for i > 0 {
		res += f.bit[i]
		i -= i & -i
	}
	return res
}

var reader = bufio.NewReader(os.Stdin)
var writer = bufio.NewWriter(os.Stdout)

func solveOne() {
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	var s, t string
	fmt.Fscan(reader, &s)
	fmt.Fscan(reader, &t)

	pos := make([][]int, 26)
	for i := 0; i < n; i++ {
		c := int(s[i] - 'a')
		pos[c] = append(pos[c], i+1)
	}
	idx := make([]int, 26)

	fw := newFenwick(n)
	for i := 1; i <= n; i++ {
		fw.add(i, 1)
	}

	const INF int64 = 1 << 60
	cost := int64(0)
	ans := INF

	for i := 0; i < n; i++ {
		cur := int(t[i] - 'a')
		for ch := 0; ch < cur; ch++ {
			if idx[ch] < len(pos[ch]) {
				p := pos[ch][idx[ch]]
				moves := fw.sum(p) - 1
				if cost+moves < ans {
					ans = cost + moves
				}
			}
		}
		if idx[cur] >= len(pos[cur]) {
			break
		}
		p := pos[cur][idx[cur]]
		cost += fw.sum(p) - 1
		fw.add(p, -1)
		idx[cur]++
	}

	if ans == INF {
		fmt.Fprintln(writer, -1)
	} else {
		fmt.Fprintln(writer, ans)
	}
}

func main() {
	defer writer.Flush()
	var q int
	fmt.Fscan(reader, &q)
	for ; q > 0; q-- {
		solveOne()
	}
}
