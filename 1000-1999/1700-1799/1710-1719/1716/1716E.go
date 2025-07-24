package main

import (
	"bufio"
	"fmt"
	"os"
)

type Node struct {
	sum, pref, suff, ans int64
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func merge(l, r Node) Node {
	res := Node{}
	res.sum = l.sum + r.sum
	res.pref = max(l.pref, l.sum+r.pref)
	res.suff = max(r.suff, r.sum+l.suff)
	cross := l.suff + r.pref
	res.ans = max(max(l.ans, r.ans), cross)
	return res
}

func build(a []int64, lvl, start int) []Node {
	if lvl == 0 {
		x := a[start]
		n := Node{sum: x}
		if x > 0 {
			n.pref = x
			n.suff = x
			n.ans = x
		}
		return []Node{n}
	}
	half := 1 << (lvl - 1)
	left := build(a, lvl-1, start)
	right := build(a, lvl-1, start+half)
	size := 1 << lvl
	res := make([]Node, size)
	maskLower := half - 1
	for m := 0; m < size; m++ {
		sub := m & maskLower
		if m&half == 0 {
			res[m] = merge(left[sub], right[sub])
		} else {
			res[m] = merge(right[sub], left[sub])
		}
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	size := 1 << n
	a := make([]int64, size)
	for i := 0; i < size; i++ {
		fmt.Fscan(reader, &a[i])
	}
	nodes := build(a, n, 0)
	var q int
	fmt.Fscan(reader, &q)
	mask := 0
	for ; q > 0; q-- {
		var k int
		fmt.Fscan(reader, &k)
		mask ^= 1 << k
		fmt.Fprintln(writer, nodes[mask].ans)
	}
}
