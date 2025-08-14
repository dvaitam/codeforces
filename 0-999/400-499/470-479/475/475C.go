package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]byte, m)
		fmt.Fscan(in, &grid[i])
	}
	sum := make([][]int, n+2)
	for i := range sum {
		sum[i] = make([]int, m+2)
	}
	for i := n - 1; i >= 0; i-- {
		for j := m - 1; j >= 0; j-- {
			val := 0
			if grid[i][j] == 'X' {
				val = 1
			}
			sum[i][j] = val + sum[i+1][j] + sum[i][j+1] - sum[i+1][j+1]
		}
	}
	total := sum[0][0]
	x, y := -1, -1
	for i := 0; i < n && x == -1; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == 'X' {
				x, y = i, j
				break
			}
		}
	}
	if x == -1 {
		fmt.Println(0)
		return
	}
	ans := n*m + 1
	b := 0
	for y+b < m && grid[x][y+b] == 'X' {
		b++
	}
	for a := 1; x+a <= n && grid[x+a-1][y] == 'X'; a++ {
		calc := func(i, j int) int {
			return sum[i][j] - sum[i+a][j] - sum[i][j+b] + sum[i+a][j+b]
		}
		check := func(i, j int) int {
			res := 0
			for {
				if calc(i+1, j) == a*b {
					res += b
					i++
				} else if calc(i, j+1) == a*b {
					res += a
					j++
				} else {
					res += calc(i, j)
					break
				}
			}
			return res
		}
		if check(x, y) == total {
			area := a * b
			if area < ans {
				ans = area
			}
			break
		}
	}
	a := 0
	for x+a < n && grid[x+a][y] == 'X' {
		a++
	}
	for bb := 1; y+bb <= m && grid[x][y+bb-1] == 'X'; bb++ {
		calc := func(i, j int) int {
			return sum[i][j] - sum[i+a][j] - sum[i][j+bb] + sum[i+a][j+bb]
		}
		check := func(i, j int) int {
			res := 0
			for {
				if calc(i+1, j) == a*bb {
					res += bb
					i++
				} else if calc(i, j+1) == a*bb {
					res += a
					j++
				} else {
					res += calc(i, j)
					break
				}
			}
			return res
		}
		if check(x, y) == total {
			area := a * bb
			if area < ans {
				ans = area
			}
		}
	}
	if ans == n*m+1 {
		fmt.Println(-1)
	} else {
		fmt.Println(ans)
	}
}
