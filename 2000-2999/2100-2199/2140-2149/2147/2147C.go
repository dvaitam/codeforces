package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		var str string
		fmt.Fscan(in, &n)
		fmt.Fscan(in, &str)

		s := []byte(str)
		forced := make([]bool, n)
		adj := make([][]int, n)

		for i := 0; i < n; i++ {
			if s[i] != '0' {
				continue
			}
			boundary := i == 0 || i == n-1
			leftZero := i > 0 && s[i-1] == '0'
			rightZero := i+1 < n && s[i+1] == '0'
			if !(boundary || leftZero || rightZero) {
				forced[i] = true
			}
		}

		for i := 0; i+2 < n; i++ {
			if s[i] == '0' && s[i+1] == '1' && s[i+2] == '0' {
				adj[i] = append(adj[i], i+2)
				adj[i+2] = append(adj[i+2], i)
			}
		}

		possible := true
		for i := 0; i < n; i++ {
			if s[i] == '0' && forced[i] && len(adj[i]) == 0 {
				possible = false
				break
			}
		}

		visited := make([]bool, n)
		for i := 0; i < n && possible; i++ {
			if s[i] != '0' || len(adj[i]) == 0 || visited[i] {
				continue
			}
			stack := []int{i}
			visited[i] = true
			comp := make([]int, 0)
			for len(stack) > 0 {
				u := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				comp = append(comp, u)
				for _, v := range adj[u] {
					if !visited[v] {
						visited[v] = true
						stack = append(stack, v)
					}
				}
			}
			sort.Ints(comp)
			if !componentPossible(comp, forced) {
				possible = false
				break
			}
		}

		if possible {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}

// componentPossible checks if a path component of '010' linked rabbits
// admits a matching that covers every forced rabbit in it.
func componentPossible(nodes []int, forced []bool) bool {
	k := len(nodes)
	dpFree := true   // node i not matched to the left
	dpTaken := false // node i matched to the left

	for idx := 0; idx < k; idx++ {
		pos := nodes[idx]
		forcedNode := forced[pos]
		nextFree := false
		nextTaken := false

		if dpTaken {
			nextFree = true
		}
		if dpFree {
			if !forcedNode {
				nextFree = true
			}
			if idx+1 < k {
				nextTaken = true
			}
		}

		dpFree, dpTaken = nextFree, nextTaken
	}

	return dpFree
}
