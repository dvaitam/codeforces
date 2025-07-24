package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 998244353

var n int
var c []int
var team []int
var ans int

func dfs(idx int, rem []int) {
	if idx == n {
		ans = (ans + 1) % MOD
		return
	}
	for j := 0; j < n; j++ {
		if team[idx] == team[j] { // can't vote for same team
			continue
		}
		if rem[j] == 0 {
			continue
		}
		rem[j]--
		dfs(idx+1, rem)
		rem[j]++
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fscan(reader, &n)
	c = make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &c[i])
	}
	team = make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &team[i])
		team[i]--
	}
	rem := make([]int, n)
	copy(rem, c)
	dfs(0, rem)
	writer := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(writer, ans)
	writer.Flush()
}
