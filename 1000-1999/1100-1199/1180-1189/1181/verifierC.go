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

func solve(n, m int, grid [][]byte) int64 {
	up := make([][]int, n)
	down := make([][]int, n)
	right := make([][]int, n)
	for i := 0; i < n; i++ {
		up[i] = make([]int, m)
		down[i] = make([]int, m)
		right[i] = make([]int, m)
	}
	for j := 0; j < m; j++ {
		for i := 0; i < n; i++ {
			if i > 0 && grid[i][j] == grid[i-1][j] {
				up[i][j] = up[i-1][j] + 1
			} else {
				up[i][j] = 1
			}
		}
	}
	for j := 0; j < m; j++ {
		for i := n - 1; i >= 0; i-- {
			if i+1 < n && grid[i][j] == grid[i+1][j] {
				down[i][j] = down[i+1][j] + 1
			} else {
				down[i][j] = 1
			}
		}
	}
	for i := 0; i < n; i++ {
		for j := m - 1; j >= 0; j-- {
			if j+1 < m && grid[i][j] == grid[i][j+1] {
				right[i][j] = right[i][j+1] + 1
			} else {
				right[i][j] = 1
			}
		}
	}
	var ans int64
	for j := 0; j < m; j++ {
		for i := 0; i < n; i++ {
			h2 := up[i][j]
			midEnd := i
			topEnd := midEnd - h2
			if topEnd < 0 {
				continue
			}
			botStart := midEnd + 1
			if botStart >= n {
				continue
			}
			h1 := up[topEnd][j]
			h3 := down[botStart][j]
			c2 := grid[midEnd][j]
			c1 := grid[topEnd][j]
			c3 := grid[botStart][j]
			if c1 == c2 || c2 == c3 {
				continue
			}
			h := h2
			if h1 < h {
				h = h1
			}
			if h3 < h {
				h = h3
			}
			minW := m
			for u := 0; u < h; u++ {
				rTop := topEnd - u
				if right[rTop][j] < minW {
					minW = right[rTop][j]
				}
				rMid := midEnd - u
				if right[rMid][j] < minW {
					minW = right[rMid][j]
				}
				rBot := botStart + u
				if right[rBot][j] < minW {
					minW = right[rBot][j]
				}
			}
			ans += int64(minW)
		}
	}
	return ans
}

func runCase(bin string, n, m int, grid [][]byte) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		sb.WriteString(string(grid[i]))
		sb.WriteByte('\n')
	}
	input := sb.String()
	expect := fmt.Sprintf("%d", solve(n, m, grid))
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expect {
		return fmt.Errorf("expected %q got %q", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(5) + 1
		m := rng.Intn(5) + 1
		grid := make([][]byte, n)
		for r := 0; r < n; r++ {
			row := make([]byte, m)
			for c := 0; c < m; c++ {
				row[c] = byte('a' + rng.Intn(3))
			}
			grid[r] = row
		}
		if err := runCase(bin, n, m, grid); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
