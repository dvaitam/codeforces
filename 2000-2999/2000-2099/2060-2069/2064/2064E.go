package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

const mod int64 = 998244353

var rng = rand.New(rand.NewSource(123456789))

type node struct {
	color     int
	prior     int
	left      *node
	right     *node
	size      int
	prefColor int
	prefLen   int
	suffColor int
	suffLen   int
}

func newNode(color int) *node {
	n := &node{color: color, prior: rng.Int(), size: 1}
	n.prefColor = color
	n.prefLen = 1
	n.suffColor = color
	n.suffLen = 1
	return n
}

func sz(n *node) int {
	if n == nil {
		return 0
	}
	return n.size
}

func recalc(n *node) {
	if n == nil {
		return
	}
	n.size = 1 + sz(n.left) + sz(n.right)

	n.prefColor = n.color
	n.prefLen = 1
	if n.left != nil {
		n.prefColor = n.left.prefColor
		n.prefLen = n.left.prefLen
		if n.left.prefLen == n.left.size && n.left.suffColor == n.color {
			n.prefLen = n.left.size + 1
			if n.right != nil && n.right.prefColor == n.color {
				n.prefLen += n.right.prefLen
			}
		}
	} else if n.right != nil && n.right.prefColor == n.color {
		n.prefLen += n.right.prefLen
	}

	n.suffColor = n.color
	n.suffLen = 1
	if n.right != nil {
		n.suffColor = n.right.suffColor
		n.suffLen = n.right.suffLen
		if n.right.suffLen == n.right.size && n.right.prefColor == n.color {
			n.suffLen = n.right.size + 1
			if n.left != nil && n.left.suffColor == n.color {
				n.suffLen += n.left.suffLen
			}
		}
	} else if n.left != nil && n.left.suffColor == n.color {
		n.suffLen += n.left.suffLen
	}
}

func split(root *node, k int) (*node, *node) {
	if root == nil {
		return nil, nil
	}
	if sz(root.left) >= k {
		left, right := split(root.left, k)
		root.left = right
		recalc(root)
		return left, root
	}
	left, right := split(root.right, k-sz(root.left)-1)
	root.right = left
	recalc(root)
	return root, right
}

func merge(a, b *node) *node {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if a.prior < b.prior {
		a.right = merge(a.right, b)
		recalc(a)
		return a
	}
	b.left = merge(a, b.left)
	recalc(b)
	return b
}

type fenwick struct {
	n   int
	bit []int
}

func newFenwick(n int) *fenwick {
	return &fenwick{n: n, bit: make([]int, n+2)}
}

func (f *fenwick) add(idx, delta int) {
	for idx <= f.n {
		f.bit[idx] += delta
		idx += idx & -idx
	}
}

func (f *fenwick) sum(idx int) int {
	res := 0
	for idx > 0 {
		res += f.bit[idx]
		idx -= idx & -idx
	}
	return res
}

func solveCase(n int, p, c []int) int64 {
	pos := make([]int, n+1)
	for i, v := range p {
		pos[v] = i + 1
	}
	bit := newFenwick(n)
	var root *node
	ans := int64(1)

	for val := n; val >= 1; val-- {
		idx := pos[val]
		leftCount := bit.sum(idx - 1)
		left, right := split(root, leftCount)
		color := c[idx-1]
		block := 1
		if left != nil && left.suffColor == color {
			block += left.suffLen
		}
		if right != nil && right.prefColor == color {
			block += right.prefLen
		}
		ans = (ans * int64(block)) % mod
		root = merge(merge(left, newNode(color)), right)
		bit.add(idx, 1)
	}
	return ans
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
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &p[i])
		}
		c := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &c[i])
		}
		fmt.Fprintln(out, solveCase(n, p, c))
	}
}
