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

func generateMatrix(rng *rand.Rand, n, m int) [][]int {
	mat := make([][]int, n)
	for i := 0; i < n; i++ {
		mat[i] = make([]int, m)
		for j := 0; j < m; j++ {
			mat[i][j] = rng.Intn(10)
		}
	}
	return mat
}

func swapSub(m [][]int, a, b, c, d, h, w int) {
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			m[a+i][b+j], m[c+i][d+j] = m[c+i][d+j], m[a+i][b+j]
		}
	}
}

func generateCase(rng *rand.Rand) (string, [][]int) {
	n := rng.Intn(4) + 1
	m := rng.Intn(4) + 1
	q := rng.Intn(4) + 1
	mat := generateMatrix(rng, n, m)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, q)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", mat[i][j])
		}
		sb.WriteByte('\n')
	}
	for k := 0; k < q; k++ {
		h := rng.Intn(n) + 1
		w := rng.Intn(m) + 1
		a := rng.Intn(n - h + 1)
		b := rng.Intn(m - w + 1)
		c := rng.Intn(n - h + 1)
		d := rng.Intn(m - w + 1)
		// ensure rectangles don't overlap or touch
		for !(c+h <= a || a+h <= c || d+w <= b || b+w <= d) {
			c = rng.Intn(n - h + 1)
			d = rng.Intn(m - w + 1)
		}
		fmt.Fprintf(&sb, "%d %d %d %d %d %d\n", a+1, b+1, c+1, d+1, h, w)
		swapSub(mat, a, b, c, d, h, w)
	}
	return sb.String(), mat
}

func matrixToString(m [][]int) string {
	var sb strings.Builder
	n := len(m)
	if n == 0 {
		return ""
	}
	for i := 0; i < n; i++ {
		for j := 0; j < len(m[i]); j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", m[i][j])
		}
		if i+1 < n {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func runCase(bin, input string, expected [][]int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outLines := strings.TrimSpace(out.String())
	expectedStr := matrixToString(expected)
	if outLines != expectedStr {
		return fmt.Errorf("expected:\n%s\n\nGot:\n%s", expectedStr, outLines)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
