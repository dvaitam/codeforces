package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func xorUpto(x int) int {
	switch x & 3 {
	case 0:
		return x
	case 1:
		return 1
	case 2:
		return x + 1
	default:
		return 0
	}
}

type trieNode struct {
	child [2]*trieNode
}

func (t *trieNode) insert(x int) {
	node := t
	for i := 20; i >= 0; i-- {
		b := (x >> i) & 1
		if node.child[b] == nil {
			node.child[b] = &trieNode{}
		}
		node = node.child[b]
	}
}

func (t *trieNode) maxXor(x int) int {
	node := t
	ans := 0
	for i := 20; i >= 0; i-- {
		if node == nil {
			break
		}
		b := (x >> i) & 1
		if node.child[1-b] != nil {
			ans |= 1 << i
			node = node.child[1-b]
		} else {
			node = node.child[b]
		}
	}
	return ans
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
	}

	for ; m > 0; m-- {
		var l, r int
		fmt.Fscan(reader, &l, &r)
		l--
		r--
		uniqMap := make(map[int]struct{})
		vals := make([]int, 0, r-l+1)
		for i := l; i <= r; i++ {
			v := arr[i]
			if _, ok := uniqMap[v]; !ok {
				uniqMap[v] = struct{}{}
				vals = append(vals, v)
			}
		}
		sort.Ints(vals)
		root := &trieNode{}
		best := 0
		for _, v := range vals {
			root.insert(xorUpto(v - 1))
			cur := root.maxXor(xorUpto(v))
			if cur > best {
				best = cur
			}
		}
		fmt.Fprintln(writer, best)
	}
}
