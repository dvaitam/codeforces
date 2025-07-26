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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func solve(n, m int, k int64, A [][]int64) int64 {
	d := min(n, m)
	total := n * m
	var ans int64
	if total == 1 {
		if A[0][0] == k {
			ans = 1
		}
		return ans
	}
	C := make([][]map[int64]int, n)
	for i := range C {
		C[i] = make([]map[int64]int, m)
		for j := range C[i] {
			C[i][j] = make(map[int64]int)
		}
	}
	var dfs1 func(r, c, depth, inc, t int, sum int64)
	dfs1 = func(r, c, depth, inc, t int, sum int64) {
		if r < 0 || r >= n || c < 0 || c >= m {
			return
		}
		sum ^= A[r][c]
		if depth == t {
			C[r][c][sum]++
			return
		}
		dfs1(r+inc, c, depth+1, inc, t, sum)
		dfs1(r, c+inc, depth+1, inc, t, sum)
	}
	var dfs2 func(r, c, depth, inc, t int, sum int64)
	dfs2 = func(r, c, depth, inc, t int, sum int64) {
		if r < 0 || r >= n || c < 0 || c >= m {
			return
		}
		sum ^= A[r][c]
		if depth == t {
			if r+inc >= 0 && r+inc < n {
				ans += int64(C[r+inc][c][sum^k])
			}
			if c+inc >= 0 && c+inc < m {
				ans += int64(C[r][c+inc][sum^k])
			}
			return
		}
		dfs2(r+inc, c, depth+1, inc, t, sum)
		dfs2(r, c+inc, depth+1, inc, t, sum)
	}
	dfs1(0, 0, 0, 1, d-1, 0)
	dfs2(n-1, m-1, 0, -1, n+m-2-d, 0)
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(5) + 1
		m := rng.Intn(5) + 1
		k := rng.Int63n(10)
		A := make([][]int64, n)
		input := fmt.Sprintf("%d %d %d\n", n, m, k)
		for r := 0; r < n; r++ {
			A[r] = make([]int64, m)
			for c := 0; c < m; c++ {
				A[r][c] = rng.Int63n(10)
				if c > 0 {
					input += " "
				}
				input += fmt.Sprintf("%d", A[r][c])
			}
			input += "\n"
		}
		expected := fmt.Sprintf("%d", solve(n, m, k, A))
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
