package main

import (
	"bufio"
	"fmt"
	"os"
)

type state struct {
	node   int
	parity int
}

type stackFrame struct {
	node int
	idx  int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	adj := make([][]int, n)
	for i := 0; i < n; i++ {
		var c int
		fmt.Fscan(in, &c)
		adj[i] = make([]int, c)
		for j := 0; j < c; j++ {
			fmt.Fscan(in, &adj[i][j])
			adj[i][j]--
		}
	}

	var s int
	fmt.Fscan(in, &s)
	s--

	visited := make([][2]bool, n)
	parentNode := make([][2]int, n)
	parentParity := make([][2]int, n)
	for i := 0; i < n; i++ {
		parentNode[i][0], parentNode[i][1] = -1, -1
		parentParity[i][0], parentParity[i][1] = -1, -1
	}

	queue := make([]state, 0, 2*n)
	queue = append(queue, state{s, 0})
	visited[s][0] = true
	head := 0
	targetNode, targetParity := -1, -1

	for head < len(queue) {
		cur := queue[head]
		head++
		v, parity := cur.node, cur.parity
		if len(adj[v]) == 0 && parity == 1 {
			targetNode, targetParity = v, parity
			break
		}
		for _, to := range adj[v] {
			np := parity ^ 1
			if !visited[to][np] {
				visited[to][np] = true
				parentNode[to][np] = v
				parentParity[to][np] = parity
				queue = append(queue, state{to, np})
			}
		}
	}

	if targetNode != -1 {
		path := make([]int, 0)
		node, parity := targetNode, targetParity
		for node != -1 {
			path = append(path, node)
			prevNode := parentNode[node][parity]
			prevParity := parentParity[node][parity]
			node = prevNode
			parity = prevParity
		}
		for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
			path[i], path[j] = path[j], path[i]
		}
		fmt.Fprintln(out, "Win")
		for i, node := range path {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, node+1)
		}
		fmt.Fprintln(out)
		return
	}

	reach := make([]bool, n)
	nodes := make([]int, 0, n)
	nodes = append(nodes, s)
	reach[s] = true
	for head = 0; head < len(nodes); head++ {
		v := nodes[head]
		for _, to := range adj[v] {
			if !reach[to] {
				reach[to] = true
				nodes = append(nodes, to)
			}
		}
	}

	color := make([]uint8, n)
	for i := 0; i < n; i++ {
		if !reach[i] || color[i] != 0 {
			continue
		}
		stack := []stackFrame{{i, 0}}
		color[i] = 1
		for len(stack) > 0 {
			top := &stack[len(stack)-1]
			v := top.node
			if top.idx < len(adj[v]) {
				to := adj[v][top.idx]
				top.idx++
				if !reach[to] {
					continue
				}
				if color[to] == 0 {
					color[to] = 1
					stack = append(stack, stackFrame{to, 0})
				} else if color[to] == 1 {
					fmt.Fprintln(out, "Draw")
					return
				}
			} else {
				color[v] = 2
				stack = stack[:len(stack)-1]
			}
		}
	}

	fmt.Fprintln(out, "Lose")
}
