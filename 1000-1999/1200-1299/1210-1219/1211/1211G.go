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

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		c := make([]int, n+1)
		d := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &c[i])
		}
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &d[i])
		}
		adj := make([][]int, n+1)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}

		diff := make([]int, 0)
		for i := 1; i <= n; i++ {
			if c[i] != d[i] {
				diff = append(diff, i)
			}
		}
		if len(diff) == 0 {
			fmt.Fprintln(out, "Yes")
			fmt.Fprintln(out, 0)
			continue
		}
		if len(diff) == 1 {
			fmt.Fprintln(out, "No")
			continue
		}
		diffSet := make(map[int]bool, len(diff))
		for _, v := range diff {
			diffSet[v] = true
		}
		deg := make(map[int]int, len(diff))
		for _, v := range diff {
			for _, to := range adj[v] {
				if diffSet[to] {
					deg[v]++
				}
			}
			if deg[v] > 2 {
				deg[v] = 3 // mark impossible
			}
		}
		bad := false
		for _, v := range diff {
			if deg[v] > 2 {
				bad = true
				break
			}
		}
		if bad {
			fmt.Fprintln(out, "No")
			continue
		}
		// connectivity check
		vis := make(map[int]bool)
		stack := []int{diff[0]}
		vis[diff[0]] = true
		for len(stack) > 0 {
			v := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			for _, to := range adj[v] {
				if diffSet[to] && !vis[to] {
					vis[to] = true
					stack = append(stack, to)
				}
			}
		}
		if len(vis) != len(diff) {
			fmt.Fprintln(out, "No")
			continue
		}
		ends := make([]int, 0, 2)
		for _, v := range diff {
			if deg[v] == 1 {
				ends = append(ends, v)
			}
		}
		if len(ends) != 2 {
			fmt.Fprintln(out, "No")
			continue
		}
		s, t := ends[0], ends[1]
		// get path s->t inside diff
		parent := make(map[int]int)
		queue := []int{s}
		parent[s] = -1
		for len(queue) > 0 {
			v := queue[0]
			queue = queue[1:]
			if v == t {
				break
			}
			for _, to := range adj[v] {
				if diffSet[to] && parent[to] == 0 {
					parent[to] = v
					queue = append(queue, to)
				}
			}
		}
		if parent[t] == 0 {
			fmt.Fprintln(out, "No")
			continue
		}
		path := []int{}
		x := t
		for x != -1 {
			path = append(path, x)
			x = parent[x]
		}
		// reverse path
		for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
			path[i], path[j] = path[j], path[i]
		}
		k := len(path)
		ok1 := true
		for i := 0; i < k-1; i++ {
			if d[path[i]] != c[path[i+1]] {
				ok1 = false
				break
			}
		}
		if ok1 && d[path[k-1]] == c[path[0]] {
			fmt.Fprintln(out, "Yes")
			fmt.Fprintln(out, k)
			for i, v := range path {
				if i > 0 {
					fmt.Fprint(out, " ")
				}
				fmt.Fprint(out, v)
			}
			fmt.Fprintln(out)
			continue
		}
		ok2 := true
		for i := 1; i < k; i++ {
			if d[path[i]] != c[path[i-1]] {
				ok2 = false
				break
			}
		}
		if ok2 && d[path[0]] == c[path[k-1]] {
			fmt.Fprintln(out, "Yes")
			fmt.Fprintln(out, k)
			for i := k - 1; i >= 0; i-- {
				if i != k-1 {
					fmt.Fprint(out, " ")
				}
				fmt.Fprint(out, path[i])
			}
			fmt.Fprintln(out)
			continue
		}
		fmt.Fprintln(out, "No")
	}
}
