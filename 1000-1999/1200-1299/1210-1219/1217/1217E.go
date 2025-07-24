package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	LOGN = 9
	INF  = int(2e9)
)

type Node struct {
	best int
	mn   [LOGN]int
}

var (
	a    []int
	tree []Node
)

func newNode() Node {
	n := Node{best: INF}
	for i := 0; i < LOGN; i++ {
		n.mn[i] = INF
	}
	return n
}

func addValue(n *Node, val int) {
	x := val
	for i := 0; i < LOGN; i++ {
		if x%10 != 0 && val < n.mn[i] {
			n.mn[i] = val
		}
		x /= 10
	}
}

func merge(a, b Node) Node {
	c := newNode()
	if a.best < c.best {
		c.best = a.best
	}
	if b.best < c.best {
		c.best = b.best
	}
	for i := 0; i < LOGN; i++ {
		if a.mn[i] < c.mn[i] {
			c.mn[i] = a.mn[i]
		}
		if b.mn[i] < c.mn[i] {
			c.mn[i] = b.mn[i]
		}
		if a.mn[i] < INF && b.mn[i] < INF {
			s := a.mn[i] + b.mn[i]
			if s < c.best {
				c.best = s
			}
		}
	}
	return c
}

func build(v, l, r int) {
	if l == r-1 {
		tree[v] = newNode()
		addValue(&tree[v], a[l])
		return
	}
	m := (l + r) / 2
	build(v*2, l, m)
	build(v*2+1, m, r)
	tree[v] = merge(tree[v*2], tree[v*2+1])
}

func update(v, l, r, pos, val int) {
	if l == r-1 {
		tree[v] = newNode()
		addValue(&tree[v], val)
		return
	}
	m := (l + r) / 2
	if pos < m {
		update(v*2, l, m, pos, val)
	} else {
		update(v*2+1, m, r, pos, val)
	}
	tree[v] = merge(tree[v*2], tree[v*2+1])
}

func query(v, l, r, L, R int) Node {
	if L == l && R == r {
		return tree[v]
	}
	m := (l + r) / 2
	if R <= m {
		return query(v*2, l, m, L, R)
	}
	if L >= m {
		return query(v*2+1, m, r, L, R)
	}
	left := query(v*2, l, m, L, m)
	right := query(v*2+1, m, r, m, R)
	return merge(left, right)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, q int
	fmt.Fscan(reader, &n, &q)
	a = make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	tree = make([]Node, 4*n)
	build(1, 0, n)

	for ; q > 0; q-- {
		var t, x, y int
		fmt.Fscan(reader, &t, &x, &y)
		x--
		if t == 1 {
			update(1, 0, n, x, y)
		} else {
			res := query(1, 0, n, x, y)
			if res.best >= INF {
				fmt.Fprintln(writer, -1)
			} else {
				fmt.Fprintln(writer, res.best)
			}
		}
	}
}
