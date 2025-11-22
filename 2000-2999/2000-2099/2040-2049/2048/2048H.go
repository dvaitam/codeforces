package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 998244353

type nodeKey struct {
	term           bool
	child0, child1 int
}

type node struct {
	term           bool
	child0, child1 int
}

// pack a pair into a uint64 with a <= b enforced by caller.
func pairKey(a, b int) uint64 {
	return uint64(a)<<32 | uint64(b)
}

type canon struct {
	nodes []node
	idMap map[nodeKey]int
	merge map[uint64]int
}

func newCanon() *canon {
	c := &canon{
		nodes: []node{{term: false, child0: 0, child1: 0}}, // id 0: empty set
		idMap: make(map[nodeKey]int),
		merge: make(map[uint64]int),
	}
	c.idMap[nodeKey{term: false, child0: 0, child1: 0}] = 0
	return c
}

func (c *canon) get(term bool, c0, c1 int) int {
	key := nodeKey{term: term, child0: c0, child1: c1}
	if id, ok := c.idMap[key]; ok {
		return id
	}
	id := len(c.nodes)
	c.nodes = append(c.nodes, node{term: term, child0: c0, child1: c1})
	c.idMap[key] = id
	return id
}

// union of two trie nodes represented by ids, iterative to avoid deep recursion.
func (c *canon) union(a, b int) int {
	if a == 0 {
		return b
	}
	if b == 0 {
		return a
	}
	if a == b {
		return a
	}
	if a > b {
		a, b = b, a
	}
	key := pairKey(a, b)
	if v, ok := c.merge[key]; ok {
		return v
	}

	type item struct{ a, b int }
	stack := []item{{a, b}}
	for len(stack) > 0 {
		top := stack[len(stack)-1]
		aa, bb := top.a, top.b
		if aa == 0 {
			if aa > bb {
				aa, bb = bb, aa
			}
			c.merge[pairKey(aa, bb)] = bb
			stack = stack[:len(stack)-1]
			continue
		}
		if bb == 0 {
			if aa > bb {
				aa, bb = bb, aa
			}
			c.merge[pairKey(aa, bb)] = aa
			stack = stack[:len(stack)-1]
			continue
		}
		if aa == bb {
			c.merge[pairKey(aa, bb)] = aa
			stack = stack[:len(stack)-1]
			continue
		}
		if aa > bb {
			aa, bb = bb, aa
		}
		k := pairKey(aa, bb)
		if _, ok := c.merge[k]; ok {
			stack = stack[:len(stack)-1]
			continue
		}

		na := c.nodes[aa]
		nb := c.nodes[bb]

		needChild := false
		var u0, u1 int

		// child 0
		xa, xb := na.child0, nb.child0
		if xa == 0 && xb == 0 {
			u0 = 0
		} else {
			if xa > xb {
				xa, xb = xb, xa
			}
			k0 := pairKey(xa, xb)
			if v, ok := c.merge[k0]; ok {
				u0 = v
			} else {
				stack = append(stack, item{xa, xb})
				needChild = true
			}
		}

		// child 1
		ya, yb := na.child1, nb.child1
		if ya == 0 && yb == 0 {
			u1 = 0
		} else {
			if ya > yb {
				ya, yb = yb, ya
			}
			k1 := pairKey(ya, yb)
			if v, ok := c.merge[k1]; ok {
				u1 = v
			} else {
				stack = append(stack, item{ya, yb})
				needChild = true
			}
		}

		if needChild {
			continue
		}

		id := c.get(na.term || nb.term, u0, u1)
		c.merge[k] = id
		stack = stack[:len(stack)-1]
	}

	if a > b {
		a, b = b, a
	}
	return c.merge[pairKey(a, b)]
}

func (c *canon) countStrings(root int) int {
	n := len(c.nodes)
	dp := make([]int, n)
	vis := make([]bool, n)
	type frame struct {
		id     int
		loaded bool
	}
	stack := []frame{{id: root, loaded: false}}
	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if vis[top.id] {
			continue
		}
		if top.loaded {
			vis[top.id] = true
			nv := c.nodes[top.id]
			val := 0
			if nv.term {
				val = 1
			}
			if nv.child0 != 0 {
				val += dp[nv.child0]
				if val >= mod {
					val -= mod
				}
			}
			if nv.child1 != 0 {
				val += dp[nv.child1]
				if val >= mod {
					val -= mod
				}
			}
			dp[top.id] = val
			continue
		}
		// load children first
		stack = append(stack, frame{id: top.id, loaded: true})
		nv := c.nodes[top.id]
		if nv.child0 != 0 {
			stack = append(stack, frame{id: nv.child0, loaded: false})
		}
		if nv.child1 != 0 {
			stack = append(stack, frame{id: nv.child1, loaded: false})
		}
	}
	return dp[root]
}

func solveCase(s string) int {
	n := len(s)
	c := newCanon()
	leaf := c.get(true, 0, 0) // set containing only empty string

	F := make([]int, n+1)
	S := make([]int, n+1)
	F[n] = leaf
	S[n] = leaf

	nextOne := n
	B := F[n] // union of F[k] for k in [i+1, nextOne] when i starts at n-1

	for i := n - 1; i >= 0; i-- {
		if s[i] == '0' {
			zeroChild := B
			oneChild := 0
			if nextOne != n {
				oneChild = S[nextOne+1]
			}
			F[i] = c.get(false, zeroChild, oneChild)
			B = c.union(F[i], B) // update union_{k=i}^{nextOne}
		} else {
			oneChild := S[i+1]
			F[i] = c.get(false, 0, oneChild)
			nextOne = i
			B = F[i] // reset union up to new nextOne
		}
		S[i] = c.union(F[i], S[i+1])
	}

	return c.countStrings(F[0])
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var s string
		fmt.Fscan(in, &s)
		fmt.Fprintln(out, solveCase(s)%mod)
	}
}
