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

func maxBishopSum(mat [][]int64) int64 {
	n := len(mat)
	m := len(mat[0])
	d1 := make([]int64, n+m)
	d2 := make([]int64, n+m)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			d1[i+j] += mat[i][j]
			d2[i-j+m-1] += mat[i][j]
		}
	}
	var ans int64
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			sum := d1[i+j] + d2[i-j+m-1] - mat[i][j]
			if sum > ans {
				ans = sum
			}
		}
	}
	return ans
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	mat := make([][]int64, n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d\n", n, m)
	for i := 0; i < n; i++ {
		mat[i] = make([]int64, m)
		for j := 0; j < m; j++ {
			mat[i][j] = int64(rng.Intn(100))
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", mat[i][j])
		}
		sb.WriteByte('\n')
	}
	ans := maxBishopSum(mat)
	return sb.String(), fmt.Sprintf("%d", ans)
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
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
