package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var s, t string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}
	fmt.Fscan(reader, &t)

	m := len(t)
	pi := make([]int, m)
	for i := 1; i < m; i++ {
		j := pi[i-1]
		for j > 0 && t[i] != t[j] {
			j = pi[j-1]
		}
		if t[i] == t[j] {
			j++
		}
		pi[i] = j
	}

	goTo := make([][26]int, m+1)
	for k := 0; k <= m; k++ {
		for c := 0; c < 26; c++ {
			if k < m && byte('a'+c) == t[k] {
				goTo[k][c] = k + 1
			} else if k == 0 {
				goTo[k][c] = 0
			} else {
				goTo[k][c] = goTo[pi[k-1]][c]
			}
		}
	}

	const inf = int(1e9)
	dp := make([]int, m)
	for i := range dp {
		dp[i] = inf
	}
	dp[0] = 0

	for i := 0; i < len(s); i++ {
		c := s[i] - 'a'
		ndp := make([]int, m)
		for j := range ndp {
			ndp[j] = inf
		}
		for k := 0; k < m; k++ {
			if dp[k] == inf {
				continue
			}
			if dp[k]+1 < ndp[0] {
				ndp[0] = dp[k] + 1
			}
			next := goTo[k][c]
			if next < m && dp[k] < ndp[next] {
				ndp[next] = dp[k]
			}
		}
		dp = ndp
	}

	ans := inf
	for _, v := range dp {
		if v < ans {
			ans = v
		}
	}
	fmt.Println(ans)
}
