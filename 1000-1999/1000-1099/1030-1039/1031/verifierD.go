package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func expectedString(n, k int, grid [][]byte) string {
	inf := n*n + 5
	dist := make([][]int, n)
	for i := 0; i < n; i++ {
		dist[i] = make([]int, n)
		for j := 0; j < n; j++ {
			dist[i][j] = inf
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			cost := 0
			if grid[i][j] != 'a' {
				cost = 1
			}
			if i == 0 && j == 0 {
				dist[i][j] = cost
			} else {
				if i > 0 {
					if d := dist[i-1][j] + cost; d < dist[i][j] {
						dist[i][j] = d
					}
				}
				if j > 0 {
					if d := dist[i][j-1] + cost; d < dist[i][j] {
						dist[i][j] = d
					}
				}
			}
		}
	}
	maxS := 0
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if dist[i][j] <= k {
				s := i + j + 2
				if s > maxS {
					maxS = s
				}
			}
		}
	}
	totalLen := 2*n - 1
	ans := make([]byte, 0, totalLen)
	var frontier [][2]int
	vis := make([][]int, n)
	for i := range vis {
		vis[i] = make([]int, n)
	}
	mark := 1
	if maxS > 0 {
		for i := 0; i < maxS-1; i++ {
			ans = append(ans, 'a')
		}
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if dist[i][j] <= k && i+j+2 == maxS {
					frontier = append(frontier, [2]int{i, j})
					vis[i][j] = mark
				}
			}
		}
	} else {
		ans = append(ans, grid[0][0])
		frontier = append(frontier, [2]int{0, 0})
		vis[0][0] = mark
		maxS = 2
	}
	if len(ans) == totalLen {
		return string(ans)
	}
	for len(ans) < totalLen {
		mark++
		best := byte('z' + 1)
		for _, p := range frontier {
			i, j := p[0], p[1]
			if i+1 < n {
				c := grid[i+1][j]
				if c < best {
					best = c
				}
			}
			if j+1 < n {
				c := grid[i][j+1]
				if c < best {
					best = c
				}
			}
		}
		next := make([][2]int, 0)
		for _, p := range frontier {
			i, j := p[0], p[1]
			if i+1 < n && vis[i+1][j] != mark && grid[i+1][j] == best {
				vis[i+1][j] = mark
				next = append(next, [2]int{i + 1, j})
			}
			if j+1 < n && vis[i][j+1] != mark && grid[i][j+1] == best {
				vis[i][j+1] = mark
				next = append(next, [2]int{i, j + 1})
			}
		}
		ans = append(ans, best)
		frontier = next
	}
	return string(ans)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(15) + 1
	k := rng.Intn(n*n + 1)
	grid := make([][]byte, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i := 0; i < n; i++ {
		row := make([]byte, n)
		for j := 0; j < n; j++ {
			row[j] = byte('a' + rng.Intn(3))
		}
		grid[i] = row
		sb.Write(row)
		sb.WriteByte('\n')
	}
	expect := expectedString(n, k, grid)
	return sb.String(), expect + "\n"
}

func runCase(bin, input, expect string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\nstderr: %s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expect)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
