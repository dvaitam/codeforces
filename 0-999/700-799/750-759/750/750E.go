package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int = 1 << 30

// Node represents DP transitions for a substring.
type Node struct {
	dp [5][5]int
}

func identityNode() Node {
	var n Node
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if i == j {
				n.dp[i][j] = 0
			} else {
				n.dp[i][j] = INF
			}
		}
	}
	return n
}

func makeNode(ch byte) Node {
	var n Node
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			n.dp[i][j] = INF
		}
	}
	// deletion option
	for i := 0; i < 5; i++ {
		n.dp[i][i] = 1
	}
	// keep without using
	for i := 0; i < 5; i++ {
		if !(i >= 3 && ch == '6') {
			if n.dp[i][i] > 0 {
				n.dp[i][i] = 0
			}
		}
	}
	// transitions along pattern "2017"
	pattern := []byte{'2', '0', '1', '7'}
	for i := 0; i < 4; i++ {
		if ch == pattern[i] {
			if n.dp[i][i+1] > 0 {
				n.dp[i][i+1] = 0
			}
		}
	}
	// restart on '2'
	if ch == '2' {
		for i := 0; i < 4; i++ {
			if n.dp[i][1] > 0 {
				n.dp[i][1] = 0
			}
		}
	}
	// for state4 (already have 2017)
	if ch == '6' {
		// must delete -> already set cost 1
	} else {
		if n.dp[4][4] > 0 {
			n.dp[4][4] = 0
		}
	}
	return n
}

func merge(a, b Node) Node {
	var res Node
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			res.dp[i][j] = INF
		}
	}
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if a.dp[i][j] >= INF {
				continue
			}
			for k := 0; k < 5; k++ {
				if b.dp[j][k] >= INF {
					continue
				}
				if v := a.dp[i][j] + b.dp[j][k]; v < res.dp[i][k] {
					res.dp[i][k] = v
				}
			}
		}
	}
	return res
}

type SegTree struct {
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
	id := identityNode()
	for i := 0; i < size; i++ {
		if i < n {
			tree[size+i] = makeNode(s[i])
		} else {
			tree[size+i] = id
		}
	}
	for i := size - 1; i >= 1; i-- {
		tree[i] = merge(tree[2*i], tree[2*i+1])
	}
	return &SegTree{size, tree}
}

func (st *SegTree) query(l, r int) Node {
	l += st.size
	r += st.size
	left := identityNode()
	right := identityNode()
	for l <= r {
		if l&1 == 1 {
			left = merge(left, st.tree[l])
			l++
		}
		if r&1 == 0 {
			right = merge(st.tree[r], right)
			r--
		}
		l >>= 1
		r >>= 1
	}
	return merge(left, right)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	var s string
	fmt.Fscan(in, &s)

	st := NewSegTree(s)
	for i := 0; i < q; i++ {
		var l, r int
		fmt.Fscan(in, &l, &r)
		node := st.query(l-1, r-1)
		ans := node.dp[0][4]
		if ans >= INF {
			fmt.Fprintln(out, -1)
		} else {
			fmt.Fprintln(out, ans)
		}
	}
}
