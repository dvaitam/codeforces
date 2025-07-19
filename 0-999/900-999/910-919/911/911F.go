package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	adj := make([][]int, n+1)
	for i := 1; i < n; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}
	// BFS from 1 to find u
	dist1 := make([]int, n+1)
	vis := make([]bool, n+1)
	queue := make([]int, n)
	head, tail := 0, 0
	queue[tail] = 1
	tail++
	vis[1] = true
	u := 1
	for head < tail {
		x := queue[head]
		head++
		for _, y := range adj[x] {
			if !vis[y] {
				vis[y] = true
				dist1[y] = dist1[x] + 1
				queue[tail] = y
				tail++
				if dist1[y] > dist1[u] {
					u = y
				}
			}
		}
	}
	// BFS from u to find v, record parent and depa
	depa := make([]int, n+1)
	parent := make([]int, n+1)
	children := make([]int, n+1)
	for i := range vis {
		vis[i] = false
	}
	head, tail = 0, 0
	queue[tail] = u
	tail++
	vis[u] = true
	parent[u] = 0
	depa[u] = 0
	v := u
	for head < tail {
		x := queue[head]
		head++
		for _, y := range adj[x] {
			if !vis[y] {
				vis[y] = true
				parent[y] = x
				depa[y] = depa[x] + 1
				children[x]++
				queue[tail] = y
				tail++
				if depa[y] > depa[v] {
					v = y
				}
			}
		}
	}
	// BFS from v to get depb
	depb := make([]int, n+1)
	for i := range vis {
		vis[i] = false
	}
	head, tail = 0, 0
	queue[tail] = v
	tail++
	vis[v] = true
	depb[v] = 0
	for head < tail {
		x := queue[head]
		head++
		for _, y := range adj[x] {
			if !vis[y] {
				vis[y] = true
				depb[y] = depb[x] + 1
				queue[tail] = y
				tail++
			}
		}
	}
	// process leaves except u, v
	stack := make([]int, 0, n)
	for i := 1; i <= n; i++ {
		if i != u && i != v && children[i] == 0 {
			stack = append(stack, i)
		}
	}
	xArr := make([]int, 0, n-1)
	yArr := make([]int, 0, n-1)
	delArr := make([]int, 0, n-1)
	var ans int64
	for len(stack) > 0 {
		p := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		par := parent[p]
		children[par]--
		if par != u && par != v && children[par] == 0 {
			stack = append(stack, par)
		}
		xArr = append(xArr, p)
		delArr = append(delArr, p)
		if depa[p] > depb[p] {
			yArr = append(yArr, u)
			ans += int64(depa[p])
		} else {
			yArr = append(yArr, v)
			ans += int64(depb[p])
		}
	}
	// process path from u to v
	cur := v
	for cur != u {
		xArr = append(xArr, u)
		yArr = append(yArr, cur)
		delArr = append(delArr, cur)
		ans += int64(depa[cur])
		cur = parent[cur]
	}
	// output
	fmt.Fprintln(out, ans)
	for i := 0; i < len(xArr); i++ {
		fmt.Fprintln(out, xArr[i], yArr[i], delArr[i])
	}
}
