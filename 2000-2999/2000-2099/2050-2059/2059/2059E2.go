package main

import (
	"bufio"
	"math/rand"
	"os"
)

type fastScanner struct {
	r *bufio.Reader
}

func newFastScanner() *fastScanner {
	return &fastScanner{r: bufio.NewReader(os.Stdin)}
}

func (fs *fastScanner) nextInt() int {
	sign, val := 1, 0
	b, err := fs.r.ReadByte()
	for err == nil && (b <= ' ' || b > '~') {
		b, err = fs.r.ReadByte()
	}
	if err != nil {
		return 0
	}
	if b == '-' {
		sign = -1
		b, err = fs.r.ReadByte()
	}
	for err == nil && b >= '0' && b <= '9' {
		val = val*10 + int(b-'0')
		b, err = fs.r.ReadByte()
	}
	return sign * val
}

type node struct {
	val         int
	prior       uint32
	left, right *node
	size        int
}

var rng = rand.New(rand.NewSource(1))

func newNode(v int) *node {
	return &node{val: v, prior: rng.Uint32(), size: 1}
}

func size(nd *node) int {
	if nd == nil {
		return 0
	}
	return nd.size
}

func pull(nd *node) {
	if nd != nil {
		nd.size = 1 + size(nd.left) + size(nd.right)
	}
}

func split(root *node, leftSize int) (*node, *node) {
	if root == nil {
		return nil, nil
	}
	if leftSize <= size(root.left) {
		l, newLeft := split(root.left, leftSize)
		root.left = newLeft
		pull(root)
		return l, root
	}
	leftSize -= size(root.left) + 1
	newRight, r := split(root.right, leftSize)
	root.right = newRight
	pull(root)
	return root, r
}

func merge(left, right *node) *node {
	if left == nil {
		return right
	}
	if right == nil {
		return left
	}
	if left.prior < right.prior {
		left.right = merge(left.right, right)
		pull(left)
		return left
	}
	right.left = merge(left, right.left)
	pull(right)
	return right
}

func buildTreap(seq []int) *node {
	var root *node
	for _, v := range seq {
		root = merge(root, newNode(v))
	}
	return root
}

func inorderFill(nd *node, dst []int, idx *int) {
	if nd == nil {
		return
	}
	inorderFill(nd.left, dst, idx)
	dst[*idx] = nd.val
	*idx = *idx + 1
	inorderFill(nd.right, dst, idx)
}

func copySegment(root **node, start, length int, dst []int) {
	left, rest := split(*root, start)
	mid, right := split(rest, length)
	idx := 0
	inorderFill(mid, dst, &idx)
	rest = merge(mid, right)
	*root = merge(left, rest)
}

func insertSequence(root *node, position int, seq []int, totalLen int) *node {
	if len(seq) == 0 {
		return root
	}
	left, right := split(root, position)
	seqNode := buildTreap(seq)
	root = merge(left, merge(seqNode, right))
	keep, _ := split(root, totalLen)
	return keep
}

func overlap(pattern, text []int, pi []int) int {
	m := len(pattern)
	if m == 0 {
		return 0
	}
	pi[0] = 0
	for i := 1; i < m; i++ {
		j := pi[i-1]
		for j > 0 && pattern[i] != pattern[j] {
			j = pi[j-1]
		}
		if pattern[i] == pattern[j] {
			j++
		}
		pi[i] = j
	}
	j := 0
	res := 0
	for _, v := range text {
		for j > 0 && v != pattern[j] {
			j = pi[j-1]
		}
		if v == pattern[j] {
			j++
		}
		res = j
		if j == m {
			j = pi[j-1]
		}
	}
	return res
}

type operation struct {
	idx int
	val int
}

func main() {
	fs := newFastScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t := fs.nextInt()
	for ; t > 0; t-- {
		n := fs.nextInt()
		m := fs.nextInt()
		totalLen := n * m
		flat := make([]int, 0, totalLen)
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				flat = append(flat, fs.nextInt())
			}
		}
		target := make([][]int, n)
		for i := 0; i < n; i++ {
			row := make([]int, m)
			for j := 0; j < m; j++ {
				row[j] = fs.nextInt()
			}
			target[i] = row
		}

		var root *node
		for _, v := range flat {
			root = merge(root, newNode(v))
		}

		cur := make([]int, m)
		pi := make([]int, m)
		ops := make([]operation, 0)

		for i := 0; i < n; i++ {
			start := i * m
			copySegment(&root, start, m, cur)
			lps := overlap(cur, target[i], pi)
			k := m - lps
			for idx := k - 1; idx >= 0; idx-- {
				ops = append(ops, operation{i + 1, target[i][idx]})
			}
			if k > 0 {
				root = insertSequence(root, start, target[i][:k], totalLen)
			}
		}

		out.WriteString(intToString(len(ops)))
		out.WriteByte('\n')
		for _, op := range ops {
			out.WriteString(intToString(op.idx))
			out.WriteByte(' ')
			out.WriteString(intToString(op.val))
			out.WriteByte('\n')
		}
	}
}

func intToString(x int) string {
	if x == 0 {
		return "0"
	}
	var buf [32]byte
	idx := len(buf)
	neg := x < 0
	if neg {
		x = -x
	}
	for x > 0 {
		idx--
		buf[idx] = byte('0' + x%10)
		x /= 10
	}
	if neg {
		idx--
		buf[idx] = '-'
	}
	return string(buf[idx:])
}
