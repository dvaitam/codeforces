package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"time"
)

type Node struct {
	val         int64
	id          int
	pr          int
	left, right *Node
	sz          int
	mx          int64
}

func sz(n *Node) int {
	if n == nil {
		return 0
	}
	return n.sz
}

func upd(n *Node) {
	if n == nil {
		return
	}
	n.sz = 1 + sz(n.left) + sz(n.right)
	n.mx = n.val
	if n.left != nil && n.left.mx > n.mx {
		n.mx = n.left.mx
	}
	if n.right != nil && n.right.mx > n.mx {
		n.mx = n.right.mx
	}
}

func merge(a, b *Node) *Node {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if a.pr > b.pr {
		a.right = merge(a.right, b)
		upd(a)
		return a
	}
	b.left = merge(a, b.left)
	upd(b)
	return b
}

func split(n *Node, k int) (l, r *Node) {
	if n == nil {
		return nil, nil
	}
	if sz(n.left) >= k {
		l0, r0 := split(n.left, k)
		n.left = r0
		upd(n)
		return l0, n
	}
	l0, r0 := split(n.right, k-sz(n.left)-1)
	n.right = l0
	upd(n)
	return n, r0
}

func findLastGE(n *Node, v int64) int {
	if n == nil || n.mx < v {
		return 0
	}
	if n.right != nil && n.right.mx >= v {
		idx := findLastGE(n.right, v)
		if idx > 0 {
			return sz(n.left) + 1 + idx
		}
	}
	if n.val >= v {
		return sz(n.left) + 1
	}
	return findLastGE(n.left, v)
}

func popFront(root **Node) *Node {
	l, r := split(*root, 1)
	*root = r
	return l
}

func insertAfter(root **Node, pos int, node *Node) {
	if pos == 0 {
		*root = merge(node, *root)
		return
	}
	l, r := split(*root, pos)
	*root = merge(merge(l, node), r)
}

func solve() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	rand.Seed(time.Now().UnixNano())
	for ; t > 0; t-- {
		var n, D int
		fmt.Fscan(in, &n, &D)
		k := make([]int64, n+1)
		s := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &k[i], &s[i])
		}
		var root *Node
		for i := 1; i <= n; i++ {
			node := &Node{val: k[i], id: i, pr: rand.Int(), sz: 1, mx: k[i]}
			root = merge(root, node)
		}
		first := make([]int, n+1)
		events := make(map[int][]int)
		served := 0
		lastTime := -1
		for tmin := 1; tmin <= D; tmin++ {
			if root != nil {
				nd := popFront(&root)
				id := nd.id
				if first[id] == 0 {
					first[id] = tmin
					served++
					if served == n {
						lastTime = tmin
					}
				}
				events[tmin+s[id]] = append(events[tmin+s[id]], id)
			}
			if ids, ok := events[tmin]; ok {
				sort.Slice(ids, func(i, j int) bool { return s[ids[i]] < s[ids[j]] })
				for _, id := range ids {
					nd := &Node{val: k[id], id: id, pr: rand.Int(), sz: 1, mx: k[id]}
					pos := findLastGE(root, k[id])
					insertAfter(&root, pos, nd)
				}
				delete(events, tmin)
			}
		}
		if served < n {
			fmt.Fprintln(out, -1)
		} else if lastTime <= D {
			fmt.Fprintln(out, lastTime)
		} else {
			fmt.Fprintln(out, -1)
		}
	}
}

func main() {
	solve()
}
