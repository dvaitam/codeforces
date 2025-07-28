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

type Move struct {
	x, y, z, c int
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solve(n, m, k int, a [][]int) string {
	var ans []Move
	for i := n; i >= 1; i-- {
		for j := 1; j <= m; j++ {
			for l := 2; l <= n-i+1; l++ {
				ans = append(ans, Move{i, j, l, a[i-1][j-1]})
			}
			for l := i + 1; l <= 2*i-1; l++ {
				ans = append(ans, Move{l, j, n - i + 1, a[i-1][j-1]})
			}
			for l := n - i + 2; l <= n+a[i-1][j-1]; l++ {
				ans = append(ans, Move{2*i - 1, j, l, a[i-1][j-1]})
			}
		}
	}
	for t := 1; t <= k; t++ {
		for i := 1; i < n; i++ {
			if i%2 == 1 || i+1 == n {
				for j := 1; j <= m; j++ {
					ans = append(ans, Move{2 * i, j, t + n - 1, t})
				}
			}
		}
		for i := 1; i < n; i++ {
			ans = append(ans, Move{2 * i, m + 1, t + n - 1, t})
			if i+1 < n {
				ans = append(ans, Move{2*i + 1, m + 1, t + n - 1, t})
			}
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(ans)))
	for _, mv := range ans {
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", mv.x, mv.y, mv.z, mv.c))
	}
	return strings.TrimSpace(sb.String())
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(2) + 2
	m := rng.Intn(2) + 2
	k := rng.Intn(2) + 1
	a := make([][]int, n)
	for i := range a {
		a[i] = make([]int, m)
		for j := range a[i] {
			a[i][j] = rng.Intn(k) + 1
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", a[i][j]))
		}
		sb.WriteByte('\n')
	}
	expected := solve(n, m, k, a)
	return sb.String(), expected
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected:\n%s\ngot:\n%s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
