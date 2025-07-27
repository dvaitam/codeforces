package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// Treap implementation for ordered set of integers

type tnode struct {
	key         int
	pr          uint32
	left, right *tnode
}

type Treap struct{ root *tnode }

func rotateRight(p *tnode) *tnode {
	q := p.left
	p.left = q.right
	q.right = p
	return q
}
func rotateLeft(p *tnode) *tnode {
	q := p.right
	p.right = q.left
	q.left = p
	return q
}

func insertNode(p *tnode, key int) *tnode {
	if p == nil {
		return &tnode{key: key, pr: rand.Uint32()}
	}
	if key < p.key {
		p.left = insertNode(p.left, key)
		if p.left.pr < p.pr {
			p = rotateRight(p)
		}
	} else if key > p.key {
		p.right = insertNode(p.right, key)
		if p.right.pr < p.pr {
			p = rotateLeft(p)
		}
	}
	return p
}

func merge(a, b *tnode) *tnode {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if a.pr < b.pr {
		a.right = merge(a.right, b)
		return a
	}
	b.left = merge(a, b.left)
	return b
}

func eraseNode(p *tnode, key int) *tnode {
	if p == nil {
		return nil
	}
	if key < p.key {
		p.left = eraseNode(p.left, key)
	} else if key > p.key {
		p.right = eraseNode(p.right, key)
	} else {
		p = merge(p.left, p.right)
	}
	return p
}

func (t *Treap) Insert(key int) { t.root = insertNode(t.root, key) }
func (t *Treap) Remove(key int) { t.root = eraseNode(t.root, key) }

func (t *Treap) Prev(key int) int {
	cur := t.root
	res := -1
	for cur != nil {
		if key <= cur.key {
			cur = cur.left
		} else {
			res = cur.key
			cur = cur.right
		}
	}
	return res
}

func (t *Treap) Next(key int) int {
	cur := t.root
	res := int(^uint(0) >> 1) // max int
	for cur != nil {
		if key >= cur.key {
			cur = cur.right
		} else {
			res = cur.key
			cur = cur.left
		}
	}
	return res
}

// Problem solution

var (
	n     int
	gS    [][]int
	gK    [][]int
	tin   []int
	tout  []int
	timer int

	set     Treap
	posToID map[int]int
	curSize int
	ans     int
)

func dfsK(u int) {
	timer++
	tin[u] = timer
	for _, v := range gK[u] {
		dfsK(v)
	}
	tout[u] = timer
}

func isAncestor(u, v int) bool { // v inside u in K tree
	return tin[u] <= tin[v] && tout[v] <= tout[u]
}

type info struct {
	key int
	id  int
}

func addNode(u int) info {
	key := tin[u]
	rep := info{key: -1, id: -1}
	succKey := set.Next(key)
	if succKey != int(^uint(0)>>1) {
		succID := posToID[succKey]
		if isAncestor(u, succID) {
			set.Remove(succKey)
			delete(posToID, succKey)
			curSize--
			rep = info{key: succKey, id: succID}
		}
	}
	prevKey := set.Prev(key)
	if prevKey != -1 {
		prevID := posToID[prevKey]
		if isAncestor(prevID, u) {
			set.Remove(prevKey)
			delete(posToID, prevKey)
			curSize--
			rep = info{key: prevKey, id: prevID}
		}
	}
	set.Insert(key)
	posToID[key] = u
	curSize++
	if curSize > ans {
		ans = curSize
	}
	return rep
}

func removeNode(u int, rep info) {
	key := tin[u]
	set.Remove(key)
	delete(posToID, key)
	curSize--
	if rep.key != -1 {
		set.Insert(rep.key)
		posToID[rep.key] = rep.id
		curSize++
	}
}

func dfsS(u int) {
	rep := addNode(u)
	for _, v := range gS[u] {
		dfsS(v)
	}
	removeNode(u, rep)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		fmt.Fscan(in, &n)
		gS = make([][]int, n+1)
		gK = make([][]int, n+1)
		tin = make([]int, n+1)
		tout = make([]int, n+1)
		for i := 2; i <= n; i++ {
			var p int
			fmt.Fscan(in, &p)
			gS[p] = append(gS[p], i)
		}
		for i := 2; i <= n; i++ {
			var p int
			fmt.Fscan(in, &p)
			gK[p] = append(gK[p], i)
		}
		timer = 0
		dfsK(1)
		set = Treap{}
		posToID = make(map[int]int)
		curSize = 0
		ans = 0
		dfsS(1)
		fmt.Fprintln(out, ans)
	}
}
