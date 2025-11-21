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

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		fmt.Fprintln(out, solve(n))
	}
}

func solve(n int) string {
	memo := make([][]int8, n+1)
	for i := range memo {
		memo[i] = make([]int8, 11)
		for j := 0; j < 11; j++ {
			memo[i][j] = -1
		}
	}
	choice := make([][]byte, n)
	for i := range choice {
		choice[i] = make([]byte, 11)
	}

	var dfs func(pos, mod int) bool
	dfs = func(pos, mod int) bool {
		if pos == n {
			return mod == 0
		}
		if memo[pos][mod] != -1 {
			return memo[pos][mod] == 1
		}
		digits := []int{3, 6}
		if pos == n-1 {
			digits = []int{6}
		}
		sign := 1
		if pos%2 == 1 {
			sign = -1
		}
		for _, d := range digits {
			newMod := (mod + sign*d) % 11
			if newMod < 0 {
				newMod += 11
			}
			if dfs(pos+1, newMod) {
				memo[pos][mod] = 1
				choice[pos][mod] = byte('0' + d)
				return true
			}
		}
		memo[pos][mod] = 0
		return false
	}

	if !dfs(0, 0) {
		return "-1"
	}

	result := make([]byte, n)
	mod := 0
	for pos := 0; pos < n; pos++ {
		digit := choice[pos][mod]
		result[pos] = digit
		sign := 1
		if pos%2 == 1 {
			sign = -1
		}
		val := int(digit - '0')
		mod = (mod + sign*val) % 11
		if mod < 0 {
			mod += 11
		}
	}
	return string(result)
}
