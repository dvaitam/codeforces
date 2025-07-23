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

func insert(root *node, x int) {
	cur := root
	for b := maxBit; b >= 0; b-- {
		bit := (x >> b) & 1
		if cur.child[bit] == nil {
			cur.child[bit] = &node{}
		}
		cur = cur.child[bit]
		cur.cnt++
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

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	prefix := 0
	root := &node{}
	insert(root, 0)
	total := int64(0)
	inserted := 1
	for i := 0; i < n; i++ {
		var a int
		fmt.Fscan(in, &a)
		prefix ^= a
		less := countLess(root, prefix, k)
		total += int64(inserted - less)
		insert(root, prefix)
		inserted++
	}
	fmt.Fprintln(out, total)
}
