package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveCase(grid [][]byte) int {
	n := len(grid)
	if n == 0 {
		return 0
	}
	m := len(grid[0])

	// Flood-fill '0' cells reachable from boundary (4-connectivity) to find exterior
	exterior := make([][]bool, n)
	for i := range exterior {
		exterior[i] = make([]bool, m)
	}
	queue := make([][2]int, 0)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if (i == 0 || i == n-1 || j == 0 || j == m-1) && grid[i][j] == '0' && !exterior[i][j] {
				exterior[i][j] = true
				queue = append(queue, [2]int{i, j})
			}
		}
	}
	dx4 := [4]int{-1, 1, 0, 0}
	dy4 := [4]int{0, 0, -1, 1}
	for idx := 0; idx < len(queue); idx++ {
		r, c := queue[idx][0], queue[idx][1]
		for d := 0; d < 4; d++ {
			nr, nc := r+dx4[d], c+dy4[d]
			if nr >= 0 && nr < n && nc >= 0 && nc < m && grid[nr][nc] == '0' && !exterior[nr][nc] {
				exterior[nr][nc] = true
				queue = append(queue, [2]int{nr, nc})
			}
		}
	}

	// Check for 2x2 blocks of '1's (valid cycle of length 4 with empty interior)
	best := 0
	for i := 0; i < n-1; i++ {
		for j := 0; j < m-1; j++ {
			if grid[i][j] == '1' && grid[i+1][j] == '1' && grid[i][j+1] == '1' && grid[i+1][j+1] == '1' {
				if 4 > best {
					best = 4
				}
			}
		}
	}

	// Find connected components of interior '0' cells (holes) using 8-connectivity
	visited0 := make([][]bool, n)
	for i := range visited0 {
		visited0[i] = make([]bool, m)
	}
	dx8 := [8]int{-1, -1, -1, 0, 0, 1, 1, 1}
	dy8 := [8]int{-1, 0, 1, -1, 1, -1, 0, 1}

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '0' && !exterior[i][j] && !visited0[i][j] {
				// BFS this interior '0' region with 8-connectivity
				comp := make([][2]int, 0)
				comp = append(comp, [2]int{i, j})
				visited0[i][j] = true
				inBorder := make(map[[2]int]bool)
				for idx := 0; idx < len(comp); idx++ {
					r, c := comp[idx][0], comp[idx][1]
					for d := 0; d < 8; d++ {
						nr, nc := r+dx8[d], c+dy8[d]
						if nr < 0 || nr >= n || nc < 0 || nc >= m {
							continue
						}
						if grid[nr][nc] == '0' && !exterior[nr][nc] && !visited0[nr][nc] {
							visited0[nr][nc] = true
							comp = append(comp, [2]int{nr, nc})
						} else if grid[nr][nc] == '1' {
							inBorder[[2]int{nr, nc}] = true
						}
					}
				}

				borderCells := make([][2]int, 0, len(inBorder))
				for cell := range inBorder {
					borderCells = append(borderCells, cell)
				}

				if len(borderCells) < 4 {
					continue
				}

				// Check that each border cell has exactly 2 border neighbors (4-connectivity)
				valid := true
				for _, cell := range borderCells {
					count := 0
					for d := 0; d < 4; d++ {
						nr, nc := cell[0]+dx4[d], cell[1]+dy4[d]
						if inBorder[[2]int{nr, nc}] {
							count++
						}
					}
					if count != 2 {
						valid = false
						break
					}
				}
				if !valid {
					continue
				}

				// Check that border cells form a single connected cycle (4-connectivity)
				visitedB := make(map[[2]int]bool)
				bq := make([][2]int, 0)
				bq = append(bq, borderCells[0])
				visitedB[borderCells[0]] = true
				for idx := 0; idx < len(bq); idx++ {
					cell := bq[idx]
					for d := 0; d < 4; d++ {
						nr, nc := cell[0]+dx4[d], cell[1]+dy4[d]
						nb := [2]int{nr, nc}
						if inBorder[nb] && !visitedB[nb] {
							visitedB[nb] = true
							bq = append(bq, nb)
						}
					}
				}
				if len(visitedB) != len(borderCells) {
					continue
				}

				if len(borderCells) > best {
					best = len(borderCells)
				}
			}
		}
	}
	return best
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(10) + 1
	m := rng.Intn(10) + 1
	grid := make([][]byte, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		row := make([]byte, m)
		for j := 0; j < m; j++ {
			if rng.Intn(2) == 0 {
				row[j] = '0'
			} else {
				row[j] = '1'
			}
		}
		grid[i] = row
		sb.WriteString(string(row) + "\n")
	}
	expect := solveCase(grid)
	return sb.String(), expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		var val int
		if _, err := fmt.Sscan(out, &val); err != nil {
			fmt.Fprintf(os.Stderr, "case %d bad output: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if val != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", i+1, exp, val, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
