package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type caseH struct {
	n, m, k     int
	left, right [][]int
	front, back [][]int
	down, up    [][]int
	input       string
}

func genGrid(rng *rand.Rand, n, m, k int) [][][]int {
	grid := make([][][]int, n)
	for x := 0; x < n; x++ {
		grid[x] = make([][]int, m)
		for y := 0; y < m; y++ {
			grid[x][y] = make([]int, k)
			for z := 0; z < k; z++ {
				if rng.Float64() < 0.5 {
					grid[x][y][z] = 0
				} else {
					grid[x][y][z] = rng.Intn(2) + 1
				}
			}
		}
	}
	return grid
}

func sensorsFromGrid(grid [][][]int) (left, right, front, back [][]int, down, up [][]int) {
	n := len(grid)
	m := len(grid[0])
	k := len(grid[0][0])
	left = make([][]int, m)
	right = make([][]int, m)
	for y := 0; y < m; y++ {
		left[y] = make([]int, k)
		right[y] = make([]int, k)
		for z := 0; z < k; z++ {
			val := 0
			for x := 0; x < n; x++ {
				if grid[x][y][z] != 0 {
					val = grid[x][y][z]
					break
				}
			}
			left[y][z] = val
			val = 0
			for x := n - 1; x >= 0; x-- {
				if grid[x][y][z] != 0 {
					val = grid[x][y][z]
					break
				}
			}
			right[y][z] = val
		}
	}
	front = make([][]int, n)
	back = make([][]int, n)
	for x := 0; x < n; x++ {
		front[x] = make([]int, k)
		back[x] = make([]int, k)
		for z := 0; z < k; z++ {
			val := 0
			for y := 0; y < m; y++ {
				if grid[x][y][z] != 0 {
					val = grid[x][y][z]
					break
				}
			}
			front[x][z] = val
			val = 0
			for y := m - 1; y >= 0; y-- {
				if grid[x][y][z] != 0 {
					val = grid[x][y][z]
					break
				}
			}
			back[x][z] = val
		}
	}
	down = make([][]int, n)
	up = make([][]int, n)
	for x := 0; x < n; x++ {
		down[x] = make([]int, m)
		up[x] = make([]int, m)
		for y := 0; y < m; y++ {
			val := 0
			for z := 0; z < k; z++ {
				if grid[x][y][z] != 0 {
					val = grid[x][y][z]
					break
				}
			}
			down[x][y] = val
			val = 0
			for z := k - 1; z >= 0; z-- {
				if grid[x][y][z] != 0 {
					val = grid[x][y][z]
					break
				}
			}
			up[x][y] = val
		}
	}
	return
}

func genCase(rng *rand.Rand) caseH {
	n := rng.Intn(2) + 1
	m := rng.Intn(2) + 1
	k := rng.Intn(2) + 1
	for n*m*k < 1 || n*m*k > 10 {
		n = rng.Intn(2) + 1
		m = rng.Intn(2) + 1
		k = rng.Intn(2) + 1
	}
	grid := genGrid(rng, n, m, k)
	left, right, front, back, down, up := sensorsFromGrid(grid)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, k)
	for y := 0; y < m; y++ {
		for z := 0; z < k; z++ {
			if z > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(left[y][z]))
		}
		sb.WriteByte('\n')
	}
	for y := 0; y < m; y++ {
		for z := 0; z < k; z++ {
			if z > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(right[y][z]))
		}
		sb.WriteByte('\n')
	}
	for x := 0; x < n; x++ {
		for z := 0; z < k; z++ {
			if z > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(front[x][z]))
		}
		sb.WriteByte('\n')
	}
	for x := 0; x < n; x++ {
		for z := 0; z < k; z++ {
			if z > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(back[x][z]))
		}
		sb.WriteByte('\n')
	}
	for x := 0; x < n; x++ {
		for y := 0; y < m; y++ {
			if y > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(down[x][y]))
		}
		sb.WriteByte('\n')
	}
	for x := 0; x < n; x++ {
		for y := 0; y < m; y++ {
			if y > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(up[x][y]))
		}
		sb.WriteByte('\n')
	}
	return caseH{n, m, k, left, right, front, back, down, up, sb.String()}
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func validate(tc caseH, out string) error {
	if strings.TrimSpace(out) == "-1" {
		return fmt.Errorf("solution exists but got -1")
	}
	tokens := strings.Fields(out)
	if len(tokens) != tc.n*tc.m*tc.k {
		return fmt.Errorf("expected %d numbers, got %d", tc.n*tc.m*tc.k, len(tokens))
	}
	idx := 0
	grid := make([][][]int, tc.n)
	for x := 0; x < tc.n; x++ {
		grid[x] = make([][]int, tc.m)
		for y := 0; y < tc.m; y++ {
			grid[x][y] = make([]int, tc.k)
			for z := 0; z < tc.k; z++ {
				v, err := strconv.Atoi(tokens[idx])
				if err != nil {
					return err
				}
				grid[x][y][z] = v
				idx++
			}
		}
	}
	l, r, f, b, d, u := sensorsFromGrid(grid)
	for y := 0; y < tc.m; y++ {
		for z := 0; z < tc.k; z++ {
			if l[y][z] != tc.left[y][z] || r[y][z] != tc.right[y][z] {
				return fmt.Errorf("LR mismatch")
			}
		}
	}
	for x := 0; x < tc.n; x++ {
		for z := 0; z < tc.k; z++ {
			if f[x][z] != tc.front[x][z] || b[x][z] != tc.back[x][z] {
				return fmt.Errorf("FB mismatch")
			}
		}
	}
	for x := 0; x < tc.n; x++ {
		for y := 0; y < tc.m; y++ {
			if d[x][y] != tc.down[x][y] || u[x][y] != tc.up[x][y] {
				return fmt.Errorf("DU mismatch")
			}
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		if err := validate(tc, got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\noutput:\n%s", i+1, err, tc.input, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
