package main

import (
	"bufio"
	"fmt"
	"os"
)

type ACNode struct {
	next [26]int
	fail int
}

type Query struct {
	node  int
	index int
}

type Edge struct {
	to   int
	char byte
}

type BIT struct {
	n    int
	tree []int
}

func NewBIT(n int) *BIT {
	return &BIT{n: n, tree: make([]int, n+2)}
}

func (b *BIT) Add(i, v int) {
	for ; i <= b.n; i += i & -i {
		b.tree[i] += v
	}
}

func (b *BIT) Sum(i int) int {
	s := 0
	for ; i > 0; i -= i & -i {
		s += b.tree[i]
	}
	return s
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	children := make([][]Edge, n+1)
	for i := 1; i <= n; i++ {
		var t int
		fmt.Fscan(reader, &t)
		if t == 1 {
			var c string
			fmt.Fscan(reader, &c)
			children[0] = append(children[0], Edge{to: i, char: c[0]})
		} else {
			var j int
			var c string
			fmt.Fscan(reader, &j, &c)
			children[j] = append(children[j], Edge{to: i, char: c[0]})
		}
	}

	var m int
	fmt.Fscan(reader, &m)

	nodes := make([]ACNode, 1)
	insert := func(s string) int {
		v := 0
		for i := 0; i < len(s); i++ {
			c := int(s[i] - 'a')
			if nodes[v].next[c] == 0 {
				nodes = append(nodes, ACNode{})
				nodes[v].next[c] = len(nodes) - 1
			}
			v = nodes[v].next[c]
		}
		return v
	}

	queries := make([][]Query, n+1)
	for qi := 0; qi < m; qi++ {
		var idx int
		var t string
		fmt.Fscan(reader, &idx, &t)
		node := insert(t)
		queries[idx] = append(queries[idx], Query{node: node, index: qi})
	}

	// build failure links
	queue := make([]int, 0)
	for c := 0; c < 26; c++ {
		v := nodes[0].next[c]
		if v != 0 {
			queue = append(queue, v)
		}
	}
	for i := 0; i < len(queue); i++ {
		v := queue[i]
		for c := 0; c < 26; c++ {
			u := nodes[v].next[c]
			if u != 0 {
				nodes[u].fail = nodes[nodes[v].fail].next[c]
				queue = append(queue, u)
			} else {
				nodes[v].next[c] = nodes[nodes[v].fail].next[c]
			}
		}
	}

	size := len(nodes)
	childrenFail := make([][]int, size)
	for v := 1; v < size; v++ {
		p := nodes[v].fail
		childrenFail[p] = append(childrenFail[p], v)
	}

	tin := make([]int, size)
	tout := make([]int, size)
	time := 0
	type Frame struct{ node, idx int }
	stack := []Frame{{0, 0}}
	for len(stack) > 0 {
		top := &stack[len(stack)-1]
		if top.idx == 0 {
			time++
			tin[top.node] = time
		}
		if top.idx < len(childrenFail[top.node]) {
			v := childrenFail[top.node][top.idx]
			top.idx++
			stack = append(stack, Frame{v, 0})
		} else {
			tout[top.node] = time
			stack = stack[:len(stack)-1]
		}
	}

	bit := NewBIT(time + 2)
	ans := make([]int, m)

	type SFrame struct{ song, state, idx int }
	sstack := []SFrame{{0, 0, 0}}
	for len(sstack) > 0 {
		fr := &sstack[len(sstack)-1]
		if fr.idx == len(children[fr.song]) {
			if fr.song != 0 {
				bit.Add(tin[fr.state], -1)
			}
			sstack = sstack[:len(sstack)-1]
			continue
		}
		e := children[fr.song][fr.idx]
		fr.idx++
		ns := nodes[fr.state].next[int(e.char-'a')]
		bit.Add(tin[ns], 1)
		for _, q := range queries[e.to] {
			res := bit.Sum(tout[q.node]) - bit.Sum(tin[q.node]-1)
			ans[q.index] = res
		}
		sstack = append(sstack, SFrame{e.to, ns, 0})
	}

	for i := 0; i < m; i++ {
		fmt.Fprintln(writer, ans[i])
	}
}
