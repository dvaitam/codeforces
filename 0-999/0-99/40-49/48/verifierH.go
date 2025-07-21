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

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func solve(n, m, a, b, c int) string {
	H := 2 * n
	W := 2 * m
	grid := make([][]byte, H)
	for i := range grid {
		grid[i] = make([]byte, W)
	}
	type tile struct {
		pat [2][2]byte
		cnt *int
	}
	cntB, cntW, cntM := a, b, c
	patterns := []tile{
		{pat: [2][2]byte{{'B', 'B'}, {'B', 'B'}}, cnt: &cntB},
		{pat: [2][2]byte{{'W', 'W'}, {'W', 'W'}}, cnt: &cntW},
		{pat: [2][2]byte{{'B', 'B'}, {'W', 'W'}}, cnt: &cntM},
		{pat: [2][2]byte{{'W', 'W'}, {'B', 'B'}}, cnt: &cntM},
		{pat: [2][2]byte{{'B', 'W'}, {'B', 'W'}}, cnt: &cntM},
		{pat: [2][2]byte{{'W', 'B'}, {'W', 'B'}}, cnt: &cntM},
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			var top [2]byte
			var left [2]byte
			hasTop := false
			hasLeft := false
			if i > 0 {
				hasTop = true
				r := 2 * i
				c0 := 2 * j
				top[0] = grid[r-1][c0]
				top[1] = grid[r-1][c0+1]
			}
			if j > 0 {
				hasLeft = true
				r0 := 2 * i
				c := 2 * j
				left[0] = grid[r0][c-1]
				left[1] = grid[r0+1][c-1]
			}
			placed := false
			for _, t := range patterns {
				if *t.cnt <= 0 {
					continue
				}
				ok := true
				if hasTop {
					if t.pat[0][0] != top[0] || t.pat[0][1] != top[1] {
						ok = false
					}
				}
				if ok && hasLeft {
					if t.pat[0][0] != left[0] || t.pat[1][0] != left[1] {
						ok = false
					}
				}
				if !ok {
					continue
				}
				r0, c0 := 2*i, 2*j
				for di := 0; di < 2; di++ {
					for dj := 0; dj < 2; dj++ {
						grid[r0+di][c0+dj] = t.pat[di][dj]
					}
				}
				*t.cnt--
				placed = true
				break
			}
			if !placed {
				for k := 2; k < len(patterns); k++ {
					t := &patterns[k]
					if *t.cnt > 0 {
						r0, c0 := 2*i, 2*j
						for di := 0; di < 2; di++ {
							for dj := 0; dj < 2; dj++ {
								grid[r0+di][c0+dj] = t.pat[di][dj]
							}
						}
						*t.cnt--
						placed = true
						break
					}
				}
			}
			if !placed {
				return ""
			}
		}
	}
	var sb strings.Builder
	for i := 0; i < H; i++ {
		sb.Write(grid[i])
		sb.WriteByte('\n')
	}
	return strings.TrimSpace(sb.String())
}

func generateCase(r *rand.Rand) (string, string) {
	n := r.Intn(3) + 1
	m := r.Intn(3) + 1
	a := r.Intn(5)
	b := r.Intn(5)
	c := r.Intn(5)
	input := fmt.Sprintf("%d %d %d %d %d\n", n, m, a, b, c)
	return input, solve(n, m, a, b, c)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(r)
		out, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Printf("Test %d failed.\nInput:\n%sExpected:\n%s\nGot:\n%s\n", i+1, in, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
