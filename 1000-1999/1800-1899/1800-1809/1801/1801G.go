package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Node struct {
	next   [26]int32
	fail   int32
	link   int32
	length int32
}

type BIT struct {
	n    int
	tree []int64
}

func NewBIT(n int) *BIT {
	return &BIT{n: n, tree: make([]int64, n+1)}
}

func (b *BIT) Add(i int, v int64) {
	for i <= b.n {
		b.tree[i] += v
		i += i & -i
	}
}

func (b *BIT) Sum(i int) int64 {
	var s int64
	for i > 0 {
		s += b.tree[i]
		i -= i & -i
	}
	return s
}

type Query struct {
	l, r int
	idx  int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	var t string
	fmt.Fscan(reader, &t)
	// Build Aho-Corasick
	nodes := make([]Node, 1)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(reader, &s)
		u := 0
		for j := 0; j < len(s); j++ {
			c := int(s[j] - 'a')
			if nodes[u].next[c] == 0 {
				nodes = append(nodes, Node{})
				nodes[u].next[c] = int32(len(nodes) - 1)
			}
			u = int(nodes[u].next[c])
		}
		nodes[u].length = int32(len(s))
	}
	// build fail links and output links
	queue := make([]int, 0)
	for c := 0; c < 26; c++ {
		v := int(nodes[0].next[c])
		if v != 0 {
			nodes[v].fail = 0
			queue = append(queue, v)
		}
	}
	for qi := 0; qi < len(queue); qi++ {
		u := queue[qi]
		f := int(nodes[u].fail)
		if nodes[f].length > 0 {
			nodes[u].link = int32(f)
		} else {
			nodes[u].link = nodes[f].link
		}
		for c := 0; c < 26; c++ {
			v := int(nodes[u].next[c])
			if v != 0 {
				nodes[v].fail = nodes[f].next[c]
				queue = append(queue, v)
			} else {
				nodes[u].next[c] = nodes[f].next[c]
			}
		}
	}

	queries := make([]Query, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &queries[i].l, &queries[i].r)
		queries[i].idx = i
	}
	sort.Slice(queries, func(i, j int) bool { return queries[i].r < queries[j].r })

	bit := NewBIT(len(t) + 2)
	ans := make([]int64, m)
	state := int32(0)
	qi := 0
	for pos := 1; pos <= len(t); pos++ {
		c := int(t[pos-1] - 'a')
		state = nodes[state].next[c]
		v := state
		for v != 0 {
			if nodes[v].length > 0 {
				start := pos - int(nodes[v].length) + 1
				if start >= 1 {
					bit.Add(start, 1)
				}
			}
			v = nodes[v].link
		}
		for qi < m && queries[qi].r == pos {
			l := queries[qi].l
			r := queries[qi].r
			ans[queries[qi].idx] = bit.Sum(r) - bit.Sum(l-1)
			qi++
		}
	}
	for ; qi < m; qi++ {
		l := queries[qi].l
		r := queries[qi].r
		ans[queries[qi].idx] = bit.Sum(r) - bit.Sum(l-1)
	}
	for i := 0; i < m; i++ {
		if i > 0 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, ans[i])
	}
	writer.WriteByte('\n')
}
