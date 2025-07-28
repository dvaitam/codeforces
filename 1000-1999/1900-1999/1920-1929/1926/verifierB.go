package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func expected(grid []string) string {
	n := len(grid)
	minR, maxR := n, -1
	minC, maxC := n, -1
	ones := 0
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] == '1' {
				ones++
				if i < minR {
					minR = i
				}
				if i > maxR {
					maxR = i
				}
				if j < minC {
					minC = j
				}
				if j > maxC {
					maxC = j
				}
			}
		}
	}
	height := maxR - minR + 1
	width := maxC - minC + 1
	if height == width {
		ok := true
		for i := minR; i <= maxR && ok; i++ {
			for j := minC; j <= maxC; j++ {
				if grid[i][j] != '1' {
					ok = false
					break
				}
			}
		}
		if ok && ones == height*width {
			return "SQUARE"
		}
	}
	return "TRIANGLE"
}

func run(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func makeSquare(n, k, r, c int) []string {
	grid := make([]string, n)
	for i := 0; i < n; i++ {
		row := make([]byte, n)
		for j := range row {
			row[j] = '0'
		}
		grid[i] = string(row)
	}
	for i := 0; i < k; i++ {
		row := []byte(grid[r+i])
		for j := 0; j < k; j++ {
			row[c+j] = '1'
		}
		grid[r+i] = string(row)
	}
	return grid
}

func makeTriangle(n, k, r, c int, down bool) []string {
	grid := make([]string, n)
	for i := 0; i < n; i++ {
		row := make([]byte, n)
		for j := range row {
			row[j] = '0'
		}
		grid[i] = string(row)
	}
	if down {
		for i := 0; i < k; i++ {
			row := []byte(grid[r+i])
			start := c + i
			for j := 0; j < 2*(k-i)-1; j++ {
				row[start+j] = '1'
			}
			grid[r+i] = string(row)
		}
	} else {
		for i := 0; i < k; i++ {
			row := []byte(grid[r+k-1-i])
			start := c + i
			for j := 0; j < 2*(k-i)-1; j++ {
				row[start+j] = '1'
			}
			grid[r+k-1-i] = string(row)
		}
	}
	return grid
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(1)
	var tests [][]string
	n := 7
	// squares
	for size := 2; len(tests) < 50 && size <= 4; size++ {
		for r := 0; r+size <= n && len(tests) < 50; r++ {
			for c := 0; c+size <= n && len(tests) < 50; c++ {
				tests = append(tests, makeSquare(n, size, r, c))
			}
		}
	}
	// triangles
	for size := 2; len(tests) < 100 && size <= 4; size++ {
		for r := 0; r+size <= n && len(tests) < 100; r++ {
			for c := 0; c+2*size-1 <= n && len(tests) < 100; c++ {
				tests = append(tests, makeTriangle(n, size, r, c, false))
			}
		}
	}
	for idx, grid := range tests {
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", len(grid)))
		for _, row := range grid {
			sb.WriteString(row)
			sb.WriteByte('\n')
		}
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		exp := expected(grid)
		if got != exp {
			fmt.Printf("test %d failed: expected=%s got=%s\n", idx+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("ok")
}
