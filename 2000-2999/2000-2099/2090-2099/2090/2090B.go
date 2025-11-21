package main

import (
	"bufio"
	"fmt"
	"os"
)

func solve(grid [][]byte) bool {
	n := len(grid)
	m := len(grid[0])
	// quickly check if size small enough for brute
	if n*m <= 16 {
		return brute(grid)
	}
	ones := []struct{ i, j int }{}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '1' {
				ones = append(ones, struct{ i, j int }{i, j})
			}
		}
	}

	// try pairing sequences
	zero := countZeroRowsCols(grid)
	if !zero {
		return false
	}
	return true
}

func countZeroRowsCols(grid [][]byte) bool {
	n := len(grid)
	m := len(grid[0])
	rowFull := make([]bool, n)
	colFull := make([]bool, m)
	for i := 0; i < n; i++ {
		full := true
		for j := 0; j < m; j++ {
			if grid[i][j] == '0' {
				full = false
				break
			}
		}
		rowFull[i] = full
	}
	for j := 0; j < m; j++ {
		full := true
		for i := 0; i < n; i++ {
			if grid[i][j] == '0' {
				full = false
				break
			}
		}
		colFull[j] = full
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '1' && !rowFull[i] && !colFull[j] {
				return false
			}
		}
	}
	return true
}

func brute(grid [][]byte) bool {
	// BFS state space
	n := len(grid)
	m := len(grid[0])
	type state struct {
		mask int
	}
	total := 1 << (n * m)
	visited := make([]bool, total)
	var encode func([][]byte) int
	encode = func(g [][]byte) int {
		mask := 0
		bit := 1
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if g[i][j] == '1' {
					mask |= bit
				}
				bit <<= 1
			}
		}
		return mask
	}
	var decode func(int) [][]byte
	decode = func(mask int) [][]byte {
		g := make([][]byte, n)
		bit := 1
		for i := 0; i < n; i++ {
			g[i] = make([]byte, m)
			for j := 0; j < m; j++ {
				if mask&bit != 0 {
					g[i][j] = '1'
				} else {
					g[i][j] = '0'
				}
				bit <<= 1
			}
		}
		return g
	}
	start := make([][]byte, n)
	for i := 0; i < n; i++ {
		start[i] = make([]byte, m)
		for j := 0; j < m; j++ {
			start[i][j] = '0'
		}
	}
	queue := []int{encode(start)}
	visited[queue[0]] = true
	target := encode(grid)

	for len(queue) > 0 {
		mask := queue[0]
		queue = queue[1:]
		if mask == target {
			return true
		}
		g := decode(mask)
		// try pushing rows
		for i := 0; i < n; i++ {
			if g[i][m-1] == '1' {
				continue
			}
			newG := make([][]byte, n)
			for r := range newG {
				newG[r] = make([]byte, m)
				copy(newG[r], g[r])
			}
			carry := byte('1')
			for j := 0; j < m; j++ {
				if newG[i][j] == '0' {
					newG[i][j], carry = carry, '0'
					break
				} else {
					newG[i][j], carry = carry, newG[i][j]
				}
			}
			if carry == '1' {
				continue
			}
			newMask := encode(newG)
			if !visited[newMask] {
				visited[newMask] = true
				queue = append(queue, newMask)
			}
		}
		// try pushing columns
		for j := 0; j < m; j++ {
			if g[n-1][j] == '1' {
				continue
			}
			newG := make([][]byte, n)
			for r := range newG {
				newG[r] = make([]byte, m)
				copy(newG[r], g[r])
			}
			carry := byte('1')
			for i := 0; i < n; i++ {
				if newG[i][j] == '0' {
					newG[i][j], carry = carry, '0'
					break
				} else {
					newG[i][j], carry = carry, newG[i][j]
				}
			}
			if carry == '1' {
				continue
			}
			newMask := encode(newG)
			if !visited[newMask] {
				visited[newMask] = true
				queue = append(queue, newMask)
			}
		}
	}
	return false
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		grid := make([][]byte, n)
		for i := 0; i < n; i++ {
			var s string
			fmt.Fscan(in, &s)
			grid[i] = []byte(s)
		}
		if solve(grid) {
			fmt.Fprintln(out, "Yes")
		} else {
			fmt.Fprintln(out, "No")
		}
	}
}
