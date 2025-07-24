package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"time"
)

// node represents a treap node storing (key, idx)
type node struct {
	key, idx, prio int
	left, right    *node
}

type Treap struct{ root *node }

func less(k1, i1, k2, i2 int) bool {
	if k1 != k2 {
		return k1 < k2
	}
	return i1 < i2
}

func split(t *node, key, idx int) (l, r *node) {
	if t == nil {
		return nil, nil
	}
	if less(t.key, t.idx, key, idx) {
		t.right, r = split(t.right, key, idx)
		l = t
	} else {
		l, t.left = split(t.left, key, idx)
		r = t
	}
	return
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
		return a
	}
	b.left = merge(a, b.left)
	return b
}

func (tr *Treap) Insert(key, idx int) {
	n := &node{key: key, idx: idx, prio: rand.Int()}
	l, r := split(tr.root, key, idx)
	tr.root = merge(merge(l, n), r)
}

func (tr *Treap) delete(t *node, key, idx int) *node {
	if t == nil {
		return nil
	}
	if t.key == key && t.idx == idx {
		return merge(t.left, t.right)
	}
	if less(t.key, t.idx, key, idx) {
		t.right = tr.delete(t.right, key, idx)
	} else {
		t.left = tr.delete(t.left, key, idx)
	}
	return t
}

func (tr *Treap) Delete(key, idx int) { tr.root = tr.delete(tr.root, key, idx) }
func (tr *Treap) Empty() bool         { return tr.root == nil }

func (tr *Treap) Min() (int, int, bool) {
	t := tr.root
	if t == nil {
		return 0, 0, false
	}
	for t.left != nil {
		t = t.left
	}
	return t.key, t.idx, true
}

func (tr *Treap) Max() (int, int, bool) {
	t := tr.root
	if t == nil {
		return 0, 0, false
	}
	for t.right != nil {
		t = t.right
	}
	return t.key, t.idx, true
}

func check(sa, ea, sb, eb []int) bool {
	n := len(sa)
	idxStart := make([]int, n)
	idxEnd := make([]int, n)
	for i := 0; i < n; i++ {
		idxStart[i] = i
		idxEnd[i] = i
	}
	sort.Slice(idxStart, func(i, j int) bool {
		if sa[idxStart[i]] == sa[idxStart[j]] {
			return idxStart[i] < idxStart[j]
		}
		return sa[idxStart[i]] < sa[idxStart[j]]
	})
	sort.Slice(idxEnd, func(i, j int) bool {
		if ea[idxEnd[i]] == ea[idxEnd[j]] {
			return idxEnd[i] < idxEnd[j]
		}
		return ea[idxEnd[i]] < ea[idxEnd[j]]
	})
	tSb := &Treap{}
	tEb := &Treap{}
	e := 0
	for _, id := range idxStart {
		for e < n && ea[idxEnd[e]] < sa[id] {
			rm := idxEnd[e]
			tSb.Delete(sb[rm], rm)
			tEb.Delete(eb[rm], rm)
			e++
		}
		if !tSb.Empty() {
			mx, _, _ := tSb.Max()
			if mx > eb[id] {
				return false
			}
			mn, _, _ := tEb.Min()
			if mn < sb[id] {
				return false
			}
		}
		tSb.Insert(sb[id], id)
		tEb.Insert(eb[id], id)
	}
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	rand.Seed(time.Now().UnixNano())

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	sa := make([]int, n)
	ea := make([]int, n)
	sb := make([]int, n)
	eb := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &sa[i], &ea[i], &sb[i], &eb[i])
	}
	if check(sa, ea, sb, eb) && check(sb, eb, sa, ea) {
		fmt.Fprintln(writer, "YES")
	} else {
		fmt.Fprintln(writer, "NO")
	}
}
