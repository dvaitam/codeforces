package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxBit = 30 // since a_i, k <= 1e9 < 2^30

type trieNode struct {
	child  [2]*trieNode
	maxPos int // maximum position among numbers in the subtree
}

func (t *trieNode) insert(x, pos int) {
	cur := t
	if pos > cur.maxPos {
		cur.maxPos = pos
	}
	for b := maxBit; b >= 0; b-- {
		v := (x >> b) & 1
		if cur.child[v] == nil {
			cur.child[v] = &trieNode{maxPos: pos}
		} else if pos > cur.child[v].maxPos {
			cur.child[v].maxPos = pos
		}
		cur = cur.child[v]
	}
}

// query returns the maximum position of a number y inserted into the trie
// such that (x xor y) >= k. Returns 0 if no such number exists.
func (t *trieNode) query(x, k int) int {
	var dfs func(node *trieNode, bit int, greater bool) int
	dfs = func(node *trieNode, bit int, greater bool) int {
		if node == nil {
			return 0
		}
		if greater || bit < 0 {
			return node.maxPos
		}
		kb := (k >> bit) & 1
		xb := (x >> bit) & 1
		res := 0
		if kb == 0 {
			// Branch where xor bit becomes 1 turns the prefix greater immediately.
			b1 := xb ^ 1
			if node.child[b1] != nil && node.child[b1].maxPos > res {
				res = node.child[b1].maxPos
			}
			// Branch where xor bit stays 0 keeps equality.
			b0 := xb
			if tmp := dfs(node.child[b0], bit-1, false); tmp > res {
				res = tmp
			}
		} else { // kb == 1
			// xor bit must be 1 to stay >= so we continue with equality.
			b1 := xb ^ 1
			if tmp := dfs(node.child[b1], bit-1, false); tmp > res {
				res = tmp
			}
		}
		return res
	}
	return dfs(t, maxBit, false)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		if k == 0 {
			fmt.Fprintln(out, 1)
			continue
		}

		root := &trieNode{}
		ans := n + 1
		for i, val := range a {
			pos := root.query(val, k)
			if pos > 0 {
				if cand := i - pos + 2; cand < ans { // i is zero-based
					ans = cand
				}
			}
			root.insert(val, i+1) // store positions as 1-based
		}

		if ans == n+1 {
			fmt.Fprintln(out, -1)
		} else {
			fmt.Fprintln(out, ans)
		}
	}
}
