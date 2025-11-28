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
	var ans int64
	for r := 0; r < n; r++ {
		for c := 0; c < m; c++ {
			// Try every possible stripe height h >= 1
			for h := 1; r+3*h <= n; h++ {
				// Try every possible width w >= 1
				for w := 1; c+w <= m; w++ {
					// Check if this h x w flag is valid
					valid := true
					
					// Colors of the three stripes
					c1 := grid[r][c]
					c2 := grid[r+h][c]
					c3 := grid[r+2*h][c]
					
					if c1 == c2 || c2 == c3 {
						valid = false
					} else {
						for i := 0; i < h; i++ {
							for j := 0; j < w; j++ {
								if grid[r+i][c+j] != c1 ||
								   grid[r+h+i][c+j] != c2 ||
								   grid[r+2*h+i][c+j] != c3 {
									valid = false
									break
								}
							}
							if !valid {
								break
							}
						}
					}
					
					if valid {
						ans++
					}
				}
			}
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
