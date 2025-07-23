package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

const maxV = 100000

// Treap node for ordered set of integers
type node struct {
	key         int
	prio        int
	left, right *node
	sz          int
}

func sz(n *node) int {
	if n == nil {
		return 0
	}
	return n.sz
}
func upd(n *node) {
	if n != nil {
		n.sz = 1 + sz(n.left) + sz(n.right)
	}
}

func split(root *node, key int) (l, r *node) {
	if root == nil {
		return nil, nil
	}
	if key <= root.key {
		l, root.left = split(root.left, key)
		upd(root)
		return l, root
	}
	root.right, r = split(root.right, key)
	upd(root)
	return root, r
}

func merge(a, b *node) *node {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if a.prio < b.prio {
		a.right = merge(a.right, b)
		upd(a)
		return a
	}
	b.left = merge(a, b.left)
	upd(b)
	return b
}

func insert(root *node, key int) *node {
	n := &node{key: key, prio: rand.Int()}
	n.sz = 1
	l, r := split(root, key)
	return merge(merge(l, n), r)
}

func erase(root *node, key int) *node {
	if root == nil {
		return nil
	}
	if key == root.key {
		return merge(root.left, root.right)
	}
	if key < root.key {
		root.left = erase(root.left, key)
	} else {
		root.right = erase(root.right, key)
	}
	upd(root)
	return root
}

func lowerBound(root *node, key int) int {
	ans := -1
	for root != nil {
		if root.key >= key {
			ans = root.key
			root = root.left
		} else {
			root = root.right
		}
	}
	return ans
}

func countLess(root *node, key int) int {
	if root == nil {
		return 0
	}
	if key <= root.key {
		return countLess(root.left, key)
	}
	return sz(root.left) + 1 + countLess(root.right, key)
}

type Treap struct{ root *node }

func (t *Treap) Insert(key int)         { t.root = insert(t.root, key) }
func (t *Treap) Erase(key int)          { t.root = erase(t.root, key) }
func (t *Treap) LowerBound(key int) int { return lowerBound(t.root, key) }
func (t *Treap) CountRange(l, r int) int {
	return countLess(t.root, r+1) - countLess(t.root, l)
}

var (
	arr    []int
	sets   [maxV + 1]*Treap
	parent [maxV + 2]int
)

func find(x int) int {
	if parent[x] != x {
		parent[x] = find(parent[x])
	}
	return parent[x]
}

func union(x, y int) {
	px := find(x)
	py := find(y)
	if px != py {
		parent[px] = py
	}
}

func decRange(l, r, x int) {
	for v := find(x + 1); v <= maxV; v = find(v + 1) {
		t := sets[v]
		if t == nil || t.root == nil {
			union(v, v+1)
			continue
		}
		for {
			pos := t.LowerBound(l)
			if pos == -1 || pos > r {
				break
			}
			t.Erase(pos)
			arr[pos] -= x
			nv := v - x
			if sets[nv] == nil {
				sets[nv] = &Treap{}
			}
			sets[nv].Insert(pos)
		}
		if t.root == nil {
			union(v, v+1)
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	arr = make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &arr[i])
		v := arr[i]
		if sets[v] == nil {
			sets[v] = &Treap{}
		}
		sets[v].Insert(i)
	}
	for i := 0; i <= maxV+1; i++ {
		parent[i] = i
	}
	for ; m > 0; m-- {
		var typ, l, r, x int
		fmt.Fscan(reader, &typ, &l, &r, &x)
		if typ == 1 {
			decRange(l, r, x)
		} else {
			if x <= maxV && sets[x] != nil {
				ans := sets[x].CountRange(l, r)
				fmt.Fprintln(writer, ans)
			} else {
				fmt.Fprintln(writer, 0)
			}
		}
	}
}
