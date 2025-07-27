package main

import (
	"bufio"
	"fmt"
	"os"
)

func generateStripes(n int, swap bool) []string {
	grid := make([]string, n)
	for i := 0; i < n; i++ {
		row := make([]byte, n)
		for j := 0; j < n; j++ {
			val := (i/2 + j/2) % 2
			if swap {
				val ^= 1
			}
			if val == 0 {
				row[j] = 'S'
			} else {
				row[j] = 'G'
			}
		}
		grid[i] = string(row)
	}
	return grid
}

func generateRings(n int, swap bool) []string {
	grid := make([]string, n)
	for i := 0; i < n; i++ {
		row := make([]byte, n)
		for j := 0; j < n; j++ {
			layer := i
			if j < layer {
				layer = j
			}
			if n-1-i < layer {
				layer = n - 1 - i
			}
			if n-1-j < layer {
				layer = n - 1 - j
			}
			val := layer % 2
			if swap {
				val ^= 1
			}
			if val == 0 {
				row[j] = 'S'
			} else {
				row[j] = 'G'
			}
		}
		grid[i] = string(row)
	}
	return grid
}

func check(grid []string, given []string) bool {
	n := len(grid)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if given[i][j] != '.' && given[i][j] != grid[i][j] {
				return false
			}
		}
	}
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	given := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &given[i])
	}
	if n%2 == 1 {
		fmt.Println("NONE")
		return
	}
	candidates := [][]string{
		generateStripes(n, false),
		generateStripes(n, true),
		generateRings(n, false),
		generateRings(n, true),
	}
	ans := [][]string{}
	for _, g := range candidates {
		if check(g, given) {
			ans = append(ans, g)
		}
	}
	if len(ans) == 0 {
		fmt.Println("NONE")
	} else if len(ans) > 1 {
		fmt.Println("MULTIPLE")
	} else {
		fmt.Println("UNIQUE")
		for i := 0; i < n; i++ {
			fmt.Println(ans[0][i])
		}
	}
}
