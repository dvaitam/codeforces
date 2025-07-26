package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int = 1e9

// Node stores DP transitions from automaton states.
type Node struct {
	mat [3][3]int
}

// nextState returns the automaton state after reading character ch
// when currently in state cur (0: none, 1: have 'a', 2: have 'ab').
func nextState(cur int, ch byte) int {
	switch cur {
	case 0:
		if ch == 'a' {
			return 1
		}
		return 0
	case 1:
		if ch == 'a' {
			return 1
		}
		if ch == 'b' {
			return 2
		}
		if ch == 'c' {
			return 1
		}
	case 2:
		if ch == 'c' {
			return 3 // forbidden state
		}
		return 2
	}
	return 3
}

// makeLeaf constructs a leaf node for given character ch.
func makeLeaf(ch byte) Node {
	n := Node{}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			n.mat[i][j] = INF
		}
	}
	letters := []byte{'a', 'b', 'c'}
	for i := 0; i < 3; i++ {
		for _, l := range letters {
			ns := nextState(i, l)
			if ns == 3 {
				continue
			}
			cost := 0
			if ch != l {
				cost = 1
			}
			if cost < n.mat[i][ns] {
				n.mat[i][ns] = cost
			}
		}
	}
	return n
}

// merge combines two nodes.
func merge(a, b Node) Node {
	res := Node{}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			res.mat[i][j] = INF
		}
	}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if a.mat[i][j] >= INF {
				continue
			}
			for k := 0; k < 3; k++ {
				if b.mat[j][k] >= INF {
					continue
				}
				val := a.mat[i][j] + b.mat[j][k]
				if val < res.mat[i][k] {
					res.mat[i][k] = val
				}
			}
		}
	}
	return res
}

// Segment tree implementation

type SegTree struct {
	n    int
	size int
	tree []Node
}

func NewSegTree(s string) *SegTree {
	n := len(s)
	size := 1
	for size < n {
		size <<= 1
	}
	tree := make([]Node, 2*size)
	st := &SegTree{n: n, size: size, tree: tree}
	for i := 0; i < n; i++ {
		st.tree[size+i] = makeLeaf(s[i])
	}
	for i := size - 1; i >= 1; i-- {
		st.tree[i] = merge(st.tree[i<<1], st.tree[i<<1|1])
	}
	return st
}

func (st *SegTree) Update(pos int, ch byte) {
	p := pos + st.size
	st.tree[p] = makeLeaf(ch)
	for p >>= 1; p >= 1; p >>= 1 {
		st.tree[p] = merge(st.tree[p<<1], st.tree[p<<1|1])
	}
}

func (st *SegTree) Answer() int {
	root := st.tree[1]
	ans := INF
	for i := 0; i < 3; i++ {
		if root.mat[0][i] < ans {
			ans = root.mat[0][i]
		}
	}
	return ans
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, q int
	if _, err := fmt.Fscan(reader, &n, &q); err != nil {
		return
	}
	var s string
	fmt.Fscan(reader, &s)

	st := NewSegTree(s)

	for ; q > 0; q-- {
		var pos int
		var c string
		fmt.Fscan(reader, &pos, &c)
		pos--
		st.Update(pos, c[0])
		fmt.Fprintln(writer, st.Answer())
	}
}
