package main

import (
	"bufio"
	"fmt"
	"os"
)

type TrieNode struct {
	child [2]*TrieNode
	cnt   int
}

type Trie struct {
	root *TrieNode
}

func newTrie() *Trie {
	return &Trie{root: &TrieNode{}}
}

func (tr *Trie) insert(x int) {
	node := tr.root
	for bit := 29; bit >= 0; bit-- {
		b := (x >> bit) & 1
		if node.child[b] == nil {
			node.child[b] = &TrieNode{}
		}
		node = node.child[b]
		node.cnt++
	}
}

func (tr *Trie) query(x int) int {
	node := tr.root
	if node.child[0] == nil && node.child[1] == nil {
		return 0
	}
	result := 0
	count := 0
	for bit := 29; bit >= 0; bit-- {
		b := (x >> bit) & 1
		t := b ^ 1
		if node.child[t] != nil {
			if count+node.child[t].cnt <= 0 {
				continue
			}
			result += (1 << bit)
			count += node.child[t].cnt
			node = node.child[t]
		} else if node.child[b] != nil {
			node = node.child[b]
		} else {
			break
		}
	}
	return result
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}

		tr := newTrie()
		for _, v := range arr {
			tr.insert(v)
		}

		best := 0
		for _, v := range arr {
			val := 0
			for _, u := range arr {
				val += v ^ u
			}
			if val > best {
				best = val
			}
		}
		fmt.Fprintln(writer, best)
	}
}
