package main

import (
	"bufio"
	"fmt"
	"os"
)

var g [][]int
var belong []int
var leafCnt []int

func dfs(v, p, root int) int {
	belong[v] = root
	cnt := 0
	isLeaf := true
	for _, to := range g[v] {
		if to == p {
			continue
		}
		isLeaf = false
		cnt += dfs(to, v, root)
	}
	if isLeaf {
		cnt = 1
	}
	leafCnt[v] = cnt
	return cnt
}

func subsetSum(weights []int, target int) bool {
	if target < 0 {
		return false
	}
	dp := make([]bool, target+1)
	dp[0] = true
	for _, w := range weights {
		for j := target; j >= w; j-- {
			if dp[j-w] {
				dp[j] = true
			}
		}
	}
	return dp[target]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	var a, b, c, d int
	fmt.Fscan(in, &a, &b, &c, &d)
	g = make([][]int, n+1)
	for i := 2; i <= n; i++ {
		var p int
		fmt.Fscan(in, &p)
		g[i] = append(g[i], p)
		g[p] = append(g[p], i)
	}
	belong = make([]int, n+1)
	leafCnt = make([]int, n+1)

	totalLeaves := 0
	for _, ch := range g[1] {
		totalLeaves += dfs(ch, 1, ch)
	}

	if totalLeaves%2 == 1 {
		fmt.Println("NO")
		return
	}
	half := totalLeaves / 2

	childA := belong[a]
	childB := belong[b] // unused but kept for clarity
	_ = childB
	childC := belong[c]
	childD := belong[d]

	subLeaves := make(map[int]int)
	for _, ch := range g[1] {
		subLeaves[ch] = leafCnt[ch]
	}

	var check func(x, y int) bool
	check = func(x, y int) bool {
		fixed := subLeaves[x] + subLeaves[y]
		target := half - fixed
		others := []int{}
		for _, ch := range g[1] {
			if ch == x || ch == y {
				continue
			}
			others = append(others, subLeaves[ch])
		}
		return subsetSum(others, target)
	}

	if check(childA, childC) || check(childA, childD) {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
