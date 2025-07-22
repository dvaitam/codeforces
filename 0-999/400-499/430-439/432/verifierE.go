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

func expected(n, m int) []string {
	grid := make([][]byte, n)
	for i := range grid {
		grid[i] = make([]byte, m)
	}
	for r := 0; r < n; r++ {
		for c := 0; c < m; c++ {
			if grid[r][c] != 0 {
				continue
			}
			var color byte
			for ch := byte('A'); ch <= 'Z'; ch++ {
				ok := true
				if r > 0 && grid[r-1][c] == ch {
					ok = false
				}
				if c > 0 && grid[r][c-1] == ch {
					ok = false
				}
				if ok {
					color = ch
					break
				}
			}
			maxs := 0
			limit := n - r
			if m-c < limit {
				limit = m - c
			}
			for s := 1; s <= limit; s++ {
				bad := false
				for i := r; i < r+s && !bad; i++ {
					for j := c; j < c+s; j++ {
						if grid[i][j] != 0 {
							bad = true
							break
						}
					}
				}
				if bad {
					break
				}
				if r-1 >= 0 {
					for j := c; j < c+s; j++ {
						if grid[r-1][j] == color {
							bad = true
							break
						}
					}
				}
				if bad {
					break
				}
				if c-1 >= 0 {
					for i := r; i < r+s; i++ {
						if grid[i][c-1] == color {
							bad = true
							break
						}
					}
				}
				if bad {
					break
				}
				if r+s < n {
					for j := c; j < c+s; j++ {
						if grid[r+s][j] == color {
							bad = true
							break
						}
					}
				}
				if bad {
					break
				}
				if c+s < m {
					for i := r; i < r+s; i++ {
						if grid[i][c+s] == color {
							bad = true
							break
						}
					}
				}
				if bad {
					break
				}
				maxs = s
			}
			for i := r; i < r+maxs; i++ {
				for j := c; j < c+maxs; j++ {
					grid[i][j] = color
				}
			}
		}
	}
	res := make([]string, n)
	for i := 0; i < n; i++ {
		res[i] = string(grid[i])
	}
	return res
}

func runCase(exe string, n, m int) error {
	input := fmt.Sprintf("%d %d\n", n, m)
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != n {
		return fmt.Errorf("expected %d lines got %d", n, len(lines))
	}
	for _, line := range lines {
		if len(line) != m {
			return fmt.Errorf("wrong line length")
		}
	}
	expectedGrid := expected(n, m)
	for i := 0; i < n; i++ {
		if lines[i] != expectedGrid[i] {
			return fmt.Errorf("line %d mismatch", i+1)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(5) + 1
		m := rng.Intn(5) + 1
		if err := runCase(exe, n, m); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
