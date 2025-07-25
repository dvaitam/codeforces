package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(reader, &s)
		grid[i] = []byte(s)
	}
	// dp2[i][j]: reachable from (i,j) to (n-1,n-1) moving only down or right
	dp2 := make([][]bool, n)
	for i := range dp2 {
		dp2[i] = make([]bool, n)
	}
	for i := n - 1; i >= 0; i-- {
		for j := n - 1; j >= 0; j-- {
			if grid[i][j] != '.' {
				continue
			}
			if i == n-1 && j == n-1 {
				dp2[i][j] = true
			} else if (i+1 < n && dp2[i+1][j]) || (j+1 < n && dp2[i][j+1]) {
				dp2[i][j] = true
			}
		}
	}
	var sb strings.Builder
	i, j := 0, 0
	for !(i == n-1 && j == n-1) {
		if i+1 < n && dp2[i+1][j] {
			sb.WriteByte('D')
			i++
		} else {
			sb.WriteByte('R')
			j++
		}
	}
	fmt.Println(sb.String())
}
