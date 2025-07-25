package main

import (
	"bufio"
	"fmt"
	"os"
)

// Node represents a trie node
type Node struct {
	next [26]int32
	cnt  int32
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(reader, &n)
	strs := make([]string, n)

	nodes := make([]Node, 1) // node 0 is root
	totalLen := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &strs[i])
		node := 0
		for j := 0; j < len(strs[i]); j++ {
			idx := strs[i][j] - 'a'
			child := nodes[node].next[idx]
			if child == 0 {
				nodes = append(nodes, Node{})
				nodes[node].next[idx] = int32(len(nodes) - 1)
				child = nodes[node].next[idx]
			}
			node = int(child)
			nodes[node].cnt++
		}
		totalLen += len(strs[i])
	}

	var sumLCP int64
	for _, s := range strs {
		node := 0
		for i := len(s) - 1; i >= 0; i-- {
			idx := s[i] - 'a'
			child := nodes[node].next[idx]
			if child == 0 {
				break
			}
			node = int(child)
			sumLCP += int64(nodes[node].cnt)
		}
	}

	totalLen64 := int64(totalLen)
	n64 := int64(n)
	result := 2*n64*totalLen64 - 2*sumLCP
	fmt.Println(result)
}
