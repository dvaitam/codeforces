package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var N, M int
		fmt.Fscan(in, &N, &M)
		grid := make([][]int, N)
		unsafe := make([][]int, N)
		for i := range grid {
			grid[i] = make([]int, M)
			unsafe[i] = make([]int, M)
		}
		if N%2 == 0 || M%2 == 0 {
			good := [][2]int{}
			for i := 0; i < N; i++ {
				for j := 0; j < M; j++ {
					if (i == 0 && j == 1) || (j == M-1 && i == N-3) || (i == N-1 && j == M-2) || (j == 0 && i == 2) {
						good = append(good, [2]int{i, j})
					}
				}
			}
			if N == 4 && M == 5 {
				good = [][2]int{{0, 1}, {2, 0}, {3, 2}, {2, 4}}
			} else if N == 5 && M == 4 {
				good = [][2]int{{1, 0}, {0, 2}, {2, 3}, {4, 2}}
			}
			cur := 1
			for _, p := range good {
				grid[p[0]][p[1]] = cur
				cur++
			}
			z := N*M + 1 - len(good)*3
			for _, p := range good {
				x, y := p[0], p[1]
				neighbors := [][2]int{}
				for dx := -1; dx <= 1; dx++ {
					for dy := -1; dy <= 1; dy++ {
						if abs(dx)+abs(dy) == 1 {
							neighbors = append(neighbors, [2]int{x + dx, y + dy})
						}
					}
				}
				unsafe[x][y] = -1
				for _, n := range neighbors {
					c, d := n[0], n[1]
					if c >= 0 && c < N && d >= 0 && d < M && grid[c][d] == 0 {
						grid[c][d] = z
						z++
						unsafe[c][d] = 1
					}
				}
			}
			for i := 0; i < N; i++ {
				for j := 0; j < M; j++ {
					if grid[i][j] == 0 {
						grid[i][j] = cur
						cur++
					}
				}
			}
		} else {
			cur := 1
			for i := 0; i < N; i++ {
				for j := 0; j < M; j++ {
					grid[i][j] = cur
					cur++
				}
			}
		}
		for i := 0; i < N; i++ {
			for j := 0; j < M; j++ {
				fmt.Print(grid[i][j], " ")
			}
			fmt.Println()
		}
		played := make([][]int, N)
		allowed := make([][]int, N)
		for i := range played {
			played[i] = make([]int, M)
			allowed[i] = make([]int, M)
		}
		for i := 0; i < N*M; i++ {
			var x, y int
			if i%2 == 0 {
				fmt.Fscan(in, &x, &y)
				x--
				y--
			} else {
				best := [2]int{-1, -1}
				for a := 0; a < N; a++ {
					for b := 0; b < M; b++ {
						if allowed[a][b] == 1 && played[a][b] == 0 && unsafe[a][b] != 1 {
							if best[0] == -1 || grid[a][b] < grid[best[0]][best[1]] {
								best = [2]int{a, b}
							}
						}
					}
				}
				x, y = best[0], best[1]
				fmt.Println(x+1, y+1)
			}
			played[x][y] = 1
			neighbors := [][2]int{}
			for dx := -1; dx <= 1; dx++ {
				for dy := -1; dy <= 1; dy++ {
					if abs(dx)+abs(dy) == 1 {
						nx, ny := x+dx, y+dy
						if nx >= 0 && nx < N && ny >= 0 && ny < M {
							neighbors = append(neighbors, [2]int{nx, ny})
						}
					}
				}
			}
			for _, n := range neighbors {
				allowed[n[0]][n[1]] = 1
			}
			if unsafe[x][y] == -1 {
				unsafe[x][y] = 0
				for _, n := range neighbors {
					unsafe[n[0]][n[1]] = 0
				}
			}
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
