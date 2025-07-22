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

func countStars(grid []string) [][]int {
	n := len(grid)
	m := len(grid[0])
	N := n - 2
	M := m - 2
	star := make([][]int, N)
	for i := 0; i < N; i++ {
		star[i] = make([]int, M)
		for j := 0; j < M; j++ {
			gi := i + 1
			gj := j + 1
			if grid[gi][gj] == '1' && grid[gi-1][gj] == '1' && grid[gi+1][gj] == '1' && grid[gi][gj-1] == '1' && grid[gi][gj+1] == '1' {
				star[i][j] = 1
			}
		}
	}
	return star
}

func solveCase(n, m, k int, grid []string) string {
	N := n - 2
	M := m - 2
	if N <= 0 || M <= 0 {
		return "0"
	}
	star := countStars(grid)
	colSum := make([]int, M)
	S := make([]int, M+1)
	var ans int64
	for l := 0; l < N; l++ {
		for j := 0; j < M; j++ {
			colSum[j] = 0
		}
		for r := l; r < N; r++ {
			for j := 0; j < M; j++ {
				colSum[j] += star[r][j]
			}
			S[0] = 0
			for j := 1; j <= M; j++ {
				S[j] = S[j-1] + colSum[j-1]
			}
			var innerSum int64
			c := 0
			for R := 1; R <= M; R++ {
				if S[R] < k {
					continue
				}
				target := S[R] - k
				for c <= R-1 && S[c] <= target {
					c++
				}
				t := c
				tri := int64(t) * int64(t+1) / 2
				innerSum += tri * int64(M+1-R)
			}
			if innerSum != 0 {
				A := int64(l+1) * int64(N-r)
				ans += A * innerSum
			}
		}
	}
	return fmt.Sprint(ans)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 3
	m := rng.Intn(5) + 3
	grid := make([]string, n)
	for i := 0; i < n; i++ {
		b := make([]byte, m)
		for j := 0; j < m; j++ {
			if rng.Float64() < 0.5 {
				b[j] = '0'
			} else {
				b[j] = '1'
			}
		}
		grid[i] = string(b)
	}
	k := rng.Intn(n*m) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
	for i := 0; i < n; i++ {
		sb.WriteString(grid[i])
		sb.WriteByte('\n')
	}
	input := sb.String()
	expected := solveCase(n, m, k, grid)
	return input, expected
}

func runCase(bin, input, expected string) error {
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
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	if outStr != expected {
		return fmt.Errorf("expected %s got %s", expected, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
