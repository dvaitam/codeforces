package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 1000000007

// checkTriple ensures middle is not median for length 3
func checkTriple(a, b, c int) bool {
	if (a < b && b < c) || (a > b && b > c) {
		return false
	}
	return true
}

// checkFive ensures middle element is not median for length 5
func checkFive(x [5]int) bool {
	// median index 2 after sorting
	tmp := [5]int{}
	copy(tmp[:], x[:])
	// insertion sort small array
	for i := 1; i < 5; i++ {
		v := tmp[i]
		j := i - 1
		for j >= 0 && tmp[j] > v {
			tmp[j+1] = tmp[j]
			j--
		}
		tmp[j+1] = v
	}
	return tmp[2] != x[2]
}

// backtracking enumeration -- very slow for large n
func countPerm(n int, fixed []int) int {
	used := make([]bool, n+1)
	for i := 0; i < n; i++ {
		if fixed[i] != -1 {
			used[fixed[i]] = true
		}
	}
	perm := make([]int, n)
	copy(perm, fixed)
	ans := 0
	var dfs func(pos int)
	dfs = func(pos int) {
		if pos == n {
			ans = (ans + 1) % mod
			return
		}
		if perm[pos] != -1 {
			// check constraints with previous numbers
			if pos >= 2 && !checkTriple(perm[pos-2], perm[pos-1], perm[pos]) {
				return
			}
			if pos >= 4 {
				var arr [5]int
				copy(arr[:], perm[pos-4:pos+1])
				if !checkFive(arr) {
					return
				}
			}
			dfs(pos + 1)
			return
		}
		for v := 1; v <= n; v++ {
			if used[v] {
				continue
			}
			perm[pos] = v
			if pos >= 2 && !checkTriple(perm[pos-2], perm[pos-1], perm[pos]) {
				perm[pos] = -1
				continue
			}
			if pos >= 4 {
				var arr [5]int
				copy(arr[:], perm[pos-4:pos+1])
				if !checkFive(arr) {
					perm[pos] = -1
					continue
				}
			}
			used[v] = true
			dfs(pos + 1)
			used[v] = false
			perm[pos] = -1
		}
	}
	dfs(0)
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		fixed := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &fixed[i])
		}
		ans := countPerm(n, fixed)
		fmt.Fprintln(out, ans)
	}
}
