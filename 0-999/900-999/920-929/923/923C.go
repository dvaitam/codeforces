package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxBits = 30

type node struct {
	child [2]int
	cnt   int
}

type trie struct {
	nodes []node
}

func newTrie(n int) *trie {
	t := &trie{nodes: make([]node, 1, n*maxBits+1)}
	return t
}

func (t *trie) insert(x int) {
	cur := 0
	t.nodes[cur].cnt++
	for b := maxBits - 1; b >= 0; b-- {
		bit := (x >> b) & 1
		nxt := t.nodes[cur].child[bit]
		if nxt == 0 {
			t.nodes = append(t.nodes, node{})
			nxt = len(t.nodes) - 1
			t.nodes[cur].child[bit] = nxt
		}
		cur = nxt
		t.nodes[cur].cnt++
	}
}

func (t *trie) remove(x int) {
	cur := 0
	t.nodes[cur].cnt--
	for b := maxBits - 1; b >= 0; b-- {
		bit := (x >> b) & 1
		cur = t.nodes[cur].child[bit]
		t.nodes[cur].cnt--
	}
}

func (t *trie) getMin(x int) int {
	cur := 0
	val := 0
	for b := maxBits - 1; b >= 0; b-- {
		bit := (x >> b) & 1
		nxt := t.nodes[cur].child[bit]
		if nxt != 0 && t.nodes[nxt].cnt > 0 {
			cur = nxt
			val |= bit << b
		} else {
			bit ^= 1
			cur = t.nodes[cur].child[bit]
			val |= bit << b
		}
	}
	t.remove(val)
	return val
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	A := make([]int, n)
	P := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &A[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &P[i])
	}

	tr := newTrie(n)
	for _, v := range P {
		tr.insert(v)
	}

	for i, a := range A {
		v := tr.getMin(a)
		if i > 0 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, a^v)
	}
	writer.WriteByte('\n')
}
