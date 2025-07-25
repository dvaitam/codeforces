package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

type pair struct{ x, y int }

func solve(grid [][]byte) ([][]byte, bool) {
	n := len(grid)
	m := len(grid[0])
	deg := make([][]int, n)
	for i := 0; i < n; i++ {
		deg[i] = make([]int, m)
	}
	dirs := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] != '.' {
				continue
			}
			cnt := 0
			for _, d := range dirs {
				ni, nj := i+d[0], j+d[1]
				if ni >= 0 && ni < n && nj >= 0 && nj < m && grid[ni][nj] == '.' {
					cnt++
				}
			}
			deg[i][j] = cnt
		}
	}
	q := make([]pair, 0, n*m)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '.' && deg[i][j] == 1 {
				q = append(q, pair{i, j})
			}
		}
	}
	head := 0
	for head < len(q) {
		p := q[head]
		head++
		i, j := p.x, p.y
		if grid[i][j] != '.' {
			continue
		}
		placed := false
		for _, d := range dirs {
			ni, nj := i+d[0], j+d[1]
			if ni >= 0 && ni < n && nj >= 0 && nj < m && grid[ni][nj] == '.' {
				if ni == i {
					if nj == j+1 {
						grid[i][j], grid[ni][nj] = '<', '>'
					} else {
						grid[ni][nj], grid[i][j] = '<', '>'
					}
				} else {
					if ni == i+1 {
						grid[i][j], grid[ni][nj] = '^', 'v'
					} else {
						grid[ni][nj], grid[i][j] = '^', 'v'
					}
				}
				for _, d2 := range dirs {
					xi, yj := i+d2[0], j+d2[1]
					if xi >= 0 && xi < n && yj >= 0 && yj < m && grid[xi][yj] == '.' {
						deg[xi][yj]--
						if deg[xi][yj] == 1 {
							q = append(q, pair{xi, yj})
						}
					}
					xi, yj = ni+d2[0], nj+d2[1]
					if xi >= 0 && xi < n && yj >= 0 && yj < m && grid[xi][yj] == '.' {
						deg[xi][yj]--
						if deg[xi][yj] == 1 {
							q = append(q, pair{xi, yj})
						}
					}
				}
				placed = true
				break
			}
		}
		if !placed {
			return nil, false
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '.' {
				return nil, false
			}
		}
	}
	return grid, true
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesD.txt")
	if err != nil {
		fmt.Println("could not read testcasesD.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	grids := make([][][]byte, t)
	results := make([][]string, t)
	ok := make([]bool, t)
	for c := 0; c < t; c++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		_, _ = strconv.Atoi(scan.Text())
		g := make([][]byte, n)
		for i := 0; i < n; i++ {
			scan.Scan()
			row := scan.Text()
			g[i] = []byte(row)
		}
		res, good := solve(copyGrid(g))
		if good {
			lines := make([]string, n)
			for i := 0; i < n; i++ {
				lines[i] = string(res[i])
			}
			results[c] = lines
		} else {
			results[c] = []string{"Not unique"}
		}
		ok[c] = good
		grids[c] = g
	}
	cmd := exec.Command(os.Args[1])
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("execution failed:", err)
		os.Exit(1)
	}
	outScan := bufio.NewScanner(bytes.NewReader(out))
	outScan.Split(bufio.ScanLines)
	for c := 0; c < t; c++ {
		if ok[c] {
			for i := 0; i < len(results[c]); i++ {
				if !outScan.Scan() {
					fmt.Printf("missing output for test %d\n", c+1)
					os.Exit(1)
				}
				if outScan.Text() != results[c][i] {
					fmt.Printf("test %d failed\n", c+1)
					os.Exit(1)
				}
			}
		} else {
			if !outScan.Scan() {
				fmt.Printf("missing output for test %d\n", c+1)
				os.Exit(1)
			}
			if outScan.Text() != "Not unique" {
				fmt.Printf("test %d failed: expected Not unique\n", c+1)
				os.Exit(1)
			}
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}

func copyGrid(g [][]byte) [][]byte {
	cp := make([][]byte, len(g))
	for i := range g {
		cp[i] = append([]byte(nil), g[i]...)
	}
	return cp
}
