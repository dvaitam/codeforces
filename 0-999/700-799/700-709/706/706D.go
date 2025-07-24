package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxBit = 30

type node struct {
	child [2]int
	cnt   int
}

var trie []node

func insert(x int) {
	idx := 0
	trie[idx].cnt++
	for b := maxBit; b >= 0; b-- {
		bit := (x >> b) & 1
		nxt := trie[idx].child[bit]
		if nxt == 0 {
			trie = append(trie, node{})
			nxt = len(trie) - 1
			trie[idx].child[bit] = nxt
		}
		idx = nxt
		trie[idx].cnt++
	}
}

func removeVal(x int) {
	idx := 0
	trie[idx].cnt--
	for b := maxBit; b >= 0; b-- {
		bit := (x >> b) & 1
		idx = trie[idx].child[bit]
		trie[idx].cnt--
	}
}

func maxXor(x int) int {
	idx := 0
	res := 0
	for b := maxBit; b >= 0; b-- {
		bit := (x >> b) & 1
		prefer := bit ^ 1
		nxt := trie[idx].child[prefer]
		if nxt != 0 && trie[nxt].cnt > 0 {
			res |= 1 << b
			idx = nxt
		} else {
			idx = trie[idx].child[bit]
		}
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	if _, err := fmt.Fscan(in, &q); err != nil {
		return
	}

	trie = make([]node, 1)
	insert(0)

	for ; q > 0; q-- {
		var op string
		var x int
		fmt.Fscan(in, &op, &x)
		switch op {
		case "+":
			insert(x)
		case "-":
			removeVal(x)
		case "?":
			ans := maxXor(x)
			fmt.Fprintln(out, ans)
		}
	}
}
