package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type cell struct{ r, c int }

func stabilize(a [][]int) [][]int {
	n := len(a)
	m := len(a[0])
	q := make([]cell, 0, n*m)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			q = append(q, cell{i, j})
		}
	}
	dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		r, c := cur.r, cur.c
		mx := math.MinInt64
		for _, d := range dirs {
			nr, nc := r+d[0], c+d[1]
			if nr >= 0 && nr < n && nc >= 0 && nc < m {
				if a[nr][nc] > mx {
					mx = a[nr][nc]
				}
			}
		}
		if mx == math.MinInt64 {
			continue
		}
		if a[r][c] > mx {
			a[r][c] = mx
			q = append(q, cell{r, c})
			for _, d := range dirs {
				nr, nc := r+d[0], c+d[1]
				if nr >= 0 && nr < n && nc >= 0 && nc < m {
					q = append(q, cell{nr, nc})
				}
			}
		}
	}
	return a
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 1
	m := rng.Intn(4) + 1
	if n*m == 1 {
		m++
	}
	a := make([][]int, n)
	for i := 0; i < n; i++ {
		a[i] = make([]int, m)
		for j := 0; j < m; j++ {
			a[i][j] = rng.Intn(20) + 1
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d\n", n, m)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", a[i][j])
		}
		sb.WriteByte('\n')
	}
	res := stabilize(cloneMatrix(a))
	var out strings.Builder
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				out.WriteByte(' ')
			}
			fmt.Fprintf(&out, "%d", res[i][j])
		}
		if i+1 < n {
			out.WriteByte('\n')
		}
	}
	out.WriteByte('\n')
	return sb.String(), out.String()
}

func cloneMatrix(a [][]int) [][]int {
	n := len(a)
	m := len(a[0])
	b := make([][]int, n)
	for i := 0; i < n; i++ {
		b[i] = make([]int, m)
		copy(b[i], a[i])
	}
	return b
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected \n%s\ngot \n%s", expected, got)
	}
	return nil
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
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
