package main

import (
	"bufio"
	"fmt"
	"os"
)

type lctPath struct {
	ch  [2][]int
	fa  []int
	val []int64
	sum []int64
}

func newLCTPath(n int) *lctPath {
	l := &lctPath{
		fa:  make([]int, n),
		val: make([]int64, n),
		sum: make([]int64, n),
	}
	l.ch[0] = make([]int, n)
	l.ch[1] = make([]int, n)
	return l
}

func (l *lctPath) isRoot(x int) bool {
	f := l.fa[x]
	return f == 0 || (l.ch[0][f] != x && l.ch[1][f] != x)
}

func (l *lctPath) pushUp(x int) {
	if x == 0 {
		return
	}
	l.sum[x] = l.val[x] + l.sum[l.ch[0][x]] + l.sum[l.ch[1][x]]
}

func (l *lctPath) rotate(x int) {
	y := l.fa[x]
	z := l.fa[y]
	dir := 0
	if l.ch[1][y] == x {
		dir = 1
	}
	b := l.ch[dir^1][x]
	if !l.isRoot(y) {
		if l.ch[0][z] == y {
			l.ch[0][z] = x
		} else if l.ch[1][z] == y {
			l.ch[1][z] = x
		}
	}
	l.fa[x] = z
	l.ch[dir^1][x] = y
	l.fa[y] = x
	l.ch[dir][y] = b
	if b != 0 {
		l.fa[b] = y
	}
	l.pushUp(y)
	l.pushUp(x)
}

func (l *lctPath) splay(x int) {
	for !l.isRoot(x) {
		y := l.fa[x]
		z := l.fa[y]
		if !l.isRoot(y) {
			if (l.ch[0][y] == x) == (l.ch[0][z] == y) {
				l.rotate(y)
			} else {
				l.rotate(x)
			}
		}
		l.rotate(x)
	}
}

func (l *lctPath) access(x int) {
	last := 0
	for v := x; v != 0; v = l.fa[v] {
		l.splay(v)
		l.ch[1][v] = last
		l.pushUp(v)
		last = v
	}
	l.splay(x)
}

func (l *lctPath) addValue(x int, delta int64) {
	l.access(x)
	l.val[x] += delta
	l.pushUp(x)
}

func (l *lctPath) queryPathSum(x int) int64 {
	l.access(x)
	return l.sum[x]
}

type lctSub struct {
	ch   [2][]int
	fa   []int
	val  []int64
	sum  []int64
	vsum []int64
}

func newLCTSub(n int) *lctSub {
	l := &lctSub{
		fa:   make([]int, n),
		val:  make([]int64, n),
		sum:  make([]int64, n),
		vsum: make([]int64, n),
	}
	l.ch[0] = make([]int, n)
	l.ch[1] = make([]int, n)
	return l
}

func (l *lctSub) isRoot(x int) bool {
	f := l.fa[x]
	return f == 0 || (l.ch[0][f] != x && l.ch[1][f] != x)
}

func (l *lctSub) pushUp(x int) {
	if x == 0 {
		return
	}
	l.sum[x] = l.val[x] + l.sum[l.ch[0][x]] + l.sum[l.ch[1][x]] + l.vsum[x]
}

func (l *lctSub) rotate(x int) {
	y := l.fa[x]
	z := l.fa[y]
	dir := 0
	if l.ch[1][y] == x {
		dir = 1
	}
	b := l.ch[dir^1][x]
	if !l.isRoot(y) {
		if l.ch[0][z] == y {
			l.ch[0][z] = x
		} else if l.ch[1][z] == y {
			l.ch[1][z] = x
		}
	}
	l.fa[x] = z
	l.ch[dir^1][x] = y
	l.fa[y] = x
	l.ch[dir][y] = b
	if b != 0 {
		l.fa[b] = y
	}
	l.pushUp(y)
	l.pushUp(x)
}

func (l *lctSub) splay(x int) {
	for !l.isRoot(x) {
		y := l.fa[x]
		z := l.fa[y]
		if !l.isRoot(y) {
			if (l.ch[0][y] == x) == (l.ch[0][z] == y) {
				l.rotate(y)
			} else {
				l.rotate(x)
			}
		}
		l.rotate(x)
	}
}

func (l *lctSub) access(x int) {
	last := 0
	for v := x; v != 0; v = l.fa[v] {
		l.splay(v)
		l.vsum[v] += l.sum[l.ch[1][v]]
		l.ch[1][v] = last
		l.vsum[v] -= l.sum[last]
		l.pushUp(v)
		last = v
	}
	l.splay(x)
}

func (l *lctSub) addValue(x int, delta int64) {
	l.access(x)
	l.val[x] += delta
	l.pushUp(x)
}

func (l *lctSub) querySubtree(x int) int64 {
	l.access(x)
	return l.val[x] + l.vsum[x]
}

func (l *lctSub) link(parent, child int) {
	l.fa[child] = parent
	l.vsum[parent] += l.sum[child]
	l.pushUp(parent)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	maxNodes := n + 5
	children := make([][26]int, maxNodes)
	fail := make([]int, maxNodes)

	lctW := newLCTPath(maxNodes)
	lctF := newLCTSub(maxNodes)

	state1 := make([]int, n+2)
	state2 := make([]int, n+2)
	verts1, verts2 := 1, 1
	root := 1
	state1[1] = root
	state2[1] = root
	totalStates := root

	// root pattern (empty string) exists in tree1 and tree2 initially
	lctW.addValue(root, 1)
	lctF.addValue(root, 1)
	ans := lctW.queryPathSum(root)

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	for i := 0; i < n; i++ {
		var t, v int
		var c string
		fmt.Fscan(in, &t, &v, &c)
		ch := int(c[0] - 'a')
		if t == 1 {
			parentState := state1[v]
			child := children[parentState][ch]
			if child == 0 {
				totalStates++
				child = totalStates
				children[parentState][ch] = child
				f := fail[parentState]
				for f != 0 && children[f][ch] == 0 {
					f = fail[f]
				}
				if children[f][ch] != 0 {
					f = children[f][ch]
				} else {
					f = root
				}
				fail[child] = f
				lctW.fa[child] = f
				lctF.link(f, child)
			}
			verts1++
			state1[verts1] = child
			lctW.addValue(child, 1)
			ans += lctF.querySubtree(child)
		} else {
			parentState := state2[v]
			s := parentState
			for s != 0 && children[s][ch] == 0 {
				s = fail[s]
			}
			if children[s][ch] != 0 {
				s = children[s][ch]
			} else {
				s = root
			}
			verts2++
			state2[verts2] = s
			lctF.addValue(s, 1)
			ans += lctW.queryPathSum(s)
		}
		fmt.Fprintln(out, ans)
	}
}
