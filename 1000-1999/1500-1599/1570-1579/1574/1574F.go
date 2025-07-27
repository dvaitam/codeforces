package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int = 998244353

func countOcc(arr []int, sub []int) int {
	if len(sub) == 0 || len(sub) > len(arr) {
		return 0
	}
	cnt := 0
	for i := 0; i+len(sub) <= len(arr); i++ {
		match := true
		for j := 0; j < len(sub); j++ {
			if arr[i+j] != sub[j] {
				match = false
				break
			}
		}
		if match {
			cnt++
		}
	}
	return cnt
}

func check(arr []int, patterns [][]int) bool {
	for _, p := range patterns {
		occP := countOcc(arr, p)
		// generate all non-empty subarrays
		for i := 0; i < len(p); i++ {
			for j := i; j < len(p); j++ {
				sub := p[i : j+1]
				occSub := countOcc(arr, sub)
				if occSub > occP {
					return false
				}
			}
		}
	}
	return true
}

var n, m, k int
var patterns [][]int
var ans int

func dfs(pos int, arr []int) {
	if pos == m {
		if check(arr, patterns) {
			ans++
			if ans >= mod {
				ans -= mod
			}
		}
		return
	}
	for v := 1; v <= k; v++ {
		arr[pos] = v
		dfs(pos+1, arr)
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	fmt.Fscan(in, &n, &m, &k)
	patterns = make([][]int, n)
	for i := 0; i < n; i++ {
		var c int
		fmt.Fscan(in, &c)
		patterns[i] = make([]int, c)
		for j := 0; j < c; j++ {
			fmt.Fscan(in, &patterns[i][j])
		}
	}
	arr := make([]int, m)
	ans = 0
	dfs(0, arr)
	fmt.Println(ans % mod)
}
