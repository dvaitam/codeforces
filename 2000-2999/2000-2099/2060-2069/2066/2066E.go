package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	maxValue = 1_000_000
)

var (
	segSize int
	cnt     []int
	sum     []int64
	best    []int64
	bestEx  []int64
	negInf  int64 = -1 << 60

	freq []int

	treapRoot *treapNode
	seed      uint32 = 123456789
)

type treapNode struct {
	key         int
	priority    uint32
	left, right *treapNode
}

func nextPriority() uint32 {
	seed = seed*1664525 + 1013904223
	return seed
}

func rotateRight(p *treapNode) *treapNode {
	l := p.left
	p.left = l.right
	l.right = p
	return l
}

func rotateLeft(p *treapNode) *treapNode {
	r := p.right
	p.right = r.left
	r.left = p
	return r
}

func treapInsert(root *treapNode, key int) *treapNode {
	if root == nil {
		return &treapNode{key: key, priority: nextPriority()}
	}
	if key < root.key {
		root.left = treapInsert(root.left, key)
		if root.left.priority < root.priority {
			root = rotateRight(root)
		}
	} else if key > root.key {
		root.right = treapInsert(root.right, key)
		if root.right.priority < root.priority {
			root = rotateLeft(root)
		}
	}
	return root
}

func treapErase(root *treapNode, key int) *treapNode {
	if root == nil {
		return nil
	}
	if key < root.key {
		root.left = treapErase(root.left, key)
	} else if key > root.key {
		root.right = treapErase(root.right, key)
	} else {
		if root.left == nil {
			return root.right
		}
		if root.right == nil {
			return root.left
		}
		if root.left.priority < root.right.priority {
			root = rotateRight(root)
			root.right = treapErase(root.right, key)
		} else {
			root = rotateLeft(root)
			root.left = treapErase(root.left, key)
		}
	}
	return root
}

func treapMax(root *treapNode) int {
	if root == nil {
		return -1
	}
	for root.right != nil {
		root = root.right
	}
	return root.key
}

func initSegTree() {
	segSize = 1
	for segSize <= maxValue+1 {
		segSize <<= 1
	}
	total := segSize * 2
	cnt = make([]int, total)
	sum = make([]int64, total)
	best = make([]int64, total)
	bestEx = make([]int64, total)
	for i := range best {
		best[i] = negInf
		bestEx[i] = negInf
	}
}

func recalc(idx int) {
	left := idx << 1
	right := left | 1
	cnt[idx] = cnt[left] + cnt[right]
	sum[idx] = sum[left] + sum[right]

	b := best[left]
	br := best[right] - sum[left]
	if br > b {
		b = br
	}
	best[idx] = b

	if cnt[right] == 0 {
		bestEx[idx] = bestEx[left]
	} else {
		val := bestEx[right] - sum[left]
		if best[left] > val {
			val = best[left]
		}
		bestEx[idx] = val
	}
}

func update(pos, delta int) {
	idx := segSize + pos
	cnt[idx] += delta
	c := cnt[idx]
	sum[idx] = int64(c) * int64(pos)
	if c == 0 {
		best[idx] = negInf
		bestEx[idx] = negInf
	} else {
		best[idx] = int64(pos)
		if c >= 2 {
			bestEx[idx] = int64(pos)
		} else {
			bestEx[idx] = negInf
		}
	}
	idx >>= 1
	for idx > 0 {
		recalc(idx)
		idx >>= 1
	}
}

func evaluate() bool {
	total := cnt[1]
	if total <= 1 {
		return true
	}
	maxDup := treapMax(treapRoot)
	if maxDup == -1 {
		return false
	}
	update(maxDup, -2)
	defer update(maxDup, 2)

	totalAfter := cnt[1]
	if totalAfter <= 1 {
		return true
	}
	required := bestEx[1]
	if required == negInf {
		return true
	}
	return int64(2*maxDup) >= required
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	initSegTree()
	freq = make([]int, maxValue+2)

	var n, q int
	fmt.Fscan(in, &n, &q)
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(in, &x)
		freq[x]++
		if freq[x] == 2 {
			treapRoot = treapInsert(treapRoot, x)
		}
		update(x, 1)
	}

	if evaluate() {
		fmt.Fprintln(out, "Yes")
	} else {
		fmt.Fprintln(out, "No")
	}

	for i := 0; i < q; i++ {
		var op string
		var x int
		fmt.Fscan(in, &op, &x)
		if op == "+" {
			freq[x]++
			if freq[x] == 2 {
				treapRoot = treapInsert(treapRoot, x)
			}
			update(x, 1)
		} else {
			if freq[x] == 2 {
				treapRoot = treapErase(treapRoot, x)
			}
			freq[x]--
			update(x, -1)
		}
		if evaluate() {
			fmt.Fprintln(out, "Yes")
		} else {
			fmt.Fprintln(out, "No")
		}
	}
}
