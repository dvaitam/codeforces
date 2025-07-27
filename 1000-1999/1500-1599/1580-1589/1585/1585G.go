package main

import (
	"bufio"
	"fmt"
	"os"
)

func solve(reader *bufio.Reader, writer *bufio.Writer) {
	var n int
	fmt.Fscan(reader, &n)
	children := make([][]int, n+1)
	parent := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &parent[i])
	}
	roots := []int{}
	for i := 1; i <= n; i++ {
		p := parent[i]
		if p == 0 {
			roots = append(roots, i)
		} else {
			children[p] = append(children[p], i)
		}
	}
	rank := make([]int, n+1)
	var dfs func(int) (int, []int)
	dfs = func(v int) (int, []int) {
		if len(children[v]) == 0 {
			rank[v] = 1
			return 1, []int{0}
		}
		if len(children[v]) == 1 {
			gChild, dpChild := dfs(children[v][0])
			rank[v] = rank[children[v][0]] + 1
			dp := append(dpChild, gChild)
			used := make(map[int]bool)
			for _, x := range dp {
				used[x] = true
			}
			g := 0
			for used[g] {
				g++
			}
			return g, dp
		}
		gChild := make([]int, len(children[v]))
		dpChild := make([][]int, len(children[v]))
		minRank := int(1<<30 - 1)
		for i, c := range children[v] {
			g, dp := dfs(c)
			gChild[i] = g
			dpChild[i] = dp
			if rank[c] < minRank {
				minRank = rank[c]
			}
		}
		rank[v] = minRank + 1
		dp := make([]int, rank[v])
		for d := 1; d <= rank[v]; d++ {
			var val int
			if d == 1 {
				for _, g := range gChild {
					val ^= g
				}
			} else {
				for i := range children[v] {
					if d-1 <= len(dpChild[i]) {
						val ^= dpChild[i][d-2]
					}
				}
			}
			dp[d-1] = val
		}
		used := make(map[int]bool)
		for _, x := range dp {
			used[x] = true
		}
		g := 0
		for used[g] {
			g++
		}
		return g, dp
	}
	result := 0
	for _, r := range roots {
		g, _ := dfs(r)
		result ^= g
	}
	if result != 0 {
		fmt.Fprintln(writer, "YES")
	} else {
		fmt.Fprintln(writer, "NO")
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		solve(reader, writer)
	}
}
