package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"time"
)

// Treap node
type Node struct {
	key, prio   int
	left, right *Node
}

func split(root *Node, key int) (l, r *Node) {
	if root == nil {
		return nil, nil
	}
	if root.key < key {
		l1, r1 := split(root.right, key)
		root.right = l1
		return root, r1
	}
	l1, r1 := split(root.left, key)
	root.left = r1
	return l1, root
}

func merge(l, r *Node) *Node {
	if l == nil {
		return r
	}
	if r == nil {
		return l
	}
	if l.prio > r.prio {
		l.right = merge(l.right, r)
		return l
	}
	r.left = merge(l, r.left)
	return r
}

func insert(root *Node, node *Node) *Node {
	if root == nil {
		return node
	}
	if node.prio > root.prio {
		l, r := split(root, node.key)
		node.left, node.right = l, r
		return node
	}
	if node.key < root.key {
		root.left = insert(root.left, node)
	} else {
		root.right = insert(root.right, node)
	}
	return root
}

func erase(root *Node, key int) *Node {
	if root == nil {
		return nil
	}
	if key == root.key {
		return merge(root.left, root.right)
	}
	if key < root.key {
		root.left = erase(root.left, key)
	} else {
		root.right = erase(root.right, key)
	}
	return root
}

func predecessor(root *Node, key int) int {
	res := -1000000000
	for root != nil {
		if root.key < key {
			if root.key > res {
				res = root.key
			}
			root = root.right
		} else {
			root = root.left
		}
	}
	return res
}

func successor(root *Node, key int) int {
	res := 1000000000
	for root != nil {
		if root.key > key {
			if root.key < res {
				res = root.key
			}
			root = root.left
		} else {
			root = root.right
		}
	}
	return res
}

func main() {
	rand.Seed(time.Now().UnixNano())
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	fmt.Fscan(reader, &n, &m)
	cover := make([]int, n)
	type evt struct{ pos, x, t int }
	endpoints := make([]evt, 0, 2*m)
	for i := 0; i < m; i++ {
		var x, l, r int
		fmt.Fscan(reader, &x, &l, &r)
		x--
		l--
		endpoints = append(endpoints, evt{l, x, -1})
		endpoints = append(endpoints, evt{r, x, 1})
	}
	sort.Slice(endpoints, func(i, j int) bool {
		return endpoints[i].pos < endpoints[j].pos
	})
	// initialize treap with sentinels -1 and n
	var root *Node
	root = insert(root, &Node{key: -1, prio: rand.Int()})
	root = insert(root, &Node{key: n, prio: rand.Int()})
	edges := make([][2]int, 0, 4*m)
	add := make([]int, 0)
	del := make([]int, 0)
	for l := 0; l < len(endpoints); {
		r := l
		for r < len(endpoints) && endpoints[r].pos == endpoints[l].pos {
			r++
		}
		add = add[:0]
		del = del[:0]
		for i := l; i < r; i++ {
			e := endpoints[i]
			if e.t == -1 {
				if cover[e.x] == 0 {
					add = append(add, e.x)
				}
				cover[e.x]++
			} else {
				cover[e.x]--
				if cover[e.x] == 0 {
					del = append(del, e.x)
				}
			}
		}
		// deletions
		for _, x := range del {
			root = erase(root, x)
		}
		// additions
		for _, x := range add {
			root = insert(root, &Node{key: x, prio: rand.Int()})
			p := predecessor(root, x)
			s := successor(root, x)
			edges = append(edges, [2]int{p, x})
			edges = append(edges, [2]int{x, s})
		}
		l = r
	}
	sort.Slice(edges, func(i, j int) bool {
		if edges[i][0] != edges[j][0] {
			return edges[i][0] < edges[j][0]
		}
		return edges[i][1] < edges[j][1]
	})
	dp := make([]int, n+2)
	g := make([]int, n+2)
	for _, e := range edges {
		u, v := e[0], e[1]
		ui, vi := u+1, v+1
		if dp[ui]+1 > dp[vi] {
			dp[vi] = dp[ui] + 1
			g[vi] = ui
		}
	}
	res := n + 1 - dp[n+1]
	fmt.Fprintln(writer, res)
	list := make([]int, 0, res)
	for i := n + 1; i > 0; {
		prev := g[i]
		for j := prev + 1; j < i; j++ {
			list = append(list, j)
		}
		i = prev
	}
	for _, v := range list {
		fmt.Fprintln(writer, v)
	}
}
