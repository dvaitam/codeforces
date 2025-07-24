package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	const letters = 26
	out := make([]int, letters)
	inDeg := make([]int, letters)
	used := make([]bool, letters)
	for i := 0; i < letters; i++ {
		out[i] = -1
		inDeg[i] = -1
	}
	for ; n > 0; n-- {
		var s string
		fmt.Fscan(in, &s)
		seen := make([]bool, letters)
		for i := 0; i < len(s); i++ {
			c := int(s[i] - 'a')
			used[c] = true
			if seen[c] {
				fmt.Println("NO")
				return
			}
			seen[c] = true
		}
		for i := 0; i+1 < len(s); i++ {
			u := int(s[i] - 'a')
			v := int(s[i+1] - 'a')
			if out[u] != -1 && out[u] != v {
				fmt.Println("NO")
				return
			}
			if inDeg[v] != -1 && inDeg[v] != u {
				fmt.Println("NO")
				return
			}
			out[u] = v
			inDeg[v] = u
		}
	}

	// detect cycles using DFS
	state := make([]int, letters) // 0 unvisited, 1 visiting, 2 done
	var bad bool
	var dfs func(int)
	dfs = func(v int) {
		state[v] = 1
		if out[v] != -1 {
			u := out[v]
			if state[u] == 1 {
				bad = true
				return
			}
			if state[u] == 0 {
				dfs(u)
				if bad {
					return
				}
			}
		}
		state[v] = 2
	}
	for i := 0; i < letters; i++ {
		if used[i] && state[i] == 0 {
			dfs(i)
			if bad {
				fmt.Println("NO")
				return
			}
		}
	}

	visited := make([]bool, letters)
	var chains []string
	for i := 0; i < letters; i++ {
		if used[i] && inDeg[i] == -1 {
			cur := i
			var b []byte
			for cur != -1 && !visited[cur] {
				visited[cur] = true
				b = append(b, byte(cur+'a'))
				cur = out[cur]
			}
			chains = append(chains, string(b))
		}
	}
	for i := 0; i < letters; i++ {
		if used[i] && !visited[i] {
			// Node in cycle or unreachable start (shouldn't happen if no cycle)
			fmt.Println("NO")
			return
		}
	}

	sort.Strings(chains)
	res := ""
	for _, s := range chains {
		res += s
	}
	fmt.Println(res)
}
