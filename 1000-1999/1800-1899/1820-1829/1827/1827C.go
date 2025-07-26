package main

import (
	"bufio"
	"fmt"
	"os"
)

// Node represents a stack node for persistent stack.
type Node struct {
	ch   byte
	prev *Node
}

type key struct {
	prev *Node
	ch   byte
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		var s string
		fmt.Fscan(in, &s)

		// counts maps stack state pointer to counts for parity 0 and 1
		counts := make(map[*Node][2]int64)
		// cache for created nodes to reuse identical states
		nodeMap := make(map[key]*Node)
		// initial state: empty stack, parity 0
		empty := (*Node)(nil)
		counts[empty] = [2]int64{1, 0}

		var ans int64
		parity := 0
		top := empty

		for i := 0; i < n; i++ {
			c := s[i]
			if top != nil && top.ch == c {
				top = top.prev
			} else {
				k := key{prev: top, ch: c}
				if node, ok := nodeMap[k]; ok {
					top = node
				} else {
					node := &Node{ch: c, prev: top}
					nodeMap[k] = node
					top = node
				}
			}
			parity ^= 1
			arr := counts[top]
			ans += arr[parity]
			arr[parity]++
			counts[top] = arr
		}
		fmt.Fprintln(out, ans)
	}
}
