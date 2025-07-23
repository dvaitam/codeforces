package main

// Package main solves problem 842D - Vitya and Strange Lesson.
// See problemD.txt for the statement.

import (
	"bufio"
	"fmt"
	"os"
)

const (
	B = 19
	L = 1 << B
)

type node struct {
	child [2]int
}

func insert(trie *[]node, val int) {
	idx := 0
	for i := B - 1; i >= 0; i-- {
		b := (val >> i) & 1
		next := (*trie)[idx].child[b]
		if next == 0 {
			next = len(*trie)
			*trie = append(*trie, node{})
			(*trie)[idx].child[b] = next
		}
		idx = next
	}
}

func minXor(trie []node, val int) int {
	idx := 0
	res := 0
	for i := B - 1; i >= 0; i-- {
		b := (val >> i) & 1
		if trie[idx].child[b] != 0 {
			idx = trie[idx].child[b]
		} else {
			res |= 1 << i
			idx = trie[idx].child[b^1]
		}
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	present := make([]bool, L)
	for i := 0; i < n; i++ {
		var v int
		fmt.Fscan(in, &v)
		if v < L {
			present[v] = true
		}
	}

	trie := []node{{}}
	for i := 0; i < L; i++ {
		if !present[i] {
			insert(&trie, i)
		}
	}

	cur := 0
	for j := 0; j < m; j++ {
		var x int
		fmt.Fscan(in, &x)
		cur ^= x
		ans := minXor(trie, cur)
		fmt.Fprintln(out, ans)
	}
}
