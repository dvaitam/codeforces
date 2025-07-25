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

func checkHorizontal(n, m int, g [][]byte) bool {
	if n%3 != 0 {
		return false
	}
	h := n / 3
	colors := make(map[byte]bool)
	for k := 0; k < 3; k++ {
		col := g[k*h][0]
		for i := k * h; i < (k+1)*h; i++ {
			for j := 0; j < m; j++ {
				if g[i][j] != col {
					return false
				}
			}
		}
		if colors[col] {
			return false
		}
		colors[col] = true
	}
	return len(colors) == 3
}

func checkVertical(n, m int, g [][]byte) bool {
	if m%3 != 0 {
		return false
	}
	w := m / 3
	colors := make(map[byte]bool)
	for k := 0; k < 3; k++ {
		col := g[0][k*w]
		for j := k * w; j < (k+1)*w; j++ {
			for i := 0; i < n; i++ {
				if g[i][j] != col {
					return false
				}
			}
		}
		if colors[col] {
			return false
		}
		colors[col] = true
	}
	return len(colors) == 3
}

func expected(n, m int, g [][]byte) string {
	if checkHorizontal(n, m, g) || checkVertical(n, m, g) {
		return "YES"
	}
	return "NO"
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 1
	m := rng.Intn(6) + 1
	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		row := make([]byte, m)
		for j := 0; j < m; j++ {
			switch rng.Intn(3) {
			case 0:
				row[j] = 'R'
			case 1:
				row[j] = 'G'
			default:
				row[j] = 'B'
			}
		}
		grid[i] = row
	}
	// Occasionally generate a correct flag
	if rng.Intn(2) == 0 {
		if rng.Intn(2) == 0 {
			// horizontal stripes
			n = 3
			m = rng.Intn(6) + 1
			grid = make([][]byte, n)
			colors := []byte{'R', 'G', 'B'}
			rng.Shuffle(3, func(i, j int) { colors[i], colors[j] = colors[j], colors[i] })
			for i := 0; i < 3; i++ {
				row := make([]byte, m)
				for j := 0; j < m; j++ {
					row[j] = colors[i]
				}
				grid[i] = row
			}
		} else {
			// vertical stripes
			m = 3
			n = rng.Intn(6) + 1
			grid = make([][]byte, n)
			colors := []byte{'R', 'G', 'B'}
			rng.Shuffle(3, func(i, j int) { colors[i], colors[j] = colors[j], colors[i] })
			for i := 0; i < n; i++ {
				row := make([]byte, m)
				for j := 0; j < 3; j++ {
					row[j] = colors[j]
				}
				grid[i] = row
			}
		}
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		sb.WriteString(string(grid[i]))
		sb.WriteByte('\n')
	}
	return sb.String(), expected(n, m, grid)
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
