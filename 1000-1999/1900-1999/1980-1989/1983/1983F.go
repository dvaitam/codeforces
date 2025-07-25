package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// simple binary trie for non-negative integers up to 2^30
// only supports insertion and querying minimum xor

type node struct {
	child [2]*node
}

func (n *node) insert(x int) {
	cur := n
	for i := 30; i >= 0; i-- {
		b := (x >> i) & 1
		if cur.child[b] == nil {
			cur.child[b] = &node{}
		}
		cur = cur.child[b]
	}
}

func (n *node) minXor(x int) int {
	cur := n
	res := 0
	for i := 30; i >= 0; i-- {
		b := (x >> i) & 1
		if cur.child[b] != nil {
			cur = cur.child[b]
		} else {
			res |= 1 << i
			cur = cur.child[1-b]
		}
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}

		var vals []int
		for l := 0; l < n; l++ {
			root := &node{}
			root.insert(a[l])
			minVal := int(^uint(0) >> 1) // max int
			for r := l + 1; r < n; r++ {
				v := root.minXor(a[r])
				if v < minVal {
					minVal = v
				}
				root.insert(a[r])
				vals = append(vals, minVal)
			}
		}
		sort.Ints(vals)
		if k <= len(vals) {
			fmt.Fprintln(writer, vals[k-1])
		} else {
			fmt.Fprintln(writer, -1)
		}
	}
}
