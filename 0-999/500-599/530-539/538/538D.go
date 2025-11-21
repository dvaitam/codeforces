package main

import (
	"bufio"
	"fmt"
	"os"
)

type move struct {
	dx int
	dy int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	grid := make([][]byte, n)
	pieces := make([][2]int, 0)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		row := []byte(s)
		grid[i] = row
		for j := 0; j < n; j++ {
			if row[j] == 'o' {
				pieces = append(pieces, [2]int{i, j})
			}
		}
	}

	if len(pieces) == 0 {
		fmt.Println("NO")
		return
	}

	size := 2*n - 1
	moveBoard := make([][]byte, size)
	for i := 0; i < size; i++ {
		moveBoard[i] = make([]byte, size)
		for j := 0; j < size; j++ {
			moveBoard[i][j] = '.'
		}
	}
	center := n - 1
	moveBoard[center][center] = 'o'

	validMoves := make([]move, 0)
	for dy := -(n - 1); dy <= n-1; dy++ {
		for dx := -(n - 1); dx <= n-1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			valid := true
			hit := false
			for _, p := range pieces {
				nr := p[0] + dy
				nc := p[1] + dx
				if nr < 0 || nr >= n || nc < 0 || nc >= n {
					continue
				}
				cell := grid[nr][nc]
				if cell == '.' {
					valid = false
					break
				}
				if cell == 'x' {
					hit = true
				}
			}
			if valid && hit {
				validMoves = append(validMoves, move{dx: dx, dy: dy})
				moveBoard[center+dy][center+dx] = 'x'
			}
		}
	}

	attacked := make([][]bool, n)
	for i := 0; i < n; i++ {
		attacked[i] = make([]bool, n)
	}

	for _, p := range pieces {
		for _, mv := range validMoves {
			nr := p[0] + mv.dy
			nc := p[1] + mv.dx
			if nr < 0 || nr >= n || nc < 0 || nc >= n {
				continue
			}
			if grid[nr][nc] == 'o' {
				continue
			}
			attacked[nr][nc] = true
		}
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] == '.' && attacked[i][j] {
				fmt.Println("NO")
				return
			}
			if grid[i][j] == 'x' && !attacked[i][j] {
				fmt.Println("NO")
				return
			}
		}
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	fmt.Fprintln(out, "YES")
	for i := 0; i < size; i++ {
		fmt.Fprintln(out, string(moveBoard[i]))
	}
}
