package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m, k, q int
	if _, err := fmt.Fscan(in, &n, &m, &k, &q); err != nil {
		return
	}
	owner := make([]int, m+1)
	parent := make([]int, n+1)
	for i := 0; i < k; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		if owner[b] == 0 {
			owner[b] = a
		} else {
			parent[a] = owner[b]
		}
	}
	children := make([][]int, n+1)
	for i := 1; i <= n; i++ {
		p := parent[i]
		if p != 0 {
			children[p] = append(children[p], i)
		}
	}
	tin := make([]int, n+1)
	tout := make([]int, n+1)
	size := make([]int, n+1)
	timer := 1
	type pair struct{ v, stage int }
	stack := make([]pair, 0)
	for i := 1; i <= n; i++ {
		if parent[i] == 0 {
			stack = append(stack, pair{i, 0})
			for len(stack) > 0 {
				cur := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				v := cur.v
				if cur.stage == 0 {
					tin[v] = timer
					timer++
					stack = append(stack, pair{v, 1})
					for j := len(children[v]) - 1; j >= 0; j-- {
						stack = append(stack, pair{children[v][j], 0})
					}
				} else {
					sz := 1
					for _, u := range children[v] {
						sz += size[u]
					}
					size[v] = sz
					tout[v] = timer
					timer++
				}
			}
		}
	}
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	isAncestor := func(a, b int) bool {
		return tin[a] <= tin[b] && tout[b] <= tout[a]
	}
	for i := 0; i < q; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		w := owner[y]
		if w != 0 && isAncestor(x, w) {
			fmt.Fprintln(out, size[x])
		} else {
			fmt.Fprintln(out, 0)
		}
	}
}
