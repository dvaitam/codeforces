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

func run(bin string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveD(grid []string) int {
	n := len(grid)
	sum := make([][]int, n+1)
	for i := range sum {
		sum[i] = make([]int, n+1)
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= n; j++ {
			add := 0
			if grid[i-1][j-1] == '#' {
				add = 1
			}
			sum[i][j] = sum[i-1][j] + sum[i][j-1] - sum[i-1][j-1] + add
		}
	}
	var dp [51][51][51][51]int
	for h := 1; h <= n; h++ {
		for w := 1; w <= n; w++ {
			for x1 := 1; x1+h-1 <= n; x1++ {
				x2 := x1 + h - 1
				for y1 := 1; y1+w-1 <= n; y1++ {
					y2 := y1 + w - 1
					cnt := sum[x2][y2] - sum[x1-1][y2] - sum[x2][y1-1] + sum[x1-1][y1-1]
					if cnt == 0 {
						dp[x1][x2][y1][y2] = 0
					} else {
						best := h
						if w > best {
							best = w
						}
						for k := x1; k < x2; k++ {
							cost := dp[x1][k][y1][y2] + dp[k+1][x2][y1][y2]
							if cost < best {
								best = cost
							}
						}
						for k := y1; k < y2; k++ {
							cost := dp[x1][x2][y1][k] + dp[x1][x2][k+1][y2]
							if cost < best {
								best = cost
							}
						}
						dp[x1][x2][y1][y2] = best
					}
				}
			}
		}
	}
	return dp[1][n][1][n]
}

func generateCase(rng *rand.Rand) ([]byte, int) {
	n := rng.Intn(5) + 1
	grid := make([]string, n)
	var b bytes.Buffer
	fmt.Fprintf(&b, "%d\n", n)
	for i := 0; i < n; i++ {
		row := make([]byte, n)
		for j := 0; j < n; j++ {
			if rng.Intn(2) == 0 {
				row[j] = '.'
			} else {
				row[j] = '#'
			}
		}
		grid[i] = string(row)
		fmt.Fprintln(&b, grid[i])
	}
	expect := solveD(grid)
	return b.Bytes(), expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input, expect := generateCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, string(input))
			os.Exit(1)
		}
		if strings.TrimSpace(out) != fmt.Sprint(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\ninput:\n%s", i, expect, strings.TrimSpace(out), string(input))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
