package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxBits = 30

type node struct {
	next [2]int
	cnt  int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var l, r int
		fmt.Fscan(in, &l, &r)
		n := r - l + 1
		nodes := make([]node, 1, n*maxBits+5)
		addNode := func() int {
			nodes = append(nodes, node{})
			return len(nodes) - 1
		}
		// insert all values in [l, r]
		for val := l; val <= r; val++ {
			cur := 0
			nodes[cur].cnt++
			for b := 0; b < maxBits; b++ {
				bit := (val >> b) & 1
				nxt := nodes[cur].next[bit]
				if nxt == 0 {
					nxt = addNode()
					nodes[cur].next[bit] = nxt
				}
				cur = nxt
				nodes[cur].cnt++
			}
		}

		ans := make([]int, n)
		var total int64

		for i := 0; i < n; i++ {
			cur := 0
			nodes[cur].cnt--
			val := l + i
			pick := 0
			for b := 0; b < maxBits; b++ {
				desired := ((val >> b) & 1) ^ 1
				nxt := nodes[cur].next[desired]
				if nxt == 0 || nodes[nxt].cnt == 0 {
					desired ^= 1
					nxt = nodes[cur].next[desired]
				}
				if desired == 1 {
					pick |= 1 << b
				}
				cur = nxt
				nodes[cur].cnt--
			}
			ans[i] = pick
			total += int64(pick | val)
		}

		fmt.Fprintln(out, total)
		for i, v := range ans {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
		fmt.Fprintln(out)
	}
}
