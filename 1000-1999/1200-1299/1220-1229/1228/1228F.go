package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func farthest(start int, adj [][]int) (int, []int) {
	n := len(adj) - 1
	parent := make([]int, n+1)
	for i := range parent {
		parent[i] = -2
	}
	q := make([]int, 0, n)
	q = append(q, start)
	parent[start] = -1
	var last int
	for len(q) > 0 {
		u := q[0]
		q = q[1:]
		last = u
		for _, v := range adj[u] {
			if parent[v] != -2 {
				continue
			}
			parent[v] = u
			q = append(q, v)
		}
	}
	return last, parent
}

func treeCenters(adj [][]int) []int {
	a, _ := farthest(1, adj)
	b, parent := farthest(a, adj)
	path := []int{}
	for cur := b; cur != -1; cur = parent[cur] {
		path = append(path, cur)
	}
	if len(path)%2 == 1 {
		return []int{path[len(path)/2]}
	}
	return []int{path[len(path)/2-1], path[len(path)/2]}
}

func verify(adj [][]int, n int, root, p int, leaf bool) bool {
	used := false
	var dfs func(u, par int) int
	dfs = func(u, par int) int {
		heights := []int{}
		for _, v := range adj[u] {
			if v == par {
				continue
			}
			h := dfs(v, u)
			if h < 0 {
				return -1
			}
			heights = append(heights, h)
		}
		if u == p {
			if leaf {
				if len(heights) != 1 || heights[0] != 0 {
					return -1
				}
				used = true
				return 1
			}
			if len(heights) != 3 {
				return -1
			}
			sort.Ints(heights)
			if !(heights[0] == heights[1] && heights[2] == heights[1]+1) {
				return -1
			}
			used = true
			return heights[2] + 1
		}
		if len(heights) == 0 {
			return 0
		}
		if len(heights) == 2 && heights[0] == heights[1] {
			return heights[0] + 1
		}
		return -1
	}
	h := dfs(root, 0)
	return used && h == n-1
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	N := (1 << n) - 2
	adj := make([][]int, N+1)
	for i := 0; i < N-1; i++ {
		var a, b int
		fmt.Fscan(reader, &a, &b)
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}

	deg := make([]int, N+1)
	leaves := 0
	for i := 1; i <= N; i++ {
		deg[i] = len(adj[i])
		if deg[i] == 1 {
			leaves++
		}
	}

	centers := treeCenters(adj)
	candRoots := append([]int(nil), centers...)
	candParents := map[int]struct{}{}
	leafCase := false

	if n == 2 {
		for _, v := range centers {
			candParents[v] = struct{}{}
		}
		leafCase = true
	} else {
		if leaves == (1 << (n - 1)) {
			var deg4 []int
			for i := 1; i <= N; i++ {
				if deg[i] == 4 {
					deg4 = append(deg4, i)
				}
			}
			if len(deg4) > 0 {
				candParents[deg4[0]] = struct{}{}
				candRoots = nil
				for i := 1; i <= N; i++ {
					if deg[i] == 2 {
						candRoots = append(candRoots, i)
					}
				}
			} else {
				for _, v := range centers {
					candParents[v] = struct{}{}
				}
			}
			leafCase = false
		} else if leaves == (1<<(n-1))-1 {
			for i := 1; i <= N; i++ {
				if deg[i] == 2 {
					candParents[i] = struct{}{}
				}
			}
			leafCase = true
		} else {
			fmt.Fprintln(writer, 0)
			return
		}
	}

	result := []int{}
	for p := range candParents {
		for _, r := range candRoots {
			if verify(adj, n, r, p, leafCase) {
				result = append(result, p)
				break
			}
		}
	}
	if len(result) == 0 {
		fmt.Fprintln(writer, 0)
		return
	}
	sort.Ints(result)
	fmt.Fprintln(writer, len(result))
	for i, v := range result {
		if i > 0 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, v)
	}
	writer.WriteByte('\n')
}
