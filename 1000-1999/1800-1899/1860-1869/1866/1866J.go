package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	n    int
	X, Y int64
	c    []int
	dp   [405][405][2]int64
	vis  [405][405][2]bool
)

func solve(l, r, dir int) int64 {
	if l > r {
		return 0
	}
	if vis[l][r][dir] {
		return dp[l][r][dir]
	}
	vis[l][r][dir] = true
	var res int64
	if dir == 0 {
		col := c[l]
		j := l
		for j <= r && c[j] == col {
			j++
		}
		res = X + solve(j, r, 0)
		for i := j; i <= r; i++ {
			if c[i] == col {
				cost := int64(i-j+1)*Y + solve(j, i-1, 1) + X + solve(i+1, r, 0)
				if cost < res {
					res = cost
				}
			}
		}
	} else {
		col := c[r]
		j := r
		for j >= l && c[j] == col {
			j--
		}
		res = X + solve(l, j, 1)
		for i := j; i >= l; i-- {
			if c[i] == col {
				cost := int64(j-i+1)*Y + solve(i+1, j, 0) + X + solve(l, i-1, 1)
				if cost < res {
					res = cost
				}
			}
		}
	}
	dp[l][r][dir] = res
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	fmt.Fscan(in, &n, &X, &Y)
	c = make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &c[i])
	}
	fmt.Println(solve(0, n-1, 0))
}
