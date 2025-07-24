package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxBit = 30

type node struct {
	child [2]*node
	cnt   int
}

func add(root *node, x, delta int) {
	cur := root
	for b := maxBit; b >= 0; b-- {
		bit := (x >> b) & 1
		if cur.child[bit] == nil {
			cur.child[bit] = &node{}
		}
		cur = cur.child[bit]
		cur.cnt += delta
	}
}

func countLess(root *node, x, k int) int {
	cur := root
	res := 0
	for b := maxBit; b >= 0; b-- {
		if cur == nil {
			break
		}
		xBit := (x >> b) & 1
		kBit := (k >> b) & 1
		if kBit == 1 {
			if cur.child[xBit] != nil {
				res += cur.child[xBit].cnt
			}
			cur = cur.child[xBit^1]
		} else {
			cur = cur.child[xBit]
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

	root := &node{}
	for i := 0; i < q; i++ {
		var t int
		if _, err := fmt.Fscan(in, &t); err != nil {
			return
		}
		if t == 1 {
			var p int
			fmt.Fscan(in, &p)
			add(root, p, 1)
		} else if t == 2 {
			var p int
			fmt.Fscan(in, &p)
			add(root, p, -1)
		} else if t == 3 {
			var p, l int
			fmt.Fscan(in, &p, &l)
			ans := countLess(root, p, l)
			fmt.Fprintln(out, ans)
		}
	}
}
