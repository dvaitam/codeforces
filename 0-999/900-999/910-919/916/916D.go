package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// persistent trie node
type trieNode struct {
	child [26]int
	val   int
}

var trie []trieNode

func cloneTrie(idx int) int {
	trie = append(trie, trieNode{})
	newIdx := len(trie) - 1
	if idx != 0 {
		trie[newIdx] = trie[idx]
	}
	return newIdx
}

func trieQuery(root int, s string) int {
	p := root
	for i := 0; i < len(s); i++ {
		if p == 0 {
			return 0
		}
		p = trie[p].child[s[i]-'a']
	}
	if p == 0 {
		return 0
	}
	return trie[p].val
}

func trieUpdate(root int, s string, val int) int {
	newRoot := cloneTrie(root)
	pPrev := root
	p := newRoot
	for i := 0; i < len(s); i++ {
		ch := s[i] - 'a'
		var nextPrev int
		if pPrev != 0 {
			nextPrev = trie[pPrev].child[ch]
		}
		next := cloneTrie(nextPrev)
		trie[p].child[ch] = next
		pPrev = nextPrev
		p = next
	}
	trie[p].val = val
	return newRoot
}

// persistent segment tree node
type stNode struct {
	left, right int
	val         int
}

var st []stNode

func cloneST(idx int) int {
	st = append(st, stNode{})
	newIdx := len(st) - 1
	if idx != 0 {
		st[newIdx] = st[idx]
	}
	return newIdx
}

func stUpdate(idx, l, r, pos, delta int) int {
	newIdx := cloneST(idx)
	if l == r {
		st[newIdx].val = st[idx].val + delta
		return newIdx
	}
	mid := (l + r) >> 1
	if pos <= mid {
		st[newIdx].left = stUpdate(st[idx].left, l, mid, pos, delta)
	} else {
		st[newIdx].right = stUpdate(st[idx].right, mid+1, r, pos, delta)
	}
	st[newIdx].val = st[st[newIdx].left].val + st[st[newIdx].right].val
	return newIdx
}

func stQuery(idx, l, r, pos int) int {
	if idx == 0 || pos <= 0 {
		return 0
	}
	if r <= pos {
		return st[idx].val
	}
	mid := (l + r) >> 1
	if pos <= mid {
		return stQuery(st[idx].left, l, mid, pos)
	}
	return st[st[idx].left].val + stQuery(st[idx].right, mid+1, r, pos)
}

const (
	typSet = iota
	typRemove
	typQuery
	typUndo
)

type Operation struct {
	t    int
	name string
	val  int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	if _, err := fmt.Fscan(in, &q); err != nil {
		return
	}
	ops := make([]Operation, q)
	values := make([]int, 0)
	for i := 0; i < q; i++ {
		var typ string
		fmt.Fscan(in, &typ)
		switch typ {
		case "set":
			var name string
			var x int
			fmt.Fscan(in, &name, &x)
			ops[i] = Operation{t: typSet, name: name, val: x}
			values = append(values, x)
		case "remove":
			var name string
			fmt.Fscan(in, &name)
			ops[i] = Operation{t: typRemove, name: name}
		case "query":
			var name string
			fmt.Fscan(in, &name)
			ops[i] = Operation{t: typQuery, name: name}
		case "undo":
			var d int
			fmt.Fscan(in, &d)
			ops[i] = Operation{t: typUndo, val: d}
		}
	}
	sort.Ints(values)
	uniq := make([]int, 0, len(values))
	for _, v := range values {
		if len(uniq) == 0 || uniq[len(uniq)-1] != v {
			uniq = append(uniq, v)
		}
	}
	compress := make(map[int]int, len(uniq))
	for i, v := range uniq {
		compress[v] = i + 1
	}
	for i := range ops {
		if ops[i].t == typSet {
			ops[i].val = compress[ops[i].val]
		}
	}

	trie = make([]trieNode, 1)
	st = make([]stNode, 1)
	trieRoots := make([]int, q+1)
	segRoots := make([]int, q+1)
	trieRoots[0] = 0
	segRoots[0] = 0

	for day := 1; day <= q; day++ {
		op := ops[day-1]
		parent := day - 1
		if op.t == typUndo {
			parent = day - op.val - 1
		}
		trRoot := trieRoots[parent]
		sgRoot := segRoots[parent]
		switch op.t {
		case typSet:
			old := trieQuery(trRoot, op.name)
			if old != 0 {
				sgRoot = stUpdate(sgRoot, 1, len(uniq), old, -1)
			}
			trRoot = trieUpdate(trRoot, op.name, op.val)
			sgRoot = stUpdate(sgRoot, 1, len(uniq), op.val, 1)
		case typRemove:
			old := trieQuery(trRoot, op.name)
			if old != 0 {
				trRoot = trieUpdate(trRoot, op.name, 0)
				sgRoot = stUpdate(sgRoot, 1, len(uniq), old, -1)
			}
		case typQuery:
			p := trieQuery(trRoot, op.name)
			if p == 0 {
				fmt.Fprintln(out, -1)
			} else {
				res := stQuery(sgRoot, 1, len(uniq), p-1)
				fmt.Fprintln(out, res)
			}
		case typUndo:
			// nothing more to do; roots already from parent
		}
		trieRoots[day] = trRoot
		segRoots[day] = sgRoot
	}
}
