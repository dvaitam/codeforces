package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m, k int
	if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
		return
	}
	grid := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &grid[i])
	}
	var cmds string
	fmt.Fscan(reader, &cmds)
	// flatten grid
	N := n * m
	passable := make([]bool, N)
	var exitID int
	passCnt := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			id := i*m + j
			if grid[i][j] != '#' {
				passable[id] = true
				passCnt++
				if grid[i][j] == 'E' {
					exitID = id
				}
			}
		}
	}
	if passCnt == 1 {
		fmt.Println(0)
		return
	}
	// direction indices: L=0,R=1,U=2,D=3
	inv := make([][][]int, 4)
	for d := 0; d < 4; d++ {
		inv[d] = make([][]int, N)
	}
	// build forward mapping and inverse lists
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			id := i*m + j
			if !passable[id] {
				continue
			}
			// L
			var tgt int
			if j-1 >= 0 && grid[i][j-1] != '#' {
				tgt = i*m + (j - 1)
			} else {
				tgt = id
			}
			inv[0][tgt] = append(inv[0][tgt], id)
			// R
			if j+1 < m && grid[i][j+1] != '#' {
				tgt = i*m + (j + 1)
			} else {
				tgt = id
			}
			inv[1][tgt] = append(inv[1][tgt], id)
			// U
			if i-1 >= 0 && grid[i-1][j] != '#' {
				tgt = (i-1)*m + j
			} else {
				tgt = id
			}
			inv[2][tgt] = append(inv[2][tgt], id)
			// D
			if i+1 < n && grid[i+1][j] != '#' {
				tgt = (i+1)*m + j
			} else {
				tgt = id
			}
			inv[3][tgt] = append(inv[3][tgt], id)
		}
	}
	// BFS over reversed commands
	visited := make([]bool, N)
	visited[exitID] = true
	curr := make([]int, 0, N)
	curr = append(curr, exitID)
	visCnt := 1
	for t := 1; t <= k; t++ {
		var dir int
		switch cmds[t-1] {
		case 'L':
			dir = 0
		case 'R':
			dir = 1
		case 'U':
			dir = 2
		case 'D':
			dir = 3
		default:
			continue
		}
		// process current reachable positions
		newCurr := make([]int, 0, 16)
		for _, x := range curr {
			for _, p := range inv[dir][x] {
				if !visited[p] {
					visited[p] = true
					newCurr = append(newCurr, p)
					visCnt++
				}
			}
		}
		if visCnt == passCnt {
			fmt.Println(t)
			return
		}
		if len(newCurr) > 0 {
			curr = append(curr, newCurr...)
		}
	}
	// not all reached
	fmt.Println(-1)
}
