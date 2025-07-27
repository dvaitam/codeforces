package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func buildGrid(n int) [][]int64 {
	g := make([][]int64, n)
	for i := 0; i < n; i++ {
		g[i] = make([]int64, n)
		for j := 0; j < n; j++ {
			if (i+j)%2 == 1 {
				g[i][j] = 1 << uint(i)
			}
		}
	}
	return g
}

func pathFromK(n int, g [][]int64, k int64) []string {
	x, y := 0, 0
	steps := make([]string, 0, 2*n-1)
	steps = append(steps, fmt.Sprintf("%d %d", 1, 1))
	for step := 0; step < 2*n-2; step++ {
		if (x+y)%2 == 0 {
			down := int64(0)
			if x+1 < n {
				down = g[x+1][y]
			}
			if down > 0 && (k&down) != 0 {
				x++
			} else {
				y++
			}
		} else {
			if x+1 < n {
				x++
			} else {
				y++
			}
		}
		steps = append(steps, fmt.Sprintf("%d %d", x+1, y+1))
	}
	return steps
}

func generateIO() (string, string) {
	n := rand.Intn(4) + 2
	q := 100
	grid := buildGrid(n)
	var in strings.Builder
	fmt.Fprintf(&in, "%d\n%d\n", n, q)
	var out strings.Builder
	fmt.Fprintf(&out, "%d\n", n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if j > 0 {
				out.WriteByte(' ')
			}
			fmt.Fprintf(&out, "%d", grid[i][j])
		}
		out.WriteByte('\n')
	}
	for i := 0; i < q; i++ {
		k := rand.Int63n(1 << uint(n-1))
		fmt.Fprintf(&in, "%d\n", k)
		steps := pathFromK(n, grid, k)
		for _, s := range steps {
			fmt.Fprintln(&out, s)
		}
	}
	return in.String(), out.String()
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return buf.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	rand.Seed(5)
	in, exp := generateIO()
	out, err := runBinary(os.Args[1], in)
	if err != nil {
		fmt.Println("Runtime error:", err)
		os.Exit(1)
	}
	if strings.TrimSpace(out) != strings.TrimSpace(exp) {
		fmt.Println("Wrong Answer")
		fmt.Println("Expected:\n" + exp)
		fmt.Println("Got:\n" + out)
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
