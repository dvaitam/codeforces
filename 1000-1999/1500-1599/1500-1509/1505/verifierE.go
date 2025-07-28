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

type testCase struct {
	input  string
	expect string
}

func solve(H, W int, grid []string) string {
	r, c := 0, 0
	berries := 0
	g := make([][]byte, H)
	for i := range g {
		g[i] = []byte(grid[i])
	}
	if g[r][c] == '*' {
		berries++
		g[r][c] = '.'
	}
	for {
		ti, tj := -1, -1
		best := 1 << 30
		for i := r; i < H; i++ {
			for j := c; j < W; j++ {
				if g[i][j] == '*' {
					d := (i - r) + (j - c)
					if d < best || (d == best && (i < ti || (i == ti && j < tj))) {
						best = d
						ti, tj = i, j
					}
				}
			}
		}
		if ti == -1 {
			break
		}
		for r < ti {
			r++
			if g[r][c] == '*' {
				berries++
				g[r][c] = '.'
			}
		}
		for c < tj {
			c++
			if g[r][c] == '*' {
				berries++
				g[r][c] = '.'
			}
		}
		if g[r][c] == '*' {
			berries++
			g[r][c] = '.'
		}
	}
	for r < H-1 {
		r++
		if g[r][c] == '*' {
			berries++
			g[r][c] = '.'
		}
	}
	for c < W-1 {
		c++
		if g[r][c] == '*' {
			berries++
			g[r][c] = '.'
		}
	}
	return fmt.Sprintf("%d", berries)
}

func genTests() []testCase {
	r := rand.New(rand.NewSource(1))
	tests := make([]testCase, 100)
	for i := range tests {
		H := r.Intn(5) + 1
		W := r.Intn(5) + 1
		grid := make([]string, H)
		for j := 0; j < H; j++ {
			b := make([]byte, W)
			for k := 0; k < W; k++ {
				if r.Intn(2) == 0 {
					b[k] = '.'
				} else {
					b[k] = '*'
				}
			}
			grid[j] = string(b)
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", H, W))
		for j := 0; j < H; j++ {
			sb.WriteString(grid[j])
			sb.WriteByte('\n')
		}
		tests[i].input = sb.String()
		tests[i].expect = solve(H, W, grid)
	}
	return tests
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return out.String() + errBuf.String(), fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	tests := genTests()
	for i, tc := range tests {
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			fmt.Print(tc.input)
			return
		}
		if got != tc.expect {
			fmt.Printf("case %d failed: expected %s got %s\n", i+1, tc.expect, got)
			fmt.Print(tc.input)
			return
		}
	}
	fmt.Println("All tests passed")
}
