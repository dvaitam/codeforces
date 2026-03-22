package main

import (
	"bufio"
	"fmt"
	"os"
)

type Node struct {
	l, r *Node
	sum  int
}

func getSum(n *Node) int {
	if n == nil {
		return 0
	}
	return n.sum
}

func update(root *Node, l, r, idx, val int) *Node {
	if l == r {
		return &Node{sum: val}
	}
	mid := l + (r-l)/2
	res := &Node{}
	if root != nil {
		res.l = root.l
		res.r = root.r
	}
	if idx <= mid {
		res.l = update(res.l, l, mid, idx, val)
	} else {
		res.r = update(res.r, mid+1, r, idx, val)
	}
	res.sum = getSum(res.l) + getSum(res.r)
	return res
}

func query(root *Node, l, r, idx int) int {
	if root == nil {
		return 0
	}
	if l == r {
		return root.sum
	}
	mid := l + (r-l)/2
	if idx <= mid {
		return query(root.l, l, mid, idx)
	}
	return query(root.r, mid+1, r, idx)
}

func getOnes(root *Node, l, r int, ones *[]int) {
	if root == nil || root.sum == 0 {
		return
	}
	if l == r {
		*ones = append(*ones, l)
		return
	}
	mid := l + (r-l)/2
	getOnes(root.l, l, mid, ones)
	getOnes(root.r, mid+1, r, ones)
}

type State struct {
	u    int
	c    int
	root *Node
}

type Info struct {
	prevU int
	prevC int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	adj := make([][]int, n+1)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	visited := make([][2]bool, n+1)
	parent := make([][2]Info, n+1)

	next := make([]int, n+2)
	prev := make([]int, n+2)
	for i := 2; i <= n; i++ {
		next[i] = i + 1
		prev[i] = i - 1
	}
	head := 2
	prev[2] = 0
	next[n] = n + 1

	removeNode := func(v int) {
		p := prev[v]
		nxt := next[v]
		if p != 0 {
			next[p] = nxt
		} else {
			head = nxt
		}
		if nxt != n+1 && nxt != 0 {
			prev[nxt] = p
		}
		prev[v] = 0
		next[v] = 0
	}

	isAdj := make([]bool, n+1)

	q := []State{{u: 1, c: 0, root: nil}}
	visited[1][0] = true

	found := false
	for len(q) > 0 {
		curr := q[0]
		q = q[1:]

		u := curr.u
		c := curr.c
		root := curr.root

		if u == n {
			found = true
			break
		}

		for _, v := range adj[u] {
			isAdj[v] = true
		}

		newRoot := update(root, 1, n, u, 1-c)

		for _, v := range adj[u] {
			dv := query(root, 1, n, v)
			if dv == c {
				if !visited[v][c] {
					visited[v][c] = true
					parent[v][c] = Info{u, c}
					q = append(q, State{v, c, newRoot})
					if c == 0 && prev[v] != 0 || head == v {
						removeNode(v)
					}
				}
			}
		}

		if c == 0 {
			var ones []int
			getOnes(root, 1, n, &ones)
			for _, v := range ones {
				if v == u {
					continue
				}
				if !isAdj[v] {
					if !visited[v][1] {
						visited[v][1] = true
						parent[v][1] = Info{u, c}
						q = append(q, State{v, 1, newRoot})
					}
				}
			}
		} else {
			var toRemove []int
			for v := head; v != 0 && v <= n; v = next[v] {
				if v == u {
					continue
				}
				if !isAdj[v] {
					if query(root, 1, n, v) == 0 {
						visited[v][0] = true
						toRemove = append(toRemove, v)
						parent[v][0] = Info{u, c}
						q = append(q, State{v, 0, newRoot})
					}
				}
			}
			for _, v := range toRemove {
				removeNode(v)
			}
		}

		for _, v := range adj[u] {
			isAdj[v] = false
		}
	}

	if !found && !visited[n][0] && !visited[n][1] {
		fmt.Println("-1")
		return
	}

	currU := n
	currC := 0
	if !visited[n][0] {
		currC = 1
	}

	path := []int{n}
	for currU != 1 || currC != 0 {
		p := parent[currU][currC]
		path = append(path, p.prevU)
		currU = p.prevU
		currC = p.prevC
	}

	fmt.Println(len(path) - 1)
	writer := bufio.NewWriter(os.Stdout)
	for i := len(path) - 1; i >= 0; i-- {
		fmt.Fprintf(writer, "%d ", path[i])
	}
	fmt.Fprintln(writer)
	writer.Flush()
}
