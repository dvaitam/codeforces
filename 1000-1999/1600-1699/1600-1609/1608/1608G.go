package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct {
	to   int
	char byte
}

var adj [][]Edge

func findPath(u, v int, parent []int) string {
	// BFS from u to find path to v
	// Since this is a tree, we can perform DFS to record parent
	n := len(adj)
	if parent == nil {
		parent = make([]int, n)
		for i := range parent {
			parent[i] = -1
		}
		stack := []int{u}
		parent[u] = u
		for len(stack) > 0 {
			x := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if x == v {
				break
			}
			for _, e := range adj[x] {
				if parent[e.to] == -1 {
					parent[e.to] = x
					stack = append(stack, e.to)
				}
			}
		}
	}
	// reconstruct path
	path := []byte{}
	cur := v
	for cur != u {
		p := parent[cur]
		// find char from p to cur
		ch := byte('?')
		for _, e := range adj[p] {
			if e.to == cur {
				ch = e.char
				break
			}
		}
		path = append(path, ch)
		cur = p
	}
	// reverse
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return string(path)
}

func countOccurrences(s, pat string) int {
	if len(pat) == 0 {
		return 0
	}
	count := 0
	for i := 0; i+len(pat) <= len(s); i++ {
		if s[i:i+len(pat)] == pat {
			count++
		}
	}
	return count
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m, q int
	if _, err := fmt.Fscan(reader, &n, &m, &q); err != nil {
		return
	}

	adj = make([][]Edge, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		var c string
		fmt.Fscan(reader, &u, &v, &c)
		ch := c[0]
		adj[u] = append(adj[u], Edge{v, ch})
		adj[v] = append(adj[v], Edge{u, ch})
	}

	strs := make([]string, m+1)
	for i := 1; i <= m; i++ {
		fmt.Fscan(reader, &strs[i])
	}

	for ; q > 0; q-- {
		var u, v, l, r int
		fmt.Fscan(reader, &u, &v, &l, &r)
		pat := findPath(u, v, nil)
		total := 0
		for i := l; i <= r; i++ {
			total += countOccurrences(strs[i], pat)
		}
		fmt.Fprintln(writer, total)
	}
}
